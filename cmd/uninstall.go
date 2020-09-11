/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"io"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/cmd/require"
	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/cli/output"
	"github.com/sitewhere/swctl/pkg/uninstall"
)

var uninstallHelp = `
Uninstall SiteWhere from your Kubernetes Cluster.
This command will uninstall:
 - SiteWhere System Namespace: sitewhere-system (default)
 - SiteWhere Custom Resources Definition.
 - SiteWhere Templates.
 - SiteWhere Operator.
 - SiteWhere Infrastructure.
`

func newUninstallCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewUninstall(cfg)
	var outfmt output.Format

	cmd := &cobra.Command{
		Use:               "uninstall",
		Short:             "uninstall SiteWhere CRD and Operators",
		Long:              uninstallHelp,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outfmt.Write(out, newUninstallWriter(results))
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&client.Minimal, "minimal", "m", false, "Minimal uninstallation.")
	f.BoolVarP(&client.Verbose, "verbose", "v", false, "Verbose uninstallation.")
	f.BoolVarP(&client.Purge, "purge", "p", false, "Purge data.")

	bindOutputFlag(cmd, &outfmt)

	return cmd
}

type uninstallWriter struct {
}

func newUninstallWriter(uninstall *uninstall.SiteWhereUninstall) *uninstallWriter {
	return nil
}

func (i *uninstallWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE", "REVISION", "UPDATED", "STATUS", "CHART", "APP VERSION")
	// for _, r := range r.releases {
	// 	table.AddRow(r.Name, r.Namespace, r.Revision, r.Updated, r.Status, r.Chart, r.AppVersion)
	// }
	return output.EncodeTable(out, table)
}

func (i *uninstallWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *uninstallWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
