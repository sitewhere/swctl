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
	"fmt"
	"strings"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"
	"github.com/sitewhere/swctl/pkg/resources/grv"
)

// CreateInstance is the action for creating a SiteWhere instance
type CreateInstance struct {
	cfg *Configuration
	// Name of the instance
	InstanceName string
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
	CRName string
	// Microservices created
	Microservices []string
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
	// nsr, err := i.createNamespaceAndResources()
	// if err != nil {
	// 	return nil, err
	// }

	inr, err := i.createInstanceResources(profile)
	if err != nil {
		return nil, err
	}
	return &instance.CreateSiteWhereInstance{
		InstanceName: i.InstanceName,
		//		Namespace:                  nsr.Namespace,
		Tag:                   i.Tag,
		Replicas:              i.Replicas,
		Debug:                 i.Debug,
		ConfigurationTemplate: i.ConfigurationTemplate,
		DatasetTemplate:       i.DatasetTemplate,
		// ServiceAccountName:         nsr.ServiceAccountName,
		// ClusterRoleName:            nsr.ClusterRoleName,
		// ClusterRoleBindingName:     nsr.ClusterRoleBindingName,
		// LoadBalanceServiceName:     nsr.LoadBalanceServiceName,
		InstanceCustomResourceName: inr.CRName,
	}, nil
}

// ExtractInstanceName returns the name of the instance that should be used.
func (i *CreateInstance) ExtractInstanceName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}

func (i *CreateInstance) createNamespaceAndResources() (*namespaceAndResourcesResult, error) {
	//	var err error
	// clientset, err := i.cfg.KubernetesClientSet()
	// if err != nil {
	// 	return nil, err
	// }
	// ns, err := resources.CreateNamespaceIfNotExists(i.Namespace, clientset)
	// if err != nil {
	// 	return nil, err
	// }
	// sa, err := resources.CreateServiceAccountIfNotExists(
	// 	i.buildInstanceServiceAccount(), clientset, i.Namespace)
	// if err != nil {
	// 	return nil, err
	// }
	// clusterRole, err := resources.CreateClusterRoleIfNotExists(
	// 	i.buildInstanceClusterRole(), clientset)
	// if err != nil {
	// 	return nil, err
	// }
	// clusterRoleBinding, err := resources.CreateClusterRoleBindingIfNotExists(
	// 	i.buildInstanceClusterRoleBinding(sa, clusterRole), clientset)
	// if err != nil {
	// 	return nil, err
	// }
	// loadBalanceService, err := resources.CreateServiceIfNotExists(
	// 	i.buildLoadBalancerService(), clientset, i.Namespace)
	// if err != nil {
	// 	return nil, err
	// }
	return &namespaceAndResourcesResult{
		// Namespace:              ns.ObjectMeta.Name,
		// ServiceAccountName:     sa.ObjectMeta.Name,
		// ClusterRoleName:        clusterRole.ObjectMeta.Name,
		// ClusterRoleBindingName: clusterRoleBinding.ObjectMeta.Name,
		// LoadBalanceServiceName: loadBalanceService.ObjectMeta.Name,
	}, nil
}

func (i *CreateInstance) createInstanceResources(profile alpha3.SiteWhereProfile) (*instanceResourcesResult, error) {
	var err error
	dynamicClientset, err := i.cfg.KubernetesDynamicClientSet()
	if err != nil {
		return nil, err
	}
	icr, err := resources.CreateSiteWhereInstanceCR(i.buildCRSiteWhereInstace(), dynamicClientset)
	if err != nil {
		return nil, err
	}
	//	var microservices = alpha3.GetSiteWhereMicroservicesList()
	// var installedMicroservices []string

	// for _, micrservice := range microservices {
	// 	var msCR *unstructured.Unstructured
	// 	if micrservice.ID == "instance-management" {
	// 		msCR = i.buildCRSiteWhereMicroserviceInstanceManagement()
	// 	} else if profile == alpha3.Default || profile != micrservice.Profile {
	// 		msCR = i.buildCRSiteWhereMicroservice(&micrservice)
	// 	}
	// 	mrc, err := resources.CreateSiteWhereMicroserviceCR(msCR, i.Namespace, dynamicClientset)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	installedMicroservices = append(installedMicroservices, mrc.GetName())
	// }
	return &instanceResourcesResult{
		CRName: icr.GetName(),
		// Microservices: installedMicroservices,
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

func (i *CreateInstance) buildLoadBalancerService() *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sitewhere-rest-http",
			Labels: map[string]string{
				"app": i.InstanceName,
			},
		},
		Spec: v1.ServiceSpec{
			Type: "LoadBalancer",
			Ports: []v1.ServicePort{
				{
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
					Protocol:   v1.ProtocolTCP,
					Name:       "http-rest",
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/instance": i.InstanceName,
				"sitewhere.io/name":          "instance-management",
			},
		},
	}
}

