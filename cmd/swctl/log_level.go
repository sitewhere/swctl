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
	"github.com/sitewhere/swctl/pkg/logs"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
)

var logLevelHelp = `
Use this command to change the log levels of a SiteWhere Microservice.
`

func newLogLevelCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewLogLevel(cfg)

	cmd := &cobra.Command{
		Use:     "log-level INSTANCE MS LEVEL [OPTIONS]",
		Short:   "change the log levels of a SiteWhere Microservice",
		Aliases: []string{"ll"},
		Long:    logsHelp,
		Args:    require.ExactArgs(3),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return compListInstances(toComplete, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client.InstanceName = args[0]
			client.MicroserviceName = args[1]
			level, err := logs.Parse(args[2])
			if err != nil {
				return err
			}
			client.Level = level
			return client.Run()
		},
	}
	f := cmd.Flags()
	f.StringArrayVar(&client.Logger, "logger", []string{}, "set loggers to change")

	return cmd
}
