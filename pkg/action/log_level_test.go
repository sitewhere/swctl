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

package action

import (
	"strings"
	"testing"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	"github.com/sitewhere/swctl/pkg/logs"

	"helm.sh/helm/v3/pkg/action"
)

func TestNewLogLeves(t *testing.T) {
	var newLogLevels = NewLogLevel(&action.Configuration{})
	if newLogLevels == nil {
		t.Fatalf("should have returned an action")
	}
}

func TestGenerateLoggingOverrides(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		action   *LogLevel
		input    []sitewhereiov1alpha4.MicroserviceLoggingEntry
		expected []sitewhereiov1alpha4.MicroserviceLoggingEntry
		err      error
	}{
		{
			name: "sigle",
			action: &LogLevel{
				Level: logs.DebugLevel,
			},
			input: []sitewhereiov1alpha4.MicroserviceLoggingEntry{
				{
					Logger: "somothing.com",
					Level:  "info",
				},
			},
			expected: []sitewhereiov1alpha4.MicroserviceLoggingEntry{
				{
					Logger: "somothing.com",
					Level:  "debug",
				},
			},
		},
		{
			name: "sigle",
			action: &LogLevel{
				Level:  logs.DebugLevel,
				Logger: []string{"somothing.com"},
			},
			input: []sitewhereiov1alpha4.MicroserviceLoggingEntry{
				{
					Logger: "somothing.com",
					Level:  "info",
				},
				{
					Logger: "notchanged.com",
					Level:  "info",
				},
			},
			expected: []sitewhereiov1alpha4.MicroserviceLoggingEntry{
				{
					Logger: "somothing.com",
					Level:  "debug",
				},
				{
					Logger: "notchanged.com",
					Level:  "info",
				},
			},
		},
	}
	for _, single := range data {
		t.Run(single.name, func(single struct {
			name     string
			action   *LogLevel
			input    []sitewhereiov1alpha4.MicroserviceLoggingEntry
			expected []sitewhereiov1alpha4.MicroserviceLoggingEntry
			err      error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := single.action.generateLoggingOverrides(single.input)
				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if len(single.expected) != len(result) {
						t.Fatalf("expected %d overrides, got %d overrides", len(single.expected), len(result))
					}
					for i, r := range result {
						e := single.expected[i]
						if e.Logger != r.Logger {
							t.Fatalf("expected logger: %s got logger: %s", e.Logger, r.Logger)
						}
						if e.Level != r.Level {
							t.Fatalf("expected level: %s got level: %s", e.Logger, r.Logger)
						}
					}
				}
			}
		}(single))
	}
}
