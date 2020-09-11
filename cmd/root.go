/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/pkg/action"
)

var globalUsage = `SiteWhere Control allow you to manage SiteWhere CE Instances.`

// NewRootCmd creates a new root command.
func NewRootCmd(actionConfig *action.Configuration, out io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "swctl",
		Short:        "SiteWhere Control CLI.",
		Long:         globalUsage,
		SilenceUsage: true,
		// This breaks completion for 'helm help <TAB>'
		// The Cobra release following 1.0 will fix this
		//ValidArgsFunction: noCompletions, // Disable file completion
	}
	flags := cmd.PersistentFlags()

	// Command completion
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.Parse(args)

	// Add subcommands
	cmd.AddCommand(
		newInstallCmd(actionConfig, out),
		newCheckInstallCmd(actionConfig, out),
		newCreateCmd(actionConfig, out),
		newDeleteCmd(actionConfig, out),
		newUninstallCmd(actionConfig, out),
		newVersionCmd(out))

	return cmd, nil
}
