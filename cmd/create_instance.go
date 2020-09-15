/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sitewhere/swctl/cmd/require"
	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/cli/output"
	"github.com/sitewhere/swctl/pkg/instance"
)

var createInstanceDesc = `
Use this command to create an Instance of SiteWhere.
For example, to create an instance with name "sitewhere" use:

  swctl create instance sitewhere

To create an instance with the minimal profile use:

	swctl create instance sitewhere -m
`

func newCreateInstanceCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewCreateInstance(cfg)
	var outfmt output.Format

	cmd := &cobra.Command{
		Use:               "instance [NAME]",
		Short:             "create an instance",
		Long:              createInstanceDesc,
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
			return outfmt.Write(out, newCreateInstanceWriter(results))
		},
	}

	addCreateInstanceFlags(cmd, cmd.Flags(), client)
	bindOutputFlag(cmd, &outfmt)

	return cmd
}

func addCreateInstanceFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.CreateInstance) {
	f.StringVarP(&client.Namespace, "namespace", "n", client.Namespace, "Namespace of the instance.")
	f.BoolVarP(&client.Minimal, "minimal", "m", client.Minimal, "Minimal installation.")
	f.StringVarP(&client.Tag, "tag", "t", client.Tag, "Docker image tag.")
	f.BoolVarP(&client.Debug, "debug", "d", client.Debug, "Debug mode.")
	f.StringVarP(&client.ConfigurationTemplate, "config-template", "c", client.ConfigurationTemplate, "Configuration template.")
	f.StringVarP(&client.DatasetTemplate, "dateset-template", "x", client.DatasetTemplate, "Dataset template.")
}

type createInstancePrinter struct {
	instance *instance.CreateSiteWhereInstance
}

func newCreateInstanceWriter(result *instance.CreateSiteWhereInstance) *createInstancePrinter {
	return &createInstancePrinter{instance: result}
}

func (s createInstancePrinter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, s.instance)
}

func (s createInstancePrinter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, s.instance)
}

func (s createInstancePrinter) WriteTable(out io.Writer) error {
	return nil
}
