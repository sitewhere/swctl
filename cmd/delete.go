/**
 * Copyright © 2014-2020 The SiteWhere Authors
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

package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/cmd/require"
	"github.com/sitewhere/swctl/pkg/action"
)

var deleteHelp = `
Delete a SiteWhere resource from a file or from stdin.

You can delete a SiteWhere instance by using:
  - swctl delete instance sitewhere
`

func newDeleteCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete",
		Short:             "delete a SiteWhere resource from a file or from stdin.",
		Long:              deleteHelp,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions, // Disable file completion
	}

	cmd.AddCommand(newDeleteInstanceCmd(cfg, out))

	return cmd
}
