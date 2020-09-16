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

package operator

import (
	"fmt"
)

// Template for generating a Operator Filename
const operatorFileTemplate = "/operator/operator-%02d.yaml"

// Number of Infrastructure Files
const operatorFileNumber = 23

var operatorFiles []string

func init() {
	for i := 1; i <= operatorFileNumber; i++ {
		var fileName = fmt.Sprintf(operatorFileTemplate, i)
		operatorFiles = append(operatorFiles, fileName)
	}
}

// GetSiteWhereOperatorFiles returns the name of SiteWhere Operator files
func GetSiteWhereOperatorFiles() []string {
	return operatorFiles
}
