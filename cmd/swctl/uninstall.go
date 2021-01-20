/**
 * Copyright © 2014-2021 The SiteWhere Authors
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

package main

import (
	"io"

	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/install"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
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

func newUninstallCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewUninstall(cfg, settings)
	var outFmt output.Format

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
			return outFmt.Write(out, newUninstallWriter(results))
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&client.Purge, "purge", "p", false, "Purge data.")

	bindOutputFlag(cmd, &outFmt)

	return cmd
}

type uninstallWriter struct {
	Results *install.SiteWhereInstall `json:"results"`
}

func newUninstallWriter(results *install.SiteWhereInstall) *uninstallWriter {
	return &uninstallWriter{Results: results}
}

func (i *uninstallWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow(color.Style{color.FgGreen, color.OpBold}.Render("SiteWhere 3.0 Uninstalled"))
	return output.EncodeTable(out, table)
}

func (i *uninstallWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *uninstallWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
