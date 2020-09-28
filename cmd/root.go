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
		newInstancesCmd(actionConfig, out),
		newUninstallCmd(actionConfig, out),
		newVersionCmd(out))

	return cmd, nil
}
