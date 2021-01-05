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

package uninstall

import "github.com/sitewhere/swctl/pkg/status"

// SiteWhereUninstall destribe the uninstallation of SiteWhere.
type SiteWhereUninstall struct {
	// Status of SiteWhere CDR installation
	CDRStatuses []status.SiteWhereStatus `json:"crdStatues,omitempty"`
	// Status of SiteWhere Templates installation
	TemplatesStatues []status.SiteWhereStatus `json:"templatesStatues,omitempty"`
	// Status of SiteWhere Operator
	OperatorStatuses []status.SiteWhereStatus `json:"operatorStatuses,omitempty"`
	// Status of SiteWhere Infrastructure
	InfrastructureStatuses []status.SiteWhereStatus `json:"infrastructureStatuses,omitempty"`
}
