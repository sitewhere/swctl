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
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
)

// ErrNotFound is the error when the configuration does not exists
var ErrNotFound = errors.New("not found")

// Configuration is the configuratio of SiteWhere Control CLI
type Configuration struct {
	// Microservices are the definitions of the microservices
	Microservices []sitewhereiov1alpha4.SiteWhereMicroserviceSpec `json:"microservices,omitempty"`
}

// PlaceHolder are the values sent to the template for replacement
type PlaceHolder struct {
	// Name of the instance
	InstanceName string
	// Number of replicas
	Replicas int32
	// Registry is the docker registry of the microservices images
	Registry string
	// Repository
	Repository string
	// Docker image tag
	Tag string
}

// GetConfigPath returns the path for SiteWhere Control CLI configuration path.
func GetConfigPath() string {
	return filepath.FromSlash(GetConfigHome() + "/default.yaml")
}

// GetConfigHome returns the home directory for the configuration
func GetConfigHome() string {
	home, _ := homedir.Dir()
	return filepath.FromSlash(home + "/.swctl")

}
