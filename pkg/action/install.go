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

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/sitewhere/swctl/pkg/install"
)

// Install is the action for installing SiteWhere
type Install struct {
	cfg *Configuration

	StatikFS http.FileSystem
	// CRD indicates if we need to install SiteWhere Custom Resource Definitions
	CRD bool
	// Template indicates if we need to install SiteWhere templates
	Template bool
	// Operator indicates if we need to install SiteWhere Operator
	Operator bool
	// Infrastructure indicates if we need to install SiteWhere Infrastructure
	Infrastructure bool
	// Minimal installation only install escential SiteWhere components
	Minimal bool
	// Wait for components to be ready before return control.
	WaitReady bool
	// Use verbose mode
	Verbose bool
}

// NewInstall constructs a new *Install
func NewInstall(cfg *Configuration) *Install {
	statikFS, _ := fs.New()
	return &Install{
		cfg:            cfg,
		StatikFS:       statikFS,
		CRD:            true,
		Template:       true,
		Operator:       true,
		Infrastructure: true,
		Minimal:        false,
		WaitReady:      false,
		Verbose:        false,
	}
}

// Run executes the install command, returning the result of the installation
func (i *Install) Run() (*install.SiteWhereInstall, error) {
	var err error
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	//clientSet, err := i.cfg.KubernetesClientSet()
	// if err != nil {
	// 	return nil, err
	// }
	// extensionsClients, err := i.cfg.KubernetesAPIExtensionClientSet()
	// if err != nil {
	// 	return nil, err
	// }
	// config, err := i.cfg.RESTClientGetter.ToRESTConfig()
	// if err != nil {
	// 	return nil, err
	// }
	var crdStatues []install.SiteWhereCRDStatus
	if i.CRD {
		// Install Custom Resource Definitions
		// crdStatues, err = install.SiteWhereCRDs(i.StatikFS, clientSet, extensionsClients, config)
		// if err != nil {
		// 	return nil, err
		// }
	}
	var templatesStatues []install.SiteWhereTemplateStatus
	if i.Template {
		// Install Templates
		// templatesStatues, err = install.SiteWhereTemplates(i.StatikFS, clientSet, extensionsClients, config)
		// if err != nil {
		// 	return nil, err
		// }
	}
	var operatorStatuses []install.SiteWhereOperatorStatus
	if i.Operator {
		// Install Operator
		// operatorStatuses, err = install.SiteWhereOperator(i.WaitReady, i.StatikFS, clientSet, extensionsClients, config)
		// if err != nil {
		// 	return nil, err
		// }
	}
	if i.Infrastructure {
		// Install Infrastructure
		// err = install.SiteWhereInfrastructure(i.Minimal, i.WaitReady, i.StatikFS, clientSet, extensionsClients, config)
		// if err != nil {
		// 	return nil, err
		// }
	}
	return &install.SiteWhereInstall{
		CDRStatues:       crdStatues,
		TemplatesStatues: templatesStatues,
		OperatorStatuses: operatorStatuses,
	}, nil
}
