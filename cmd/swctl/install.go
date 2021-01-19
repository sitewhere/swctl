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

var installHelp = `
Use this command to install SiteWhere 3.0 on a Kubernetes Cluster.
This command will install:
 - SiteWhere System Namespace: sitewhere-system (default)
 - SiteWhere Custom Resources Definition.
 - SiteWhere Templates.
 - SiteWhere Operator.
 - SiteWhere Infrastructure.
`

func newInstallCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewInstall(cfg, settings)
	var outFmt output.Format

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
			return outFmt.Write(out, newInstallWriter(client.SkipCRD, client.SkipTemplate, client.SkipOperator, client.SkipInfrastructure, results))
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&client.WaitReady, "wait", "w", false, "Wait for components to be ready before return control.")
	f.BoolVar(&client.SkipCRD, "skip-crd", false, "Skip Custom Resource Definition installation.")
	f.BoolVar(&client.SkipTemplate, "skip-templates", false, "Skip Templates installation.")
	f.BoolVar(&client.SkipOperator, "skip-operator", false, "Skip Operator installation.")
	f.BoolVar(&client.SkipInfrastructure, "skip-infra", false, "Skip Infrastructure installation.")
	bindOutputFlag(cmd, &outFmt)

	return cmd
}

type installWriter struct {
	SkipCRD            bool
	SkipTemplate       bool
	SkipOperator       bool
	SkipInfrastructure bool
	Results            *install.SiteWhereInstall `json:"results"`
}

func newInstallWriter(skipCRD bool, skipTemplate bool, skipOperator bool, skipInfrastructure bool, results *install.SiteWhereInstall) *installWriter {
	return &installWriter{
		SkipCRD:            skipCRD,
		SkipTemplate:       skipTemplate,
		SkipOperator:       skipOperator,
		SkipInfrastructure: skipInfrastructure,
		Results:            results,
	}
}

func (i *installWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("COMPONENT", "STATUS")
	if !i.SkipCRD {
		table.AddRow("Custom Resource Definitions", color.Info.Render("Installed"))
	}
	if !i.SkipTemplate {
		table.AddRow("Templates", color.Info.Render("Installed"))
	}
	if !i.SkipOperator {
		table.AddRow("Operator", color.Info.Render("Installed"))
	}
	if !i.SkipInfrastructure {
		table.AddRow("Infrastructure", color.Info.Render("Installed"))
	}
	table.AddRow(color.Style{color.FgGreen, color.OpBold}.Render("SiteWhere 3.0 Installed"))
	return output.EncodeTable(out, table)
}

func (i *installWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *installWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
