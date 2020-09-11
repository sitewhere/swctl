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

	"github.com/sitewhere/swctl/cmd/require"
	"github.com/sitewhere/swctl/pkg/action"
)

var deleteHelp = `
Delete a SiteWhere resource from a file or from stdin.

You can delete a SiteWhere instance by using:
  - swctl delete instance sitewhere
`

func newDeleteCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete",
		Short:             "delete a SiteWhere resource from a file or from stdin.",
		Long:              deleteHelp,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions, // Disable file completion
	}

	cmd.AddCommand(newDeleteInstanceCmd(cfg, out))

	return cmd
}
