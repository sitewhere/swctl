/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"fmt"
	"io"
	"text/template"

	"github.com/sitewhere/swctl/cmd/require"
	"github.com/sitewhere/swctl/internal/version"

	"github.com/spf13/cobra"
	// goVersion "go.hein.dev/go-version"
)

const versionDesc = `
Version will output the current build information
`

type versionOptions struct {
	short    bool
	template string
}

func newVersionCmd(out io.Writer) *cobra.Command {
	o := &versionOptions{}

	cmd := &cobra.Command{
		Use:               "version",
		Short:             "print the swctl version information",
		Long:              versionDesc,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&o.short, "short", false, "print the version number")
	f.StringVar(&o.template, "template", "", "template for version string format")
	f.BoolP("client", "c", true, "display client version information")
	f.MarkHidden("client")

	return cmd
}

func (o *versionOptions) run(out io.Writer) error {
	if o.template != "" {
		tt, err := template.New("_").Parse(o.template)
		if err != nil {
			return err
		}
		return tt.Execute(out, version.Get())
	}
	fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		if len(v.GitCommit) >= 7 {
			return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
		}
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", v)
}

// // versionCmd represents the version command
// var (
// 	shortened     = false
// 	version       = ""
// 	commit        = ""
// 	date          = ""
// 	versionOutput = "json"
// 	versionCmd    = &cobra.Command{
// 		Use:   "version",
// 		Short: "Version will output the current build information",
// 		Long:  ``,
// 		Run:   printSiteWhereVersion,
// 	}
// )

// func init() {
// 	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
// 	versionCmd.Flags().StringVarP(&versionOutput, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")
// 	rootCmd.AddCommand(versionCmd)
// }

// func printSiteWhereVersion(_ *cobra.Command, _ []string) {
// 	resp := goVersion.FuncWithOutput(shortened, version, commit, date, versionOutput)
// 	fmt.Print(resp)
// 	return
// }
