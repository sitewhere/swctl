/**
 * Copyright © 2014-2020 The SiteWhere Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
