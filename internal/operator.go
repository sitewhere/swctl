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

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Template for generating a Operator Filename
const operatorFileTemplate = "/operator/operator-%02d.yaml"

// Number of Infrastructure Files
const operatorFileNumber = 23

// InstallSiteWhereOperator Install SiteWhere Operator resource file in the cluster
func InstallSiteWhereOperator(config *rest.Config, statikFS http.FileSystem) error {
	var err error

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return err
	}

	_, err = CreateNamespaceIfNotExists(sitewhereSystemNamespace, clientset)
	if err != nil {
		return err
	}

	for i := 1; i <= operatorFileNumber; i++ {
		var operatorResource = fmt.Sprintf(operatorFileTemplate, i)
		err = InstallResourceFromFile(operatorResource, config, statikFS)
		if err != nil {
			return err
		}
	}

	return nil
}
