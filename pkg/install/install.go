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

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/sitewhere/swctl/pkg/kube"
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
	// Status of SiteWhere Infrastructure
	InfrastructureStatuses []status.SiteWhereStatus `json:"infrastructureStatuses,omitempty"`
}

func installFiles(statikFS http.FileSystem,
	parentPath string,
	fi os.FileInfo,
	KubeClient kube.Interface) ([]status.SiteWhereStatus, error) {

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
			installResult, err := installFiles(statikFS, dirName, fileInfo, KubeClient)
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
		// Open the resource file
		res, err := KubeClient.Build(deployFile, false)
		if err != nil {
			return nil, err
		}
		if _, err := KubeClient.Create(res); err != nil {
			// If the error is Resource already exists, continue.
			if errors.IsAlreadyExists(err) {
				log.Printf("Resource %s is already present. Skipping.", fileName)
			} else {
				var deployStatus = status.SiteWhereStatus{
					Name:   fileName,
					Status: status.Unknown,
				}
				result = append(result, deployStatus)
			}
		} else {
			var deployStatus = status.SiteWhereStatus{
				Name:   fileName,
				Status: status.Installed,
				//		ObjectMeta: createObject,
			}
			result = append(result, deployStatus)
		}

	}
	return result, nil
}
