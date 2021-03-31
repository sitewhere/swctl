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

// Package profile defines SiteWhere Structures for Install Profiles
package profile

// SiteWhereProfile profiles to use the application
type SiteWhereProfile string

const (
	// Default profile, use all microservices
	Default SiteWhereProfile = "Default"
	// Minimal profile, use a reduce set of microservices
	Minimal SiteWhereProfile = "Minimal"
	// Debug profile, use all microservices in debug mode
	Debug SiteWhereProfile = "Debug"
)

var All []SiteWhereProfile = []SiteWhereProfile{
	Default,
	Minimal,
	Debug,
}
