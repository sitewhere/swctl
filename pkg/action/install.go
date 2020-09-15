/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package action

import (
	"net/http"

	"github.com/rakyll/statik/fs"

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/sitewhere/swctl/pkg/install"
	"github.com/sitewhere/swctl/pkg/resources"
)

// Install is the action for installing SiteWhere
type Install struct {
	cfg *Configuration

	StatikFS http.FileSystem

	// Minimal installation only install escential SiteWhere components.
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
		cfg:       cfg,
		StatikFS:  statikFS,
		Minimal:   false,
		WaitReady: false,
		Verbose:   false,
	}
}

// Run executes the install command, returning the result of the installation
func (i *Install) Run() (*install.SiteWhereInstall, error) {
	var err error
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
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
	// Install Custom Resource Definitions
	err = resources.InstallSiteWhereCRDs(i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	// Install Templates
	err = resources.InstallSiteWhereTemplates(i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	// Install Operator
	err = resources.InstallSiteWhereOperator(i.WaitReady, i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	// Install Infrastructure
	err = resources.InstallSiteWhereInfrastructure(i.Minimal, i.WaitReady, i.StatikFS, clientset, apiextensionsClientset, config)
	if err != nil {
		return nil, err
	}
	return &install.SiteWhereInstall{}, nil
}
