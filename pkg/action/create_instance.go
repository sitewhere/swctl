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

	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
)

// CreateInstance is the action for creating a SiteWhere instance
type CreateInstance struct {
	cfg *Configuration
	// Name of the instance
	InstanceName string
	// Name of the tenant
	TenantName string
	// Namespace to use
	Namespace string
	// Use minimal profile. Initialize only essential microservices.
	Minimal bool
	// Number of replicas
	Replicas int64
	// Docker image tag
	Tag string
	// Use debug mode
	Debug bool
	// Configuration Template
	ConfigurationTemplate string
	// Dataset template
	DatasetTemplate string
}

type namespaceAndResourcesResult struct {
	// Namespace created
	Namespace string
	// Service Account created
	ServiceAccountName string
	// Custer Role created
	ClusterRoleName string
	// Cluster Role Binding created
	ClusterRoleBindingName string
	// LoadBalancer Service created
	LoadBalanceServiceName string
}

type instanceResourcesResult struct {
	// Custom Resource Name of the instance
	InstanceName string
}

// SiteWhere Docker Image default tag name
const dockerImageDefaultTag = "3.0.0.beta3"

// Default configuration Template
const defaultConfigurationTemplate = "default"

// Default Dataset template
const defaultDatasetTemplate = "default"

// NewCreateInstance constructs a new *Install
func NewCreateInstance(cfg *Configuration) *CreateInstance {
	return &CreateInstance{
		cfg:                   cfg,
		InstanceName:          "",
		TenantName:            "default",
		Namespace:             "",
		Minimal:               false,
		Replicas:              1,
		Tag:                   dockerImageDefaultTag,
		Debug:                 false,
		ConfigurationTemplate: defaultConfigurationTemplate,
		DatasetTemplate:       defaultDatasetTemplate,
	}
}

// Run executes the list command, returning a set of matches.
func (i *CreateInstance) Run() (*instance.CreateSiteWhereInstance, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	var profile alpha3.SiteWhereProfile = alpha3.Default
	if i.Namespace == "" {
		i.Namespace = i.InstanceName
	}
	if i.Tag == "" {
		i.Tag = dockerImageDefaultTag
	}
	if i.ConfigurationTemplate == "" {
		i.ConfigurationTemplate = defaultConfigurationTemplate
	}
	if i.Minimal {
		profile = alpha3.Minimal
		i.ConfigurationTemplate = "minimal"
	}
	return i.createSiteWhereInstance(profile)
}

func (i *CreateInstance) createSiteWhereInstance(profile alpha3.SiteWhereProfile) (*instance.CreateSiteWhereInstance, error) {
	inr, err := i.createInstanceResources(profile)
	if err != nil {
		return nil, err
	}
	return &instance.CreateSiteWhereInstance{
		InstanceName:               i.InstanceName,
		Tag:                        i.Tag,
		Replicas:                   i.Replicas,
		Debug:                      i.Debug,
		ConfigurationTemplate:      i.ConfigurationTemplate,
		DatasetTemplate:            i.DatasetTemplate,
		InstanceCustomResourceName: inr.InstanceName,
	}, nil
}

// ExtractInstanceName returns the name of the instance that should be used.
func (i *CreateInstance) ExtractInstanceName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}

func (i *CreateInstance) createInstanceResources(profile alpha3.SiteWhereProfile) (*instanceResourcesResult, error) {
	var err error

	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	_, err = resources.CreateNamespaceIfNotExists(i.Namespace, clientset)
	if err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	client, err := i.cfg.ControllerClient()
	if err != nil {
		return nil, err
	}

	swInstanceCR := i.buildCRSiteWhereInstace()
	ctx := context.TODO()

	if err := client.Create(ctx, swInstanceCR); err != nil {
		if apierrors.IsAlreadyExists(err) {
			i.cfg.Log(fmt.Sprintf("Instance %s is already present. Skipping.", swInstanceCR.GetName()))
		} else {
			return nil, err
		}
	}

	return &instanceResourcesResult{
		InstanceName: i.InstanceName,
	}, nil
}

func (i *CreateInstance) buildInstanceServiceAccount() *v1.ServiceAccount {
	saName := fmt.Sprintf("sitewhere-instance-service-account-%s", i.Namespace)
	return &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: saName,
			Labels: map[string]string{
				"app": i.InstanceName,
			},
		},
	}
}

func (i *CreateInstance) buildInstanceClusterRole() *rbacV1.ClusterRole {
	roleName := fmt.Sprintf("sitewhere-instance-clusterrole-%s", i.InstanceName)
	return &rbacV1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleName,
			Labels: map[string]string{
				"app": i.InstanceName,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{
					"sitewhere.io",
				},
				Resources: []string{
					"instances",
					"instances/status",
					"microservices",
					"tenants",
					"tenantengines",
					"tenantengines/status",
				},
				Verbs: []string{
					"*",
				},
			}, {
				APIGroups: []string{
					"templates.sitewhere.io",
				},
				Resources: []string{
					"instanceconfigurations",
					"instancedatasets",
					"tenantconfigurations",
					"tenantengineconfigurations",
					"tenantdatasets",
					"tenantenginedatasets",
				},
				Verbs: []string{
					"*",
				},
			}, {
				APIGroups: []string{
					"scripting.sitewhere.io",
				},
				Resources: []string{
					"scriptcategories",
					"scripttemplates",
					"scripts",
					"scriptversions",
				},
				Verbs: []string{
					"*",
				},
			}, {
				APIGroups: []string{
					"apiextensions.k8s.io",
				},
				Resources: []string{
					"customresourcedefinitions",
				},
				Verbs: []string{
					"*",
				},
			},
		},
	}
}

func (i *CreateInstance) buildInstanceClusterRoleBinding(serviceAccount *v1.ServiceAccount,
	clusterRole *rbacV1.ClusterRole) *rbacV1.ClusterRoleBinding {
	roleBindingName := fmt.Sprintf("sitewhere-instance-clusterrole-binding-%s", i.InstanceName)
	return &rbacV1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleBindingName,
			Labels: map[string]string{
				"app": i.InstanceName,
			},
		},
		Subjects: []rbacV1.Subject{
			{
				Kind:      "ServiceAccount",
				Namespace: i.Namespace,
				Name:      serviceAccount.ObjectMeta.Name,
			},
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRole.ObjectMeta.Name,
		},
	}
}

func (i *CreateInstance) buildCRSiteWhereInstace() *sitewhereiov1alpha4.SiteWhereInstance {
	return &sitewhereiov1alpha4.SiteWhereInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       sitewhereiov1alpha4.SiteWhereInstanceKind,
			APIVersion: sitewhereiov1alpha4.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: i.InstanceName,
		},
		Spec: sitewhereiov1alpha4.SiteWhereInstanceSpec{
			ConfigurationTemplate: i.ConfigurationTemplate,
			DatasetTemplate:       i.DatasetTemplate,
			DockerSpec: &sitewhereiov1alpha4.DockerSpec{
				Registry:   sitewhereiov1alpha4.DefaultDockerSpec.Registry,
				Repository: sitewhereiov1alpha4.DefaultDockerSpec.Repository,
				Tag:        i.Tag,
			},
		},
	}
}
