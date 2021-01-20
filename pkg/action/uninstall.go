/**
 * Copyright © 2014-2021 The SiteWhere Authors
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
	"log"
	"os"

	"github.com/sitewhere/swctl/pkg/install"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

// Uninstall is the action for installing SiteWhere
type Uninstall struct {
	cfg *action.Configuration

	settings *cli.EnvSettings

	// CRD indicates if we need to uninstall SiteWhere Custom Resource Definitions
	CRD bool
	// Infrastructure indicates if we need to install SiteWhere Infrastructure
	Infrastructure bool
	// Operator indicates if we need to install SiteWhere Operator
	Operator bool
	// Template indicates if we need to install SiteWhere templates
	Template bool

	// Minimal installation only install escential SiteWhere components.
	Minimal bool
	// Use verbose mode
	Verbose bool
	// Purge data
	Purge bool
}

// NewUninstall constructs a new *Uninstall
func NewUninstall(cfg *action.Configuration, settings *cli.EnvSettings) *Uninstall {
	return &Uninstall{
		cfg:            cfg,
		settings:       settings,
		CRD:            true,
		Template:       true,
		Operator:       true,
		Infrastructure: true,
		Minimal:        false,
		Verbose:        false,
		Purge:          false,
	}
}

// Run executes the uninstall command, returning the result of the uninstallation
func (i *Uninstall) Run() (*install.SiteWhereInstall, error) {
	var err error
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	return i.uninstallRelease()
}

func (i *Uninstall) uninstallRelease() (*install.SiteWhereInstall, error) {
	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces

	var logConf action.DebugLog
	if i.Verbose {
		logConf = log.Printf
	} else {
		logConf = Discardf
	}

	if err := actionConfig.Init(i.settings.RESTClientGetter(), sitewhereSystemNamespace, os.Getenv("HELM_DRIVER"), logConf); err != nil {
		return nil, err
	}

	uninstallAction := action.NewUninstall(actionConfig)

	res, err := uninstallAction.Run(sitewhereReleaseName)

	if err != nil {
		return nil, err
	}

	return &install.SiteWhereInstall{
		Release:   res.Release.Name,
		Namespace: res.Release.Namespace,
	}, nil
}
