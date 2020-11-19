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

package install

import (
	"log"
	"net/http"
	"os"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"

	"github.com/sitewhere/swctl/pkg/resources"
	"github.com/sitewhere/swctl/pkg/status"
)

// Status Status of a installable item.
type Status string

const (
	// Installed The item is installed.
	Installed Status = "Installed"
	// Uninstalled The item is not longer installed.
	Uninstalled = "Uninstalled"
	// Unknown We cannot know if the item is installed or not.
	Unknown = "Unknown"
)

// SiteWhereInstall destribe the installation of SiteWhere.
type SiteWhereInstall struct {
	// Status of SiteWhere CDR installation
	CDRStatuses []status.SiteWhereStatus `json:"crdStatues,omitempty"`
	// Status of SiteWhere Templates installation
	TemplatesStatues []status.SiteWhereStatus `json:"templatesStatues,omitempty"`
	// Status of SiteWhere Operator
	OperatorStatuses []status.SiteWhereStatus `json:"operatorStatuses,omitempty"`
}

func installFiles(statikFS http.FileSystem,
	parentPath string,
	fi os.FileInfo,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) ([]status.SiteWhereStatus, error) {

	var result []status.SiteWhereStatus

	if fi.IsDir() {
		dirName := parentPath + string(os.PathSeparator) + fi.Name()
		log.Printf("Installing Resources from %s", dirName)
		r, err := statikFS.Open(dirName)
		if err != nil {
			return nil, err
		}
		files, err := r.Readdir(-1)
		if err != nil {
			return nil, err
		}
		for _, fileInfo := range files {
			installResult, err := installFiles(statikFS, dirName, fileInfo, clientset, apiextensionsClientset, config)
			if err != nil && !errors.IsAlreadyExists(err) {
				return nil, err
			}
			result = append(result, installResult...)
		}
	} else {
		var fileName = parentPath + string(os.PathSeparator) + fi.Name()
		log.Printf("Installing Resources %s", fileName)
		deployFile, err := statikFS.Open(fileName)
		if err != nil {
			return nil, err
		}
		createObject, err := resources.InstallResourceFromFile(deployFile, fileName, statikFS, clientset, apiextensionsClientset, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			var deployStatus = status.SiteWhereStatus{
				Name:   fileName,
				Status: status.Unknown,
			}
			result = append(result, deployStatus)
		} else {
			var deployStatus = status.SiteWhereStatus{
				Name:       fileName,
				Status:     status.Installed,
				ObjectMeta: createObject,
			}
			result = append(result, deployStatus)
		}
	}
	return result, nil
}
