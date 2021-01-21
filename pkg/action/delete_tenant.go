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
	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	"github.com/sitewhere/swctl/pkg/tenant"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	"helm.sh/helm/v3/pkg/action"
)

// DeleteTenant is the action for deleting a SiteWhere tenant
type DeleteTenant struct {
	cfg *action.Configuration
	// Name of the instance
	InstanceName string
	// Name of the tenant
	TenantName string
}

// NewDeleteTenant constructs a new *Install
func NewDeleteTenant(cfg *action.Configuration) *DeleteTenant {
	return &DeleteTenant{
		cfg:          cfg,
		InstanceName: "",
		TenantName:   "",
	}
}

// Run executes the list command, returning a set of matches.
func (i *DeleteTenant) Run() (*tenant.CreateSiteWhereTenant, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}

	// Check if the instance exists
	client, err := ControllerClient(i.cfg)
	if err != nil {
		return nil, err
	}

	var swTenantCR sitewhereiov1alpha4.SiteWhereTenant
	err = client.Get(context.TODO(), k8sClient.ObjectKey{Namespace: i.InstanceName, Name: i.TenantName}, &swTenantCR)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()

	if err := client.Delete(ctx, &swTenantCR); err != nil {
		if apierrors.IsNotFound(err) {
			i.cfg.Log(fmt.Sprintf("Tenant %s is not present. Skipping.", swTenantCR.GetName()))
		} else {
			return nil, err
		}
	}

	return &tenant.CreateSiteWhereTenant{
		InstanceName: i.InstanceName,
		TenantName:   i.TenantName,
	}, nil
}

// ExtractTenantName returns the name of the instance that should be used.
func (i *DeleteTenant) ExtractTenantName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}
