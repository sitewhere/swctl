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

package action

import (
	"net/http"

	"github.com/rakyll/statik/fs"

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/sitewhere/swctl/pkg/check"
)

// CheckInstall is the action for check SiteWhere installation
type CheckInstall struct {
	cfg *Configuration

	StatikFS http.FileSystem
	// Use verbose mode
	Verbose bool
}

// NewCheckInstall constructs a new *Install
func NewCheckInstall(cfg *Configuration) *CheckInstall {
	statikFS, _ := fs.New()
	return &CheckInstall{
		cfg:      cfg,
		StatikFS: statikFS,
		Verbose:  false,
	}
}

// Run executes the list command, returning a set of matches.
func (i *CheckInstall) Run() (*check.SiteWhereInstall, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	clientSet, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	extensionsClients, err := i.cfg.KubernetesAPIExtensionClientSet()
	if err != nil {
		return nil, err
	}
	config, err := i.cfg.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	// Install Custom Resource Definitions
	_, err = check.SiteWhereCRDs(i.StatikFS, clientSet, extensionsClients, config)
	if err != nil {
		return nil, err
	}

	return &check.SiteWhereInstall{
		// CDRStatues: crdStatues,
	}, nil
}
