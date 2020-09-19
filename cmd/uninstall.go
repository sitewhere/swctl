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

	"github.com/gookit/color"
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
	f.BoolVarP(&client.Purge, "purge", "p", false, "Purge data.")

	bindOutputFlag(cmd, &outfmt)

	return cmd
}

type uninstallWriter struct {
	Results *uninstall.SiteWhereUninstall `json:"results"`
}

func newUninstallWriter(results *uninstall.SiteWhereUninstall) *uninstallWriter {
	return &uninstallWriter{Results: results}
}

func (i *uninstallWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("COMPONENT", "STATUS")
	table.AddRow("Custom Resource Definitions", color.Info.Render("Uninstalled"))
	table.AddRow("Templates", color.Info.Render("Uninstalled"))
	table.AddRow("Operator", color.Info.Render("Uninstalled"))
	table.AddRow("Infrastructure", color.Info.Render("Uninstalled"))
	table.AddRow(color.Style{color.FgGreen, color.OpBold}.Render("SiteWhere 3.0 Uninstalled"))
	return output.EncodeTable(out, table)
}

func (i *uninstallWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *uninstallWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
