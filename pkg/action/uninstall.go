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

package action

import (
	"net/http"

	"github.com/rakyll/statik/fs"

	"k8s.io/apimachinery/pkg/api/errors"

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/sitewhere/swctl/pkg/resources"
	"github.com/sitewhere/swctl/pkg/uninstall"
)

// Uninstall is the action for installing SiteWhere
type Uninstall struct {
	cfg *Configuration

	StatikFS http.FileSystem

	// Minimal installation only install escential SiteWhere components.
	Minimal bool
	// Use verbose mode
	Verbose bool
	// Purge data
	Purge bool
}

// NewUninstall constructs a new *Uninstall
func NewUninstall(cfg *Configuration) *Uninstall {
	statikFS, _ := fs.New()
	return &Uninstall{
		cfg:      cfg,
		StatikFS: statikFS,
		Minimal:  false,
		Verbose:  false,
		Purge:    false,
	}
}

// Run executes the uninstall command, returning the result of the uninstallation
func (i *Uninstall) Run() (*uninstall.SiteWhereUninstall, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	apiextensionsClientset, err := i.cfg.KubernetesAPIExtensionClientSet()
	if err != nil {
		return nil, err
	}
	config, err := i.cfg.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	// Uninstall Infrastructure
	err = resources.UninstallSiteWhereInfrastructure(i.Minimal, i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	// Uninstall Operator
	err = resources.UninstallSiteWhereOperator(i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	// Uninstall Custom Resource Definitions
	err = resources.UninstallSiteWhereCRDs(i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	if i.Minimal {
		err = resources.DeleteSiteWhereNamespaceIfExists(clientset)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	}
	return &uninstall.SiteWhereUninstall{}, nil
}
