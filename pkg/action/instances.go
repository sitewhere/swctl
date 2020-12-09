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
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/sitewhere/swctl/pkg/instance"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
)

// Instances is the action for listing SiteWhere instances
type Instances struct {
	cfg *Configuration
}

// NewInstances constructs a new *Instances
func NewInstances(cfg *Configuration) *Instances {
	return &Instances{
		cfg: cfg,
	}
}

// Run executes the install command, returning the result of the installation
func (i *Instances) Run() (*instance.ListSiteWhereInstance, error) {
	var client, err = i.cfg.ControllerClient()
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	var swInstancesList sitewhereiov1alpha4.SiteWhereInstanceList

	if err := client.List(ctx, &swInstancesList); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}
	}
	return &instance.ListSiteWhereInstance{
		Instances: swInstancesList.Items,
	}, nil
}
