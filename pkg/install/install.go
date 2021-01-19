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

package install

// Status Status of a installable item.
type Status string

const (
	// Installed The item is installed.
	Installed Status = "Installed"
	// Uninstalled The item is not longer installed.
	Uninstalled = "Uninstalled"
	// Unknown We cannot know if the item is installed or not.
	Unknown = "Unknown"
)

// SiteWhereInstall destribe the installation of SiteWhere.
type SiteWhereInstall struct {
	// Release
	Release string `json:"release,omitempty"`
	// Namespace
	Namespace string `json:"namespace,omitempty"`
}
