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

	"k8s.io/apimachinery/pkg/api/errors"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// Template for generating a CRD Filename
const crdFileTemplate = "/crd/crd-%02d.yaml"

// Number of CRD Files
const crdFileNumber = 14

// InstallSiteWhereCRDs Install SiteWhere Custom Resource Definitions
func InstallSiteWhereCRDs(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error
	for i := 1; i <= crdFileNumber; i++ {
		var fileName = fmt.Sprintf(crdFileTemplate, i)
		err = InstallResourceFromFile(fileName, statikFS, clientset, apiextensionsClientset, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// UninstallSiteWhereCRDs Uninstall SiteWhere Custom Resource Definitions
func UninstallSiteWhereCRDs(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error
	for i := 1; i <= crdFileNumber; i++ {
		var crdName = fmt.Sprintf(crdFileTemplate, i)
		UninstallResourceFromFile(crdName, statikFS, clientset, apiextensionsClientset, config)
		if err != nil {
			return err
		}
	}
	// if config.IsVerbose() {
	// 	fmt.Print("SiteWhere Custom Resources Definition: ")
	// 	color.Info.Println("Uninstalled")
	// }
	return nil
}
