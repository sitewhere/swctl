/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package action

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/sitewhere/swctl/internal"
	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/instance"
)

// DeleteInstance is the action for creating a SiteWhere instance
type DeleteInstance struct {
	cfg *Configuration

	// Purge Instance data
	Pruge bool
}

// NewDeleteInstance constructs a new *Install
func NewDeleteInstance(cfg *Configuration) *DeleteInstance {
	return &DeleteInstance{
		cfg:   cfg,
		Pruge: false,
	}
}

// Run executes the list command, returning a set of matches.
func (i *DeleteInstance) Run() (*instance.DeleteSiteWhereInstance, error) {
	return &instance.DeleteSiteWhereInstance{}, nil
}

func commandDeleteInstanceRun(cmd *cobra.Command, args []string) {
	// name := args[0]

	// if namespace == "" {
	// 	namespace = name
	// }

	// instance := alpha3.SiteWhereInstance{
	// 	Name:                  name,
	// 	Namespace:             namespace,
	// 	ConfigurationTemplate: "default",
	// 	DatasetTemplate:       "default"}

	// deleteSiteWhereInstance(&instance)
}

func deleteSiteWhereInstance(instance *alpha3.SiteWhereInstance) error {
	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes Config: %v\n", err)
		return err
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return err
	}

	err = deleteSiteWhereMicroservicesResources(instance, client)
	if err != nil {
		fmt.Printf("Error deleting instance: %v\n", err)
		return err
	}

	//TODO
	// if purgeInstance {
	// 	err = deleteSiteWhereResources(instance, client)
	// 	if err != nil {
	// 		fmt.Printf("Error deleting instance: %v\n", err)
	// 		return err
	// 	}

	// 	err = deleteSiteWhereNamespace(instance, config)
	// 	if err != nil {
	// 		fmt.Printf("Error deleting namespace: %v\n", err)
	// 		return err
	// 	}
	// }

	fmt.Printf("SiteWhere Instance '%s' deleted\n", instance.Name)

	return err
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

func deleteSiteWhereNamespace(instance *alpha3.SiteWhereInstance, config *rest.Config) error {
	var err error

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return err
	}

	var ns *v1.Namespace
	ns, err = existNamespace(instance.Namespace, clientset)
	if err != nil {
		fmt.Printf("Error Deleting Namespace: %s, %v", instance.Namespace, err)
		return err
	}

	var namespace = ns.ObjectMeta.Name

	// delete namespace
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})

	return err
}

func existNamespace(namespace string, clientset *kubernetes.Clientset) (*v1.Namespace, error) {
	var err error
	var ns *v1.Namespace

	ns, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return ns, nil
}
