/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package internal Implements swctl internal use only functions
package internal

import (
	"fmt"

	"github.com/gookit/color"

	"k8s.io/apimachinery/pkg/api/errors"
)

// Template for generating a Template Filename
const templateFileTemplate = "/templates/template-%02d.yaml"

// Number of CRD Files
const templatesFileNumber = 37

// InstallSiteWhereTemplates Install SiteWhere Templates CRD
func InstallSiteWhereTemplates(config SiteWhereConfiguration) error {
	var err error
	for i := 1; i <= templatesFileNumber; i++ {
		var templateName = fmt.Sprintf(templateFileTemplate, i)
		CreateCustomResourceFromFile(templateName, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	if config.IsVerbose() {
		fmt.Print("SiteWhere Templates: ")
		color.Info.Println("Installed")
	}
	return nil
}
