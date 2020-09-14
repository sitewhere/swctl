/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"fmt"
	"net/http"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	errors "k8s.io/apimachinery/pkg/api/errors"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// Template for generating a Operator Filename
const operatorFileTemplate = "/operator/operator-%02d.yaml"

// Number of Infrastructure Files
const operatorFileNumber = 23

// InstallSiteWhereOperator Install SiteWhere Operator resource file in the cluster
func InstallSiteWhereOperator(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	_, err = CreateNamespaceIfNotExists(sitewhereSystemNamespace, clientset)
	if err != nil {
		return err
	}

	for i := 1; i <= operatorFileNumber; i++ {
		var operatorResource = fmt.Sprintf(operatorFileTemplate, i)
		err = InstallResourceFromFile(operatorResource, statikFS, clientset, apiextensionsClientset, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	err = waitForDeploymentAvailable(clientset, "sitewhere-operator", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	// if config.IsVerbose() {
	// 	fmt.Print("Deployment sitewhere-operator: ")
	// 	color.Info.Println("Available")
	// }
	err = waitForDeploymentAvailable(clientset, "strimzi-cluster-operator", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	// if config.IsVerbose() {
	// 	fmt.Print("Deployment strimzi-cluster-operator: ")
	// 	color.Info.Println("Available")
	// }
	// if config.IsVerbose() {
	// 	fmt.Print("SiteWhere Operator: ")
	// 	color.Info.Println("Installed")
	// }
	return nil
}

// UninstallSiteWhereOperator Uninstall SiteWhere Operator resource file in the cluster
func UninstallSiteWhereOperator(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	for i := 1; i <= operatorFileNumber; i++ {
		var operatorResource = fmt.Sprintf(operatorFileTemplate, i)
		err = UninstallResourceFromFile(operatorResource, statikFS, clientset, apiextensionsClientset, config)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}

	// if config.Verbose {
	// 	fmt.Print("SiteWhere Operator: ")
	// 	color.Info.Println("Uninstalled")
	// }
	return nil
}
