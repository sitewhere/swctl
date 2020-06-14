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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// SiteWhereInstance represents an Instacen in SiteWhere
type SiteWhereInstance struct {
	Name                  string
	Namespace             string
	ConfigurationTemplate string
	DatasetTemplate       string
}

var clientset *kubernetes.Clientset
var apixClient *apixv1beta1client.ApiextensionsV1beta1Client
var (
	sitewhereInstanceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha3",
		Resource: "instances",
	}
)

var frmtAttr = "%-25s: %-32s\n"

// instancesCmd represents the instances command
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage SiteWhere Instance",
	Long:  `Manage SiteWhere Instance.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("requires one or zero arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			handleListInstances()
		} else {
			name := args[0]
			handleInstance(name)
		}
	},
}

func init() {
	rootCmd.AddCommand(instancesCmd)
}

func handleListInstances() {
	var err error

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := getKubeConfig(kubeconfig)
	client, err := dynamic.NewForConfig(config)
	res := client.Resource(sitewhereInstanceGVR)
	options := metav1.ListOptions{}
	sitewhereInstaces, err := res.List(context.TODO(), options)

	if err != nil {
		log.Printf("Error reading SiteWhere Instances: %v", err)
		return
	}

	template := "%-20s%-20s%-20s%-20s\n"
	fmt.Printf(template, "NAME", "NAMESPACE", "CONFIG TMPL", "DATESET TMPL")

	for _, instance := range sitewhereInstaces.Items {
		sitewhereInstace := extractFromResource(&instance)
		fmt.Printf(
			template,
			sitewhereInstace.Name,
			sitewhereInstace.Namespace,
			sitewhereInstace.ConfigurationTemplate,
			sitewhereInstace.DatasetTemplate)
	}
}

func handleInstance(instanceName string) {
	var err error

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := getKubeConfig(kubeconfig)
	client, err := dynamic.NewForConfig(config)
	res := client.Resource(sitewhereInstanceGVR)
	options := metav1.GetOptions{}
	sitewhereInstace, err := res.Get(context.TODO(), instanceName, options)

	if err != nil {
		fmt.Printf(
			"SiteWhere Instace %s NOT FOUND.\n",
			instanceName)
		return
	}

	printSiteWhereInstance(sitewhereInstace)
}

func printSiteWhereInstance(crSiteWhereInstace *unstructured.Unstructured) {
	sitewhereInstace := extractFromResource(crSiteWhereInstace)

	fmt.Printf(frmtAttr, "Instance Name", sitewhereInstace.Name)
	fmt.Printf(frmtAttr, "Instance Namespace", sitewhereInstace.Namespace)
	fmt.Printf(frmtAttr, "Configuration Template", sitewhereInstace.ConfigurationTemplate)
	fmt.Printf(frmtAttr, "Dataset Template", sitewhereInstace.DatasetTemplate)
}

func extractFromResource(crSiteWhereInstace *unstructured.Unstructured) *SiteWhereInstance {
	metadata, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "metadata")
	var result = SiteWhereInstance{}

	if err != nil {
		log.Printf("Error reading metadata for %s: %v", crSiteWhereInstace, err)
		return nil
	}
	if !exists {
		log.Printf("Metadata not found for for SiteWhere Instance: %s", crSiteWhereInstace)
	} else {
		extractSiteWhereInstanceMetadata(metadata, &result)
	}
	spec, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "spec")
	if err != nil {
		log.Printf("Error reading spec for %s: %v", crSiteWhereInstace, err)
		return nil
	}
	if !exists {
		log.Printf("Spec not found for for SiteWhere Instance: %s", crSiteWhereInstace)
	} else {
		extractSiteWhereInstanceSpec(spec, &result)
	}

	return &result
}

func extractSiteWhereInstanceMetadata(metadata map[string]interface{}, instance *SiteWhereInstance) {
	name := extractSiteWhereInstanceName(metadata)
	instance.Name = fmt.Sprintf("%v", name)
}

func extractSiteWhereInstanceName(metadata map[string]interface{}) interface{} {
	return metadata["name"]
}

func extractSiteWhereInstanceSpec(spec map[string]interface{}, instance *SiteWhereInstance) {
	instanceNamespace := spec["instanceNamespace"]
	configurationTemplate := spec["configurationTemplate"]
	datasetTemplate := spec["datasetTemplate"]

	instance.Namespace = fmt.Sprintf("%v", instanceNamespace)
	instance.ConfigurationTemplate = fmt.Sprintf("%v", configurationTemplate)
	instance.DatasetTemplate = fmt.Sprintf("%v", datasetTemplate)
}

// Buid a Kubernetes Config from a filepath
func getKubeConfig(pathToCfg string) (*rest.Config, error) {
	if pathToCfg == "" {
		// in cluster access
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", pathToCfg)
}

func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(pathToCfg)

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getClientV1Beta1(pathToCfg string) (*apixv1beta1client.ApiextensionsV1beta1Client, error) {
	config, err := getKubeConfig(pathToCfg)

	if err != nil {
		return nil, err
	}

	return apixv1beta1client.NewForConfig(config)
}
