/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/client-go/rest"

	"github.com/sitewhere/swctl/internal"
	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var (
	purgeInstance = false // Purge instance
	// deleteInstanceCmd represents the instance command
	deleteInstanceCmd = &cobra.Command{
		Use:   "instance",
		Short: "Delete SiteWhere Instance",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires one argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if namespace == "" {
				namespace = name
			}

			instance := alpha3.SiteWhereInstance{
				Name:                  name,
				Namespace:             namespace,
				ConfigurationTemplate: "default",
				DatasetTemplate:       "default"}

			deleteSiteWhereInstance(&instance)
		},
	}
)

func init() {
	deleteInstanceCmd.Flags().BoolVarP(&purgeInstance, "purge", "p", false, "Purge instance.")
	deleteCmd.AddCommand(deleteInstanceCmd)
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

	err = deleteSiteWhereResources(instance, client)
	if err != nil {
		fmt.Printf("Error deleting instance: %v\n", err)
		return err
	}

	if purgeInstance {
		err = deleteSiteWhereNamespace(instance, config)
		if err != nil {
			fmt.Printf("Error deleting namespace: %v\n", err)
			return err
		}
	}

	fmt.Printf("SiteWhere Instance '%s' deleted\n", instance.Name)

	return err
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
