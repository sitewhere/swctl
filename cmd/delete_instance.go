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

var deleteInstanceDesc = `
Use this command to delete a SiteWhere Instance. 
Use can use purge flag to remove the namespace of the instance.
`

func newDeleteInstanceCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewDeleteInstance(cfg)
	var outfmt output.Format

	cmd := &cobra.Command{
		Use:               "instance [NAME]",
		Short:             "delete an instance",
		Long:              deleteInstanceDesc,
		Args:              require.ExactArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(_ *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outfmt.Write(out, newDeleteInstanceWriter(results))
		},
	}

	addDeleteInstanceFlags(cmd, cmd.Flags(), client)
	bindOutputFlag(cmd, &outfmt)

	return cmd
}

func addDeleteInstanceFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.DeleteInstance) {
	f.BoolVarP(&client.Pruge, "purge", "p", false, "Purge instance.")
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
	return nil
}
