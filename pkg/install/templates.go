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
	"fmt"
	"net/http"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	errors "k8s.io/apimachinery/pkg/api/errors"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"

	"github.com/sitewhere/swctl/pkg/resources"
)

// Template for generating a Template Filename
const templateFileTemplate = "/templates/template-%02d.yaml"

// Number of CRD Files
const templatesFileNumber = 37

// SiteWhereTemplates Install SiteWhere Templates CRD
func SiteWhereTemplates(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) ([]SiteWhereTemplateStatus, error) {
	var err error
	var result []SiteWhereTemplateStatus
	for i := 1; i <= templatesFileNumber; i++ {
		var templateName = fmt.Sprintf(templateFileTemplate, i)
		_, err = resources.CreateCustomResourceFromFile(templateName, statikFS, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			return nil, err
		}
		var templateStatus = SiteWhereTemplateStatus{
			Name:   templateName,
			Status: Installed,
		}
		result = append(result, templateStatus)
	}
	return result, nil
}