func (i *CreateInstance) buildCRSiteWhereInstace() *unstructured.Unstructured {
	sitewhereInstanceGVR := grv.SiteWhereInstanceGRV()
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "SiteWhereInstance",
			"apiVersion": sitewhereInstanceGVR.Group + "/" + sitewhereInstanceGVR.Version,
			"metadata": map[string]interface{}{
				"name": i.InstanceName,
			},
			"spec": map[string]interface{}{
				"configurationTemplate": i.ConfigurationTemplate,
				"datasetTemplate":       i.DatasetTemplate,
				"dockerSpec": map[string]interface{}{
					"registry":   "docker.io",
					"repository": "sitewhere",
					"tag":        i.Tag,
				},
			},
		},
	}
}

func (i *CreateInstance) buildCRSiteWhereMicroserviceInstanceManagement() *unstructured.Unstructured {
	sitewhereMicroserviceGVR := grv.SiteWhereMicroserviceGRV()
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "SiteWhereMicroservice",
			"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
			"metadata": map[string]interface{}{
				"name":      "instance-management-microservice",
				"namespace": i.Namespace,
				"labels": map[string]interface{}{
					"sitewhere.io/instance":        i.InstanceName,
					"sitewhere.io/functional-area": "instance-management",
				},
			},
			"spec": map[string]interface{}{
				"replicas":    i.Replicas,
				"name":        "Instance Management",
				"description": "Handles APIs for managing global aspects of an instance",
				"icon":        "language",
				"logging": map[string]interface{}{
					"overrides": []map[string]interface{}{
						{
							"logger": "com.sitewhere",
							"level":  "info",
						},
						{
							"logger": "com.sitewhere.grpc.client",
							"level":  "info",
						},
						{
							"logger": "com.sitewhere.microservice.grpc",
							"level":  "info",
						},
						{
							"logger": "com.sitewhere.microservice.kafka",
							"level":  "info",
						},
						{
							"logger": "org.redisson",
							"level":  "info",
						},
						{
							"level":  "info",
							"logger": "com.sitewhere.instance",
						},
						{
							"level":  "info",
							"logger": "com.sitewhere.web",
						},
					},
				},
				"configuration": map[string]interface{}{
					"userManagement": map[string]interface{}{
						"syncopeHost":            "sitewhere-syncope.sitewhere-system",
						"syncopePort":            8080,
						"jwtExpirationInMinutes": 60,
					},
				},
				"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
					"chartName":      "sitewhere-0.3.0",
					"releaseName":    i.InstanceName,
					"releaseService": "Tiller",
				},
				"podSpec": map[string]interface{}{
					"imageRegistry":   "docker.io",
					"imageRepository": "sitewhere",
					"imageTag":        i.Tag,
					"imagePullPolicy": "IfNotPresent",
					"ports": []map[string]interface{}{
						{
							"containerPort": 8080,
						},
						{
							"containerPort": 9000,
						},
						{
							"containerPort": 9090,
						},
					},
					"env": []map[string]interface{}{
						{
							"name": "sitewhere.config.k8s.name",
							"valueFrom": map[string]interface{}{
								"fieldRef": map[string]interface{}{
									"fieldPath": "metadata.name",
								},
							},
						},
						{
							"name": "sitewhere.config.k8s.namespace",
							"valueFrom": map[string]interface{}{
								"fieldRef": map[string]interface{}{
									"fieldPath": "metadata.namespace",
								},
							},
						},
						{
							"name": "sitewhere.config.k8s.pod.ip",
							"valueFrom": map[string]interface{}{
								"fieldRef": map[string]interface{}{
									"fieldPath": "status.podIP",
								},
							},
						},
					},
				},
				"serviceSpec": map[string]interface{}{
					"type": "ClusterIP",
					"ports": []map[string]interface{}{
						{
							"port":       8080,
							"targetPort": 8080,
							"protocol":   "TCP",
							"name":       "http-rest",
						},
						{
							"port":       9000,
							"targetPort": 9000,
							"protocol":   "TCP",
							"name":       "grpc-api",
						},
						{
							"port":       9090,
							"targetPort": 9090,
							"protocol":   "TCP",
							"name":       "http-metrics",
						},
					},
				},
				"debug": map[string]interface{}{
					"enabled":  i.Debug,
					"jdwpPort": 8001,
					"jmxPort":  1101,
				},
			},
		},
	}
}

