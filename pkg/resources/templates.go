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

// Template for generating a Template Filename
const templateFileTemplate = "/templates/template-%02d.yaml"

// Number of CRD Files
const templatesFileNumber = 37

// InstallSiteWhereTemplates Install SiteWhere Templates CRD
func InstallSiteWhereTemplates(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error
	for i := 1; i <= templatesFileNumber; i++ {
		var templateName = fmt.Sprintf(templateFileTemplate, i)
		err = CreateCustomResourceFromFile(templateName, statikFS, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}
