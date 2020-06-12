/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// instancesCmd represents the instances command
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage SiteWhere Instance",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "sitewhere"
		}
		handleInstance(name)
	},
}

func init() {
	rootCmd.AddCommand(instancesCmd)
	instancesCmd.Flags().StringP("name", "n", viper.GetString("SW_INSTACE_NAME"), "SiteWhere Instance Name")
}

func handleInstance(instanceName string) {
	fmt.Println("Handling Instance: ", instanceName)
}
