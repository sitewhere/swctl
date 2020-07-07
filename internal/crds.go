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
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/rest"
)

// Template for generating a CRD Filename
const crdFileTemplate = "/crd/crd-%02d.yaml"

// Number of CRD Files
const crdFileNumber = 14

// InstallSiteWhereCRDs Install SiteWhere Custom Resource Definitions
func InstallSiteWhereCRDs(config *rest.Config, statikFS http.FileSystem) error {
	var err error
	for i := 1; i <= crdFileNumber; i++ {
		var crdName = fmt.Sprintf(crdFileTemplate, i)
		err = InstallResourceFromFile(crdName, config, statikFS)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// UninstallSiteWhereCRDs Uninstall SiteWhere Custom Resource Definitions
func UninstallSiteWhereCRDs(config *rest.Config, statikFS http.FileSystem) error {
	var err error
	for i := 1; i <= crdFileNumber; i++ {
		var crdName = fmt.Sprintf(crdFileTemplate, i)
		UninstallResourceFromFile(crdName, config, statikFS)
		if err != nil {
			return err
		}
	}
	return nil
}
