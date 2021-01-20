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
	"github.com/sitewhere/swctl/pkg/install"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
)

var checkInstallHelp = `
Use this command to check the install of SiteWhere 3.0 on a Kubernetes Cluster.
`

func newCheckInstallCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewCheckInstall(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:     "check-install",
		Short:   "Check Install SiteWhere CRD and Operators",
		Aliases: []string{"check"},
		Long:    checkInstallHelp,
		Args:    require.NoArgs,
		//ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outFmt.Write(out, newCheckInstallWriter(results))
		},
	}

	bindOutputFlag(cmd, &outFmt)

	return cmd
}

type checkInstallWriter struct {
}

func newCheckInstallWriter(install *install.SiteWhereInstall) *checkInstallWriter {
	return &checkInstallWriter{}
}

func (i *checkInstallWriter) WriteTable(out io.Writer) error {
	//return output.EncodeTable(out, table)
	return nil
}

func (i *checkInstallWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *checkInstallWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}
