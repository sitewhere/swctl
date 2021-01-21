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
	"github.com/spf13/pflag"

	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/tenant"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
)

var deleteTenantDesc = `
Use this command to delete a Tenant of a SiteWhere Instance.
For example, to delete a tenant with name "tenant2" for instance "sitewhere" use:
 
swctl delete tenant tenant2 --instance=sitewhere

`

func newDeleteTenantCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewDeleteTenant(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:               "tenant [NAME]",
		Short:             "delete a tenant for an instance",
		Long:              deleteTenantDesc,
		Args:              require.ExactArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(_ *cobra.Command, args []string) error {
			tenantNameName, err := client.ExtractTenantName(args)
			if err != nil {
				return err
			}
			client.TenantName = tenantNameName
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outFmt.Write(out, newDeleteTenantWriter(results))
		},
	}

	addDeleteTenantFlags(cmd, cmd.Flags(), client)
	bindOutputFlag(cmd, &outFmt)

	return cmd
}

type deleteTenantPrinter struct {
	instance *tenant.CreateSiteWhereTenant
}

func newDeleteTenantWriter(result *tenant.CreateSiteWhereTenant) *deleteTenantPrinter {
	return &deleteTenantPrinter{instance: result}
}

func (s deleteTenantPrinter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, s.instance)
}

func (s deleteTenantPrinter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, s.instance)
}

func (s deleteTenantPrinter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("INSTANCE", "TENANT", "STATUS")
	table.AddRow(s.instance.InstanceName, s.instance.TenantName, color.Info.Render("Deleted"))
	return output.EncodeTable(out, table)
}

func addDeleteTenantFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.DeleteTenant) {
	f.StringVarP(&client.InstanceName, "instance", "i", client.InstanceName, "Instance name")
	cmd.MarkFlagRequired("instance")
}
