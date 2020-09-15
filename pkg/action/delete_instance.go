/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"
)

// DeleteInstance is the action for creating a SiteWhere instance
type DeleteInstance struct {
	cfg *Configuration
	// Name of the instance
	InstanceName string
	// Purge Instance data
	Purge bool
}

// NewDeleteInstance constructs a new *Install
func NewDeleteInstance(cfg *Configuration) *DeleteInstance {
	return &DeleteInstance{
		cfg:          cfg,
		InstanceName: "",
		Purge:        false,
	}
}

// Run executes the list command, returning a set of matches.
func (i *DeleteInstance) Run() (*instance.DeleteSiteWhereInstance, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	dynamicClientset, err := i.cfg.KubernetesDynamicClientSet()
	if err != nil {
		return nil, err
	}

	instanceToDelete := alpha3.SiteWhereInstance{
		Name:                  i.InstanceName,
		Namespace:             i.InstanceName,
		ConfigurationTemplate: "default",
		DatasetTemplate:       "default"}

	err = deleteSiteWhereMicroservicesResources(&instanceToDelete, dynamicClientset)
	if err != nil {
		return nil, err
	}
	if i.Purge {
		err = deleteSiteWhereResources(&instanceToDelete, dynamicClientset)
		if err != nil {
			return nil, err
		}
		err = resources.DeleteNamespaceIfExists(i.InstanceName, clientset)
		if err != nil {
			return nil, err
		}
	}
	return &instance.DeleteSiteWhereInstance{}, nil
}

// ExtractInstanceName returns the name of the instance that should be used.
func (i *DeleteInstance) ExtractInstanceName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}

func deleteSiteWhereMicroservicesResources(instance *alpha3.SiteWhereInstance, client dynamic.Interface) error {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(instance.Namespace)

	microservices, err := res.List(context.TODO(), metav1.ListOptions{})

	// delete instance
	if k8serror.IsNotFound(err) {
		errorMessage := fmt.Sprintf("SiteWhere Microservices '%s' not found in namespace '%s'", instance.Name, instance.Namespace)
		return errors.New(errorMessage)
	}
	if err != nil {
		return err
	}

	for _, microservice := range microservices.Items {
		metadata, exists, err := unstructured.NestedMap(microservice.Object, "metadata")

		if err != nil {
			fmt.Printf("Error reading metadata for %s: %v\n", instance.Name, err)
			return nil
		}
		if !exists {
			fmt.Printf("Metadata not found for for SiteWhere Instance: %s", instance.Name)
		} else {
			name, exists, err := unstructured.NestedString(metadata, "name")

			if err != nil {
				fmt.Printf("Error reading metadata for %s: %v\n", instance.Name, err)
				return nil
			}
			if !exists {
				fmt.Printf("Metadata not found for for SiteWhere Instance: %s", instance.Name)
			} else {
				err = res.Delete(context.TODO(), name, metav1.DeleteOptions{})

				if k8serror.IsNotFound(err) {
					errorMessage := fmt.Sprintf("SiteWhere Microservice '%s' not found in namespace '%s'", name, instance.Namespace)
					return errors.New(errorMessage)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func deleteSiteWhereResources(instance *alpha3.SiteWhereInstance, client dynamic.Interface) error {

	res := client.Resource(sitewhereInstanceGVR)

	_, err := res.Get(context.TODO(), instance.Name, metav1.GetOptions{})

	// delete instance
	if k8serror.IsNotFound(err) {
		errorMessage := fmt.Sprintf("SiteWhere Instance '%s' not found", instance.Name)
		return errors.New(errorMessage)
	}
	if err != nil {
		return err
	}
	return res.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
}

func deleteSiteWhereNamespace(instance *alpha3.SiteWhereInstance, clientset kubernetes.Interface) error {
	var err error

	var ns *v1.Namespace
	ns, err = clientset.CoreV1().Namespaces().Get(context.TODO(), instance.Namespace, metav1.GetOptions{})
	if err != nil {
		// fmt.Printf("Error Deleting Namespace: %s, %v", instance.Namespace, err)
		return err
	}

	var namespace = ns.ObjectMeta.Name

	// delete namespace
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})

	return err
}
