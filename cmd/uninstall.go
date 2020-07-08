/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"fmt"

	"github.com/rakyll/statik/fs"

	"github.com/sitewhere/swctl/internal"
	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var (
	purge        = false // Purge data
	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall SiteWhere from your Kubernetes Cluster",
		Long: `Uninstall SiteWhere from your Kubernetes Cluster.
This command will uninstall:
 - SiteWhere System Namespace: sitewhere-system (default)
 - SiteWhere Custom Resources Definition.
 - SiteWhere Templates.
 - SiteWhere Operator.
 - SiteWhere Infrastructure.`,
		Run: uninstallSiteWhereCommand,
	}
)

func init() {
	uninstallCmd.Flags().BoolVarP(&purge, "purge", "p", false, "Purge data.")
	rootCmd.AddCommand(uninstallCmd)
}

// uninstallSiteWhereCommand Performs the steps necessary to uninstall SiteWhere
func uninstallSiteWhereCommand(_ *cobra.Command, _ []string) {
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

	// Uninstall Infrastructure
	err = internal.UninstallSiteWhereInfrastructure(config, statikFS)
	if err != nil {
		fmt.Printf("Error Uninstalling SiteWhere Infrastucture: %v\n", err)
		return
	}

	// Uninstall Operator
	err = internal.UninstallSiteWhereOperator(config, statikFS)
	if err != nil {
		fmt.Printf("Error Uninstalling SiteWhere Operator: %v\n", err)
		return
	}

	// Uninstall Custom Resource Definitions
	internal.UninstallSiteWhereCRDs(config, statikFS)

	fmt.Println("SiteWhere 3.0 uninstalled.")
}
