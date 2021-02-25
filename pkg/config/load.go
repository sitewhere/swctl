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
	"io/ioutil"
	"os"
)

// LoadConfigurationTemplate loads the configuration template from
// ~/swctl/deafult.yaml file. If the files does not exist
// it returns the error ErrNotFound
func LoadConfigurationTemplate(placeHolder *PlaceHolder) (string, error) {
	configPath := GetConfigPath()

	f, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotFound
		}
		return "", err
	}
	defer f.Close()

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// LoadConfigurationOrDefault loads the configuration from
// ~/swctl/config file or load the default configuration
func LoadConfigurationOrDefault(placeHolder *PlaceHolder) (*Configuration, error) {
	templateContext, err := LoadConfigurationTemplate(placeHolder)
	if err != nil {
		templateContext = defaultTemplate
	}
	return FromTemplate(templateContext, placeHolder)
}
