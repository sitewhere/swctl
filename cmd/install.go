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
	"github.com/sitewhere/swctl/pkg/install"
)

var installHelp = `
Use this command to install SiteWhere 3.0 on a Kubernetes Cluster.
This command will install:
 - SiteWhere System Namespace: sitewhere-system (default)
 - SiteWhere Custom Resources Definition.
 - SiteWhere Templates.
 - SiteWhere Operator.
 - SiteWhere Infrastructure.
`

func newInstallCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewInstall(cfg)
	var outfmt output.Format

	cmd := &cobra.Command{
		Use:               "install",
		Short:             "Install SiteWhere CRD and Operators",
		Long:              installHelp,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outfmt.Write(out, newInstallWriter(results))
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&client.Minimal, "minimal", "m", false, "Minimal installation.")
	f.BoolVarP(&client.WaitReady, "wait", "w", false, "Wait for components to be ready before return control.")
	bindOutputFlag(cmd, &outfmt)

	return cmd
}

type installWriter struct {
}

func newInstallWriter(install *install.SiteWhereInstall) *installWriter {
	return nil
}

func (i *installWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE", "REVISION", "UPDATED", "STATUS", "CHART", "APP VERSION")
	// for _, r := range r.releases {
	// 	table.AddRow(r.Name, r.Namespace, r.Revision, r.Updated, r.Status, r.Chart, r.AppVersion)
	// }
	return output.EncodeTable(out, table)
}

func (i *installWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *installWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
