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
	"github.com/sitewhere/swctl/pkg/instance"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
)

var instancesHelp = `
Use this command to list SiteWhere Intances.
`

func newInstancesCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewInstances(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:               "instances [NAME]",
		Short:             "show SiteWhere instances",
		Long:              instancesHelp,
		Args:              require.MaximumNArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {

			instanceName, err := client.ExtractInstanceNameArg(args)
			client.InstanceName = instanceName

			results, err := client.Run()
			if err != nil {
				return err
			}
			return outFmt.Write(out, newInstancesWriter(results))
		},
	}
	bindOutputFlag(cmd, &outFmt)
	return cmd
}

type instancesWriter struct {
	// Instances found
	Instances []sitewhereiov1alpha4.SiteWhereInstance

	//Microservices found
	SiteWhereMicroservice []sitewhereiov1alpha4.SiteWhereMicroservice
}

func newInstancesWriter(result *instance.ListSiteWhereInstance) *instancesWriter {
	return &instancesWriter{
		Instances:             result.Instances,
		SiteWhereMicroservice: result.SiteWhereMicroservice,
	}
}

func (i *instancesWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE", "CONFIG TMPL", "DATESET TMPL", "TM STATUS", "UM STATUS")
	for _, item := range i.Instances {
		tmState := renderState(item.Status.TenantManagementBootstrapState)
		umStatus := renderState(item.Status.UserManagementBootstrapState)
		table.AddRow(item.Name, item.Name, item.Spec.ConfigurationTemplate, item.Spec.DatasetTemplate, tmState, umStatus)
	}

	table.AddRow("", "", "", "", "")
	output.EncodeTable(out, table)

	if len(i.Instances) == 1 && len(i.SiteWhereMicroservice) > 0 {
		i.WriteMicroserviceInfo(out)
		i.WriteInstanceDetailInfo(out, i.Instances[0])
	}
	return nil
}

func (i *instancesWriter) WriteMicroserviceInfo(out io.Writer) {
	microserviceTable := uitable.New()
	microserviceTable.AddRow("MICROSERVICE", "NAMESPACE", "DEPLOYMENT")
	for _, item := range i.SiteWhereMicroservice {
		microserviceTable.AddRow(item.Spec.Name, item.ObjectMeta.Namespace, item.Status.Deployment)
	}
	microserviceTable.AddRow("", "", "")
	output.EncodeTable(out, microserviceTable)
	output.EncodeYAML(out, i.Instances[0].Spec.DockerSpec)
	output.EncodeYAML(out, i.Instances[0].Spec.Configuration)

}

func (i *instancesWriter) WriteInstanceDetailInfo(out io.Writer, instance sitewhereiov1alpha4.SiteWhereInstance) {
	output.EncodeYAML(out, instance.Spec.DockerSpec)
	output.EncodeYAML(out, instance.Spec.Configuration)
}

func (i *instancesWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *instancesWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}

func renderState(state sitewhereiov1alpha4.BootstrapState) string {
	switch state {
	case "Unknown":
		return color.Warn.Render("Unknown")
	case "Bootstrapped":
		return color.Info.Render("Bootstrapped")
	case "NotBootstrapped":
		return color.Error.Render("Not Bootstrapped")
	default:
		return ""
	}
}
