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

	"github.com/sitewhere/swctl/cmd/swctl/require"
	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/tenant"

	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
)

var createTenantDesc = `
Use this command to create a Tenant of SiteWhere.
For example, to create a tenant with name "sitewhereTenant" use:
 
swctl create tenant sitewhereTenant

`

func newCreateTenantCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewCreateTenant(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:               "tenant [NAME]",
		Short:             "create an tenant",
		Long:              createTenantDesc,
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
			return outFmt.Write(out, newCreateTenantWriter(results))
		},
	}

	addCreateTenantFlags(cmd, cmd.Flags(), client)
	bindOutputFlag(cmd, &outFmt)

	return cmd
}

type createTenantPrinter struct {
	instance *tenant.CreateSiteWhereTenant
}

func newCreateTenantWriter(result *tenant.CreateSiteWhereTenant) *createTenantPrinter {
	return &createTenantPrinter{instance: result}
}

func (s createTenantPrinter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, s.instance)
}

func (s createTenantPrinter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, s.instance)
}

func (s createTenantPrinter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("INSTANCE", "TENANT", "STATUS")
	table.AddRow(s.instance.InstanceName, s.instance.TenantName, color.Info.Render("Installed"))
	return output.EncodeTable(out, table)
}

func addCreateTenantFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.CreateTenant) {
	f.StringVarP(&client.InstanceName, "instance", "i", client.InstanceName, "Instance name")
	f.StringSliceVarP(&client.AuthorizedUserIds, "authorizedUserIds", "u", client.AuthorizedUserIds, "Authorized User Ids")
	f.StringVarP(&client.AuthenticationToken, "authenticationToken", "t", client.AuthenticationToken, "AuthenticationToken")
	f.StringVarP(&client.ConfigurationTemplate, "configurationTemplate", "c", client.ConfigurationTemplate, "Configuration Template")
	f.StringVarP(&client.DatasetTemplate, "datasetTemplate", "d", client.DatasetTemplate, "Dataset Template")

	cmd.MarkFlagRequired("instance")
}
