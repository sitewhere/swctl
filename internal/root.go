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
