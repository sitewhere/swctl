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
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/cmd/helm/require"
)

const completionDesc = `
Generate autocompletion scripts for SiteWhere Control CLI for the specified shell.
`
const bashCompDesc = `
Generate the autocompletion script for SiteWhere Control CLI for the bash shell.

To load completions in your current shell session:
$ source <(swctl completion bash)

To load completions for every new session, execute once:
Linux:
  $ swctl completion bash > /etc/bash_completion.d/swctl
MacOS:
  $ swctl completion bash > /usr/local/etc/bash_completion.d/swctl
`

const (
	noDescFlagName = "no-descriptions"
	noDescFlagText = "disable completion descriptions"
)

var disableCompDescriptions bool

func newCompletionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "generate autocompletion scripts for the specified shell",
		Long:  completionDesc,
		Args:  require.NoArgs,
	}

	bash := &cobra.Command{
		Use:                   "bash",
		Short:                 "generate autocompletion script for bash",
		Long:                  bashCompDesc,
		Args:                  require.NoArgs,
		DisableFlagsInUseLine: true,
		ValidArgsFunction:     noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompletionBash(out, cmd)
		},
	}

	cmd.AddCommand(bash)

	return cmd
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	err := cmd.Root().GenBashCompletion(out)

	if binary := filepath.Base(os.Args[0]); binary != "swctl" {
		renamedBinaryHook := `
# Hook the command used to generate the completion script
# to the swctl completion function to handle the case where
# the user renamed the swctl binary
if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_swctl %[1]s
else
    complete -o default -o nospace -F __start_swctl %[1]s
fi
`
		fmt.Fprintf(out, renamedBinaryHook, binary)
	}

	return err
}

// Function to disable file completion
func noCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}
