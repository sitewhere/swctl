/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// deleteInstanceCmd represents the instance command
var deleteInstanceCmd = &cobra.Command{
	Use:   "delete instance",
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

func init() {
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

	err = deleteSiteWhereNamespace(instance, config)
	if err != nil {
		fmt.Printf("Error deleting namespace: %v\n", err)
		return err
	}

	fmt.Printf("SiteWhere Instance '%s' deleted\n", instance.Name)

	return err
}

func deleteSiteWhereResources(instance *alpha3.SiteWhereInstance, client dynamic.Interface) error {

	res := client.Resource(sitewhereInstanceGVR)

	sitewhereInstances, err := res.Get(context.TODO(), instance.Name, metav1.GetOptions{})

	// delete instance
	if sitewhereInstances != nil {
		err = res.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
	}

	return err
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
