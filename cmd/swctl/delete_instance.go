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

package main

import (
	"io"

	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/instance"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
)

var deleteInstanceDesc = `
Use this command to delete a SiteWhere Instance. 
Use can use purge flag to remove the namespace of the instance.
`

func newDeleteInstanceCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewDeleteInstance(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:               "instance [NAME]",
		Short:             "delete an instance",
		Long:              deleteInstanceDesc,
		Args:              require.ExactArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(_ *cobra.Command, args []string) error {
			instanceName, err := client.ExtractInstanceName(args)
			if err != nil {
				return err
			}
			client.InstanceName = instanceName
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outFmt.Write(out, newDeleteInstanceWriter(results))
		},
	}

	addDeleteInstanceFlags(cmd, cmd.Flags(), client)
	bindOutputFlag(cmd, &outFmt)

	return cmd
}

func addDeleteInstanceFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.DeleteInstance) {
	f.BoolVarP(&client.Purge, "purge", "p", client.Purge, "Purge instance.")
}

type deleteInstancePrinter struct {
	instance *instance.DeleteSiteWhereInstance
}

func newDeleteInstanceWriter(result *instance.DeleteSiteWhereInstance) *deleteInstancePrinter {
	return &deleteInstancePrinter{instance: result}
}

func (s deleteInstancePrinter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, s.instance)
}

func (s deleteInstancePrinter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, s.instance)
}

func (s deleteInstancePrinter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("INSTANCE", "STATUS")
	table.AddRow(s.instance.InstanceName, color.Info.Render("Deleted"))
	return output.EncodeTable(out, table)
}
