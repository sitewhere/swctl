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

	"github.com/sitewhere/swctl/internal/infra"
)

// SiteWhereInfrastructure Uninstall SiteWhere Infrastructure components in the cluster
func SiteWhereInfrastructure(minimal bool,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	if minimal {
		return infra.UninstallSiteWhereInfrastructureMinimal(statikFS, clientset, apiextensionsClientset, config)
	}
	return infra.UninstallSiteWhereInfrastructureDefault(statikFS, clientset, apiextensionsClientset, config)
	// if config.Verbose {
	// 	fmt.Print("SiteWhere Infrastructure: ")
	// 	color.Info.Println("Uninstalled")
	// }
	// return err
}
