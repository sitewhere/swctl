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

package logs

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		value    string
		expected Level
		err      error
	}{
		{
			name:     "debug",
			value:    "debug",
			expected: DebugLevel,
		},
		{
			name:     "debug-capital",
			value:    "deBug",
			expected: DebugLevel,
		},
		{
			name:     "info",
			value:    "info",
			expected: InfoLevel,
		},
		{
			name:     "warn",
			value:    "warn",
			expected: WarnLevel,
		},
		{
			name:     "error",
			value:    "error",
			expected: ErrorLevel,
		},
		{
			name:     "falta",
			value:    "fatal",
			expected: FatalLevel,
		},
		{
			name:     "off",
			value:    "off",
			expected: OffLevel,
		},
		{
			name:     "no-found",
			value:    "dfdsffd",
			expected: NoLevel,
			err:      fmt.Errorf("no logging level match for %s", "dfdsffd"),
		},
	}

	for _, single := range data {
		t.Run(single.name, func(single struct {
			name     string
			value    string
			expected Level
			err      error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := Parse(single.value)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result != single.expected {
						t.Fatalf("expected %s, got %s", single.expected, result)
					}
				}
			}
		}(single))
	}
}

func TestString(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		value    Level
		expected string
	}{
		{
			name:     "debug",
			value:    DebugLevel,
			expected: "debug",
		},
		{
			name:     "info",
			value:    InfoLevel,
			expected: "info",
		},
		{
			name:     "warn",
			value:    WarnLevel,
			expected: "warn",
		},
		{
			name:     "error",
			value:    ErrorLevel,
			expected: "error",
		},
		{
			name:     "falta",
			value:    FatalLevel,
			expected: "fatal",
		},
		{
			name:     "off",
			value:    OffLevel,
			expected: "off",
		},
		{
			name:     "no-level",
			value:    NoLevel,
			expected: "",
		},
	}
	for _, single := range data {
		t.Run(single.name, func(single struct {
			name     string
			value    Level
			expected string
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result := single.value.String()
				if result != single.expected {
					t.Fatalf("expected %s namespace, got %s", single.expected, result)
				}
			}
		}(single))
	}
}
