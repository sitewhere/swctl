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
	"fmt"
	"strings"

	"github.com/pkg/errors"
	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	"github.com/sitewhere/swctl/pkg/tenant"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateTenant is the action for creating a SiteWhere tenant
type CreateTenant struct {
	cfg *Configuration
	// Name of the instance
	InstanceName string
	// Name of the tenant
	TenantName string
	// AuthenticationToken is the token used for authenticating the tenant
	AuthenticationToken string
	// Authorized are the IDs of the users that are authorized to use the tenant
	AuthorizedUserIds []string
	// ConfigurationTemplate is the configuration template used for the tenant
	ConfigurationTemplate string
	// DatasetTemplate is the dataset template used for the tenant
	DatasetTemplate string
}

type tenantResourcesResult struct {
	// Custom Resource Name of the instance
	TenantName string
}

// NewCreateTenant constructs a new *Install
func NewCreateTenant(cfg *Configuration) *CreateTenant {
	return &CreateTenant{
		cfg:                   cfg,
		InstanceName:          "",
		TenantName:            "",
		ConfigurationTemplate: "default",
		DatasetTemplate:       "construction",
	}
}

// Run executes the list command, returning a set of matches.
func (i *CreateTenant) Run() (*tenant.CreateSiteWhereTenant, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}

	//Revisar si existe la instancia
	client, err := i.cfg.ControllerClient()
	if err != nil {
		return nil, err
	}

	var swInstance sitewhereiov1alpha4.SiteWhereInstance
	err = client.Get(context.TODO(), k8sClient.ObjectKey{Name: i.InstanceName}, &swInstance)
	if err != nil {
		return nil, err
	}

	swTenantCR := i.buildCRSiteWhereTenant()
	ctx := context.TODO()

	if err := client.Create(ctx, swTenantCR); err != nil {
		if apierrors.IsAlreadyExists(err) {
			i.cfg.Log(fmt.Sprintf("Tenant %s is already present. Skipping.", swTenantCR.GetName()))
		} else {
			return nil, err
		}
	}

	return &tenant.CreateSiteWhereTenant{
		InstanceName: i.InstanceName,
		TenantName:   i.TenantName,
	}, nil
}

func (i *CreateTenant) buildCRSiteWhereTenant() *sitewhereiov1alpha4.SiteWhereTenant {
	return &sitewhereiov1alpha4.SiteWhereTenant{
		TypeMeta: metav1.TypeMeta{
			Kind:       sitewhereiov1alpha4.SiteWhereTenantKind,
			APIVersion: sitewhereiov1alpha4.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.TenantName,
			Namespace: i.InstanceName,
		},
		Spec: sitewhereiov1alpha4.SiteWhereTenantSpec{
			Name:                  i.TenantName,
			AuthenticationToken:   i.AuthenticationToken,
			AuthorizedUserIds:     i.AuthorizedUserIds,
			DatasetTemplate:       i.DatasetTemplate,
			ConfigurationTemplate: i.ConfigurationTemplate,
		},
	}
}

// ExtractTenantName returns the name of the instance that should be used.
func (i *CreateTenant) ExtractTenantName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}
