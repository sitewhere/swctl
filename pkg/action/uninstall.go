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
	"github.com/spf13/cobra"

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
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
	return &uninstall.SiteWhereUninstall{}, nil
}

// uninstallSiteWhereCommand Performs the steps necessary to uninstall SiteWhere
func uninstallSiteWhereCommand(_ *cobra.Command, _ []string) {
	// config, err := internal.GetKubeConfigFromKubeconfig()
	// if err != nil {
	// 	fmt.Printf("Error getting Kubernetes Config: %v\n", err)
	// 	return
	// }

	// statikFS, err := fs.New()
	// if err != nil {
	// 	fmt.Printf("Error Reading Resources: %v\n", err)
	// 	return
	// }

	// var sitewhereConfig = internal.SiteWhereInstallConfiguration{
	// 	Minimal:          minimalUninstall,
	// 	Verbose:          verboseUninstall,
	// 	KubernetesConfig: config,
	// 	StatikFS:         statikFS,
	// }

	// // Uninstall Infrastructure
	// err = internal.UninstallSiteWhereInfrastructure(&sitewhereConfig)
	// if err != nil {
	// 	fmt.Printf("Error Uninstalling SiteWhere Infrastucture: %v\n", err)
	// 	return
	// }

	// // Uninstall Operator
	// err = internal.UninstallSiteWhereOperator(&sitewhereConfig)
	// if err != nil {
	// 	fmt.Printf("Error Uninstalling SiteWhere Operator: %v\n", err)
	// 	return
	// }

	// // Uninstall Custom Resource Definitions
	// internal.UninstallSiteWhereCRDs(&sitewhereConfig)

	// if purge {
	// 	err = internal.DeleteSiteWhereNamespaceIfExists(sitewhereConfig.GetClientset())
	// 	if err != nil && !errors.IsNotFound(err) {
	// 		fmt.Printf("Error Uninstalling SiteWhere Namespace: %v\n", err)
	// 		return
	// 	}
	// }

	// color.Style{color.FgGreen, color.OpBold}.Println("\nSiteWhere 3.0 Uninstalled")
}