func (i *CreateInstance) buildCRSiteWhereMicroservice(microservice *alpha3.SiteWhereMicroservice) *unstructured.Unstructured {
	msName := fmt.Sprintf("%s-microservice", microservice.ID)
	sitewhereMicroserviceGVR := grv.SiteWhereMicroserviceGRV()
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "SiteWhereMicroservice",
			"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
			"metadata": map[string]interface{}{
				"name":      msName,
				"namespace": i.Namespace,
				"labels": map[string]interface{}{
					"sitewhere.io/instance":        i.InstanceName,
					"sitewhere.io/functional-area": microservice.ID,
				},
			},
			"spec": map[string]interface{}{
				"configuration": map[string]interface{}{},
				"replicas":      i.Replicas,
				"multitenant":   true,
				"name":          microservice.Name,
				"description":   microservice.Description,
				"icon":          microservice.Icon,
				"logging": map[string]interface{}{
					"overrides": []map[string]interface{}{
						{
							"logger": "com.sitewhere",
							"level":  "info",
						},
						{
							"logger": "com.sitewhere.grpc.client",
							"level":  "info",
						},
						{
							"logger": "com.sitewhere.microservice.grpc",
							"level":  "info",
						},
						{
							"logger": "com.sitewhere.microservice.kafka",
							"level":  "info",
						},
						{
							"logger": "org.redisson",
							"level":  "info",
						},
						{
							"level":  "info",
							"logger": microservice.Logger,
						},
					},
				},
				"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
					"chartName":      "sitewhere-0.3.0",
					"releaseName":    i.InstanceName,
					"releaseService": "Tiller",
				},
				"podSpec": map[string]interface{}{
					"imageRegistry":   "docker.io",
					"imageRepository": "sitewhere",
					"imageTag":        i.Tag,
					"imagePullPolicy": "IfNotPresent",
					"ports": []map[string]interface{}{
						{
							"containerPort": 9000,
						},
						{
							"containerPort": 9090,
						},
					},
					"env": []map[string]interface{}{
						{
							"name": "sitewhere.config.k8s.name",
							"valueFrom": map[string]interface{}{
								"fieldRef": map[string]interface{}{
									"fieldPath": "metadata.name",
								},
							},
						},
						{
							"name": "sitewhere.config.k8s.namespace",
							"valueFrom": map[string]interface{}{
								"fieldRef": map[string]interface{}{
									"fieldPath": "metadata.namespace",
								},
							},
						},
						{
							"name": "sitewhere.config.k8s.pod.ip",
							"valueFrom": map[string]interface{}{
								"fieldRef": map[string]interface{}{
									"fieldPath": "status.podIP",
								},
							},
						},
					},
				},
				"serviceSpec": map[string]interface{}{
					"type": "ClusterIP",
					"ports": []map[string]interface{}{
						{
							"port":       9000,
							"targetPort": 9000,
							"protocol":   "TCP",
							"name":       "grpc-api",
						},
						{
							"port":       9090,
							"targetPort": 9090,
							"protocol":   "TCP",
							"name":       "http-metrics",
						},
					},
				},
				"debug": map[string]interface{}{
					"enabled":  i.Debug,
					"jdwpPort": 8000 + microservice.PortOffset,
					"jmxPort":  1100 + microservice.PortOffset,
				},
			},
		},
	}
}
