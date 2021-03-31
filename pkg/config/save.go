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
	"os"
)

// CreateDefaultConfiguration saves the default configuration to the config file
func CreateDefaultConfiguration() error {
	return createConfiguration(GetConfigPath(), defaultTemplate)
}

// CreateMinimalConfiguration saves the minimal configuration to the config file
func CreateMinimalConfiguration() error {
	return createConfiguration(GetMinimalConfigPath(), minimalTemplate)
}

// CreateDebugConfiguration saves the debug configuration to the config file
func CreateDebugConfiguration() error {
	return createConfiguration(GetDebugConfigPath(), debugTemplate)
}

// createConfiguration saves the default configuration to the config file
func createConfiguration(configPath string, templete string) error {
	var err error
	configHome := GetConfigHome()
	err = os.Mkdir(configHome, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(templete))
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
