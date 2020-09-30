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

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	// "os/exec"
	// "runtime"
	"strings"
	"testing"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/internal/test"
	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/cli"
	kubefake "github.com/sitewhere/swctl/pkg/kube/fake"
)

func runTestCmd(t *testing.T, tests []cmdTestCase) {
	t.Helper()
	for _, tt := range tests {
		for i := 0; i <= tt.repeat; i++ {
			t.Run(tt.name, func(t *testing.T) {
				defer resetEnv()()

				t.Logf("running cmd (attempt %d): %s", i+1, tt.cmd)
				_, out, err := executeActionCommandC(tt.cmd)
				if (err != nil) != tt.wantError {
					t.Errorf("expected error, got '%v'", err)
				}
				if tt.golden != "" {
					test.AssertGoldenString(t, out, tt.golden)
				}
			})
		}
	}
}

func executeActionCommandC(cmd string) (*cobra.Command, string, error) {
	return executeActionCommandStdinC(nil, cmd)
}

func executeActionCommandStdinC(in *os.File, cmd string) (*cobra.Command, string, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)

	actionConfig := &action.Configuration{
		KubeClient: &kubefake.PrintingKubeClient{Out: ioutil.Discard},
		Log:        func(format string, v ...interface{}) {},
	}

	root, err := newRootCmd(actionConfig, buf, args)
	if err != nil {
		return nil, "", err
	}

	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	oldStdin := os.Stdin
	if in != nil {
		root.SetIn(in)
		os.Stdin = in
	}

	c, err := root.ExecuteC()

	result := buf.String()

	os.Stdin = oldStdin

	return c, result, err
}

// cmdTestCase describes a test case that.
type cmdTestCase struct {
	name      string
	cmd       string
	golden    string
	wantError bool
	// Number of repeats (in case a feature was previously flaky and the test checks
	// it's now stably producing identical results). 0 means test is run exactly once.
	repeat int
}

func resetEnv() func() {
	origEnv := os.Environ()
	return func() {
		os.Clearenv()
		for _, pair := range origEnv {
			kv := strings.SplitN(pair, "=", 2)
			os.Setenv(kv[0], kv[1])
		}
		settings = cli.New()
	}
}
