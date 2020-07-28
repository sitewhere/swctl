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

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/kubernetes"
)

// Template for generating a Operator Filename
const operatorFileTemplate = "/operator/operator-%02d.yaml"

// Number of Infrastructure Files
const operatorFileNumber = 23

// InstallSiteWhereOperator Install SiteWhere Operator resource file in the cluster
func InstallSiteWhereOperator(config SiteWhereConfiguration) error {
	var err error

	_, err = CreateNamespaceIfNotExists(sitewhereSystemNamespace, config.GetClientset())
	if err != nil {
		return err
	}

	for i := 1; i <= operatorFileNumber; i++ {
		var operatorResource = fmt.Sprintf(operatorFileTemplate, i)
		err = InstallResourceFromFile(operatorResource, config)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	err = waitForDeploymentAvailable(config.GetClientset(), "sitewhere-operator", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Deploymene sitewhere-operator: Available\n")
	}
	err = waitForDeploymentAvailable(config.GetClientset(), "strimzi-cluster-operator", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Deploymene strimzi-cluster-operator: Available\n")
	}
	if config.IsVerbose() {
		fmt.Printf("SiteWhere Operator: Installed\n")
	}
	return nil
}

// UninstallSiteWhereOperator Uninstall SiteWhere Operator resource file in the cluster
func UninstallSiteWhereOperator(config *SiteWhereInstallConfiguration) error {
	var err error

	clientset, err := kubernetes.NewForConfig(config.KubernetesConfig)
	if err != nil {
		return err
	}

	for i := 1; i <= operatorFileNumber; i++ {
		var operatorResource = fmt.Sprintf(operatorFileTemplate, i)
		err = UninstallResourceFromFile(operatorResource, config.KubernetesConfig, config.StatikFS)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}

	err = DeleteNamespaceIfExists(sitewhereSystemNamespace, clientset)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if config.Verbose {
		fmt.Printf("SiteWhere Operator: Uninstalled\n")
	}
	return nil
}
