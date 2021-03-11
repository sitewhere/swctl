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

	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/pkg/action"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
)

var logsHelp = `
Use this command to show the logs of a SiteWhere Microservice.
`

func newLogsCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewLogs(cfg)

	cmd := &cobra.Command{
		Use:   "logs [OPTIONS] INSTANCE MS",
		Short: "show the logs of a SiteWhere Microservice",
		Long:  logsHelp,
		Args:  require.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return compListInstances(toComplete, cfg)
			} else if len(args) == 1 {
				return compListMicroservices(toComplete, args[0], cfg)
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client.InstanceName = args[0]
			client.MicroserviceName = args[1]
			return client.Run()
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&client.Follow, "follow", "f", false, "Specify if the logs should be streamed.")
	return cmd
}
