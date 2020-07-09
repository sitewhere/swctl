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
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List instances",
	Long:  ``,
	Run:   instanceListSiteWhereCommand,
}

func init() {
	instancesCmd.AddCommand(listCmd)
}

func instanceListSiteWhereCommand(cmd *cobra.Command, args []string) {
	fmt.Println("list called")
}
