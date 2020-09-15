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
