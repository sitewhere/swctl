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

// Template for generating a Infrastucture Filename
const infraFileTemplate = "/infra-min/infra-min-%02d.yaml"

// Number of Infrastructure Files
const infraFileNumber = 28

// InstallSiteWhereInfrastructure Install SiteWhere Infrastructure components in the cluster
func InstallSiteWhereInfrastructure(config *rest.Config, statikFS http.FileSystem) error {
	var err error
	for i := 1; i <= infraFileNumber; i++ {
		var infraResource = fmt.Sprintf(infraFileTemplate, i)
		err = InstallResourceFromFile(infraResource, config, statikFS)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// UninstallSiteWhereInfrastructure Uninstall SiteWhere Infrastructure components in the cluster
func UninstallSiteWhereInfrastructure(config *rest.Config, statikFS http.FileSystem) error {
	var err error
	for i := 1; i <= infraFileNumber; i++ {
		var infraResource = fmt.Sprintf(infraFileTemplate, i)
		err = UninstallResourceFromFile(infraResource, config, statikFS)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}
