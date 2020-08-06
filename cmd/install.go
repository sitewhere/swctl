/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/internal"
	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
)

// installCmd represents the install command
var (
	minimalInstall = false // Use minimal install profile. Initialize only essential microservices.
	verboseInstall = false // Use verbose installation
	installCmd     = &cobra.Command{
		Use:   "install",
		Short: "Install SiteWhere CRD and Operator",
		Long: `Use this command to install SiteWhere 3.0 on a Kubernetes Cluster.
This command will install:
 - SiteWhere System Namespace: sitewhere-system (default)
 - SiteWhere Custom Resources Definition.
 - SiteWhere Templates.
 - SiteWhere Operator.
 - SiteWhere Infrastructure.`,
		Run: installSiteWhereCommand,
	}
)

func init() {
	installCmd.Flags().BoolVarP(&minimalInstall, "minimal", "m", false, "Minimal installation.")
	installCmd.Flags().BoolVarP(&verboseInstall, "verbose", "v", false, "Verbose installation.")
	rootCmd.AddCommand(installCmd)
}

func installSiteWhereCommand(cmd *cobra.Command, args []string) {
	var err error

	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes Config: %v\n", err)
		return
	}

	statikFS, err := fs.New()
	if err != nil {
		fmt.Printf("Error Reading Resources: %v\n", err)
		return
	}

	var sitewhereConfig internal.SiteWhereConfiguration = &internal.SiteWhereInstallConfiguration{
		Minimal:          minimalInstall,
		Verbose:          verboseInstall,
		KubernetesConfig: config,
		StatikFS:         statikFS,
	}

	// Install Custom Resource Definitions
	err = internal.InstallSiteWhereCRDs(sitewhereConfig)
	if err != nil {
		fmt.Printf("Error Installing SiteWhere CRDs: %v\n", err)
		return
	}

	// Install Templates
	err = internal.InstallSiteWhereTemplates(sitewhereConfig)
	if err != nil {
		fmt.Printf("Error Installing SiteWhere Templates: %v\n", err)
		return
	}

	// Install Operator
	err = internal.InstallSiteWhereOperator(sitewhereConfig)
	if err != nil {
		fmt.Printf("Error Installing SiteWhere Operator: %v\n", err)
		return
	}

	// Install Infrastructure
	err = internal.InstallSiteWhereInfrastructure(sitewhereConfig)
	if err != nil {
		fmt.Printf("Error Installing SiteWhere Infrastucture: %v\n", err)
		return
	}

	color.Style{color.FgGreen, color.OpBold}.Println("\nSiteWhere 3.0 Installed")
}
