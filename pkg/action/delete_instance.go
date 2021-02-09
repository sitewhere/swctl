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

package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // GPC Auth
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Auth

	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"

	"helm.sh/helm/v3/pkg/action"
)

// DeleteInstance is the action for creating a SiteWhere instance
type DeleteInstance struct {
	cfg *action.Configuration
	// Name of the instance
	InstanceName string
	// Purge Instance data
	Purge bool
}

// NewDeleteInstance constructs a new *Install
func NewDeleteInstance(cfg *action.Configuration) *DeleteInstance {
	return &DeleteInstance{
		cfg:          cfg,
		InstanceName: "",
		Purge:        false,
	}
}

// Run executes the list command, returning a set of matches.
func (i *DeleteInstance) Run() (*instance.DeleteSiteWhereInstance, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	var client, err = ControllerClient(i.cfg)
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	var swInstance sitewhereiov1alpha4.SiteWhereInstance
	if err := client.Get(ctx, types.NamespacedName{Name: i.InstanceName}, &swInstance); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("sitewhere instance '%s' not found", i.InstanceName)
		}

		return nil, err
	}
	if err := client.Delete(ctx, &swInstance); err != nil {
		return nil, err
	}

	if i.Purge {
		clientset, err := i.cfg.KubernetesClientSet()
		if err != nil {
			return nil, err
		}
		err = resources.DeleteNamespaceIfExists(i.InstanceName, clientset)
		if err != nil {
			return nil, err
		}
	}

	return &instance.DeleteSiteWhereInstance{
		InstanceName: i.InstanceName,
		Namespace:    i.InstanceName,
	}, nil
}

// ExtractInstanceName returns the name of the instance that should be used.
func (i *DeleteInstance) ExtractInstanceName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}
