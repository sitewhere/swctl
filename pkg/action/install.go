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
)

// Install is the action for installing SiteWhere
type Install struct {
	cfg *Configuration

	StatikFS http.FileSystem

	// Minimal installation only install escential SiteWhere components.
	Minimal bool
	// Use verbose mode
	Verbose bool
}

// NewInstall constructs a new *Install
func NewInstall(cfg *Configuration) *Install {
	statikFS, _ := fs.New()
	return &Install{
		cfg:      cfg,
		StatikFS: statikFS,
		Minimal:  false,
		Verbose:  false,
	}
}

// Run executes the install command, returning the result of the installation
func (i *Install) Run() (*install.SiteWhereInstall, error) {

	// // Install Custom Resource Definitions
	// err = internal.InstallSiteWhereCRDs(sitewhereConfig, i)
	// if err != nil {
	// 	fmt.Printf("Error Installing SiteWhere CRDs: %v\n", err)
	// 	return nil, err
	// }

	// // Install Templates
	// err = internal.InstallSiteWhereTemplates(sitewhereConfig)
	// if err != nil {
	// 	fmt.Printf("Error Installing SiteWhere Templates: %v\n", err)
	// 	return nil, err
	// }

	// // Install Operator
	// err = internal.InstallSiteWhereOperator(sitewhereConfig)
	// if err != nil {
	// 	fmt.Printf("Error Installing SiteWhere Operator: %v\n", err)
	// 	return nil, err
	// }

	// // Install Infrastructure
	// err = internal.InstallSiteWhereInfrastructure(sitewhereConfig)
	// if err != nil {
	// 	fmt.Printf("Error Installing SiteWhere Infrastucture: %v\n", err)
	// 	return nil, err
	// }

	// color.Style{color.FgGreen, color.OpBold}.Println("\nSiteWhere 3.0 Installed")

	return &install.SiteWhereInstall{}, nil
}
