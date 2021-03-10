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
)

// Level is the debug level of a microservice
type Level string

const (
	// NoLevel is the abscense of an logging level
	NoLevel Level = ""
	// DebugLevel is the debug level
	DebugLevel Level = "debug"
	// InfoLevel is the info level
	InfoLevel Level = "info"
	// WarnLevel is the warn level
	WarnLevel Level = "warn"
	// ErrorLevel is the error level
	ErrorLevel Level = "error"
	// FatalLevel is the fatal level
	FatalLevel Level = "fatal"
	// OffLevel is the off level
	OffLevel Level = "off"
)

// Parse parse a string to a Level
func Parse(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "off":
		return OffLevel, nil
	}
	return NoLevel, fmt.Errorf("no logging level match for %s", level)
}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case OffLevel:
		return "off"
	}
	return ""
}
