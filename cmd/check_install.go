/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"io"

	// "github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/cmd/require"

	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/cli/output"
	"github.com/sitewhere/swctl/pkg/install"
)

var checkInstallHelp = `
Use this command to check the install of SiteWhere 3.0 on a Kubernetes Cluster.
`

func newCheckInstallCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewCheckInstall(cfg)
	var outfmt output.Format

	cmd := &cobra.Command{
		Use:     "check-install",
		Short:   "Check Install SiteWhere CRD and Operators",
		Aliases: []string{"check"},
		Long:    checkInstallHelp,
		Args:    require.NoArgs,
		//ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outfmt.Write(out, newCheckInstallWriter(results))
		},
	}

	bindOutputFlag(cmd, &outfmt)

	return cmd
}

type checkInstallWriter struct {
}

func newCheckInstallWriter(install *install.SiteWhereInstall) *checkInstallWriter {
	return &checkInstallWriter{}
}

func (i *checkInstallWriter) WriteTable(out io.Writer) error {
	//return output.EncodeTable(out, table)
	return nil
}

func (i *checkInstallWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *checkInstallWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
