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

package uninstall

import (
	"net/http"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	//"github.com/sitewhere/swctl/internal/crds"
	//"github.com/sitewhere/swctl/pkg/resources"
)

// SiteWhereCRDs Uninstall SiteWhere Custom Resource Definitions
func SiteWhereCRDs(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	// var err error
	// for _, crdFile := range crds.GetSiteWhereCRDFiles() {
	// resources.UninstallResourceFromFile(crdFile, statikFS, clientset, apiextensionsClientset, config)
	// if err != nil {
	// 	return err
	// }
	//	}
	// if config.IsVerbose() {
	// 	fmt.Print("SiteWhere Custom Resources Definition: ")
	// 	color.Info.Println("Uninstalled")
	// }
	return nil
}
