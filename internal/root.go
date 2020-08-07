/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package internal Implements swctl internal use only functions
package internal

import (
	"net/http"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// SiteWhereConfiguration define a configuration
type SiteWhereConfiguration interface {
	IsVerbose() bool
	IsMinimal() bool
	GetConfig() *rest.Config
	GetStatikFS() http.FileSystem
	GetClientset() kubernetes.Interface
	GetApiextensionsClient() apiextensionsclientset.Interface
}

// SiteWhereInstallConfiguration Hold install configuration data
type SiteWhereInstallConfiguration struct {
	Minimal          bool
	Verbose          bool
	KubernetesConfig *rest.Config
	StatikFS         http.FileSystem
}

// InfraTemplateResource template for resources files
type InfraTemplateResource struct {
	FileTemplate string
	FileCount    int
	Enabled      bool
}

// IsVerbose Verbose value
func (c *SiteWhereInstallConfiguration) IsVerbose() bool {
	return c.Verbose
}

// IsMinimal Minimal install
func (c *SiteWhereInstallConfiguration) IsMinimal() bool {
	return c.Minimal
}

// GetConfig Kubernetes Config
func (c *SiteWhereInstallConfiguration) GetConfig() *rest.Config {
	return c.KubernetesConfig
}

// GetStatikFS Statik FS
func (c *SiteWhereInstallConfiguration) GetStatikFS() http.FileSystem {
	return c.StatikFS
}

// GetClientset Kubernetes clienset
func (c *SiteWhereInstallConfiguration) GetClientset() kubernetes.Interface {
	clienset, err := kubernetes.NewForConfig(c.KubernetesConfig)
	if err != nil {
		return nil
	}
	return clienset
}

// GetApiextensionsClient Kubernetes API Extension clienset
func (c *SiteWhereInstallConfiguration) GetApiextensionsClient() apiextensionsclientset.Interface {
	apiextensionsClient, err := apiextensionsclientset.NewForConfig(c.KubernetesConfig)
	if err != nil {
		return nil
	}
	return apiextensionsClient
}
