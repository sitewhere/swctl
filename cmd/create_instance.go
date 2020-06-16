/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// createInstanceCmd represents the instance command
var (
	namespace         = ""
	createInstanceCmd = &cobra.Command{
		Use:   "instance",
		Short: "Create SiteWhere Instance",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires one argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			createSiteWhereInstance(name)
		},
	}
)

func init() {
	createInstanceCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the instance.")
	createCmd.AddCommand(createInstanceCmd)
}

func createSiteWhereInstance(name string) {

}
