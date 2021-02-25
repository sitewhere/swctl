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

package config

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v2"
)

// FromTemplate renders the configuration from a template
func FromTemplate(templateContent string, placeHolder *PlaceHolder) (*Configuration, error) {
	tmpl, err := template.New(placeHolder.InstanceName).Parse(templateContent)
	if err != nil {
		return nil, err
	}
	var templatedBuffer bytes.Buffer
	err = tmpl.Execute(&templatedBuffer, placeHolder)
	if err != nil {
		return nil, err
	}
	var cfg Configuration
	err = yaml.Unmarshal(templatedBuffer.Bytes(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
