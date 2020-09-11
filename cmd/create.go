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

var createHelp = `
Create a SiteWhere resource from a file or from stdin.

You can create a SiteWhere instance by using:
  - swctl create instance sitewhere
`

func newCreateCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Short:             "create a SiteWhere resource from a file or from stdin.",
		Long:              createHelp,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions, // Disable file completion
	}

	cmd.AddCommand(newCreateInstanceCmd(cfg, out))

	return cmd
}
