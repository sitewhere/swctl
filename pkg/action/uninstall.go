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
	"github.com/sitewhere/swctl/pkg/resources"
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

	// Use verbose mode
	Verbose bool
	// Purge data
	Purge bool
}

// NewUninstall constructs a new *Uninstall
func NewUninstall(cfg *action.Configuration, settings *cli.EnvSettings) *Uninstall {
	return &Uninstall{
		cfg:      cfg,
		settings: settings,
		Verbose:  false,
		Purge:    false,
	}
}

// Run executes the uninstall command, returning the result of the uninstallation
func (i *Uninstall) Run() (*install.SiteWhereInstall, error) {
	var err error
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	result, err := i.uninstallRelease()
	if err != nil {
		return nil, err
	}
	if i.Purge {
		clientSet, err := i.cfg.KubernetesClientSet()
		if err != nil {
			return nil, err
		}
		err = resources.DeleteSiteWhereNamespaceIfExists(clientSet)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
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
