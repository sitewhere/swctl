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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
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
	// IstioInject if true, label namespace for instio inject
	IstioInject bool
	// Minimal use minimal profile. Initialize only essential microservices.
	Minimal bool
	// Number of replicas
	Replicas int32
	// Registry is the docker registry of the microservices images
	Registry string
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

const (
	// Client Secret key
	clientSecretKey = "client-secret"

	// sitewhereGatewayName is the FQDN of sitewhere gateway
	//	sitewhereGatewayName = "sitewhere-gateway.sitewhere-system.svc.cluster.local"
	sitewhereGatewayName = "sitewhere-system/sitewhere-gateway"
)

// NewCreateInstance constructs a new *Install
func NewCreateInstance(cfg *Configuration) *CreateInstance {
	return &CreateInstance{
		cfg:                   cfg,
		InstanceName:          "",
		TenantName:            "default",
		Namespace:             "",
		IstioInject:           false,
		Minimal:               false,
		Replicas:              1,
		Tag:                   dockerImageDefaultTag,
		Registry:              sitewhereiov1alpha4.DefaultDockerSpec.Registry,
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
	_, err = resources.CreateNamespaceIfNotExists(i.Namespace, i.IstioInject, clientset)
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

	err = i.AddIstioVirtualService()
	if err != nil {
		return nil, err
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

	var defaultMicroservices = i.renderDefaultMicroservices()

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
				Registry:   i.Registry,
				Repository: sitewhereiov1alpha4.DefaultDockerSpec.Repository,
				Tag:        i.Tag,
			},
			Microservices: defaultMicroservices,
		},
	}
}

func (i *CreateInstance) renderDefaultMicroservices() []sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result []sitewhereiov1alpha4.SiteWhereMicroserviceSpec

	if i.Minimal {
		result = []sitewhereiov1alpha4.SiteWhereMicroserviceSpec{
			renderAssetManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderCommandDeliveryMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderDeviceManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderEventManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderEventSourcesMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderInboundProcessingMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderInstanceManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderOutboundConnectorsMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
		}
	} else {
		result = []sitewhereiov1alpha4.SiteWhereMicroserviceSpec{
			renderAssetManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderBatchOperationsMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderCommandDeliveryMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderDeviceManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderDeviceRegistrationMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderDeviceStateMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderEventManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderEventSourcesMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderInboundProcessingMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderInstanceManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderLabelGenerationMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderOutboundConnectorsMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
			renderScheduleManagementMicroservice(i.Replicas, i.Tag, i.Registry, i.InstanceName),
		}
	}

	return result
}

func renderDefaultMicroserviceTemplate(replicas int32,
	tag string,
	registry string,
	instanceName string,
	functionalArea string,
	multitenant bool,
	name string,
	nescription string,
	icon string,
	debugOffset int) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var clusterIPType = corev1.ServiceTypeClusterIP
	return sitewhereiov1alpha4.SiteWhereMicroserviceSpec{
		FunctionalArea: functionalArea,
		Replicas:       replicas,
		Multitenant:    multitenant,
		Name:           name,
		Description:    nescription,
		Icon:           icon,
		PodSpec: &sitewhereiov1alpha4.MicroservicePodSpecification{
			DockerSpec: &sitewhereiov1alpha4.DockerSpec{
				Registry:   registry,
				Repository: sitewhereiov1alpha4.DefaultDockerSpec.Repository,
				Tag:        tag,
			},
			ImagePullPolicy: corev1.PullIfNotPresent,
			Ports: []corev1.ContainerPort{
				corev1.ContainerPort{
					ContainerPort: 9000,
					Protocol:      corev1.ProtocolTCP,
				},
				corev1.ContainerPort{
					ContainerPort: 9090,
					Protocol:      corev1.ProtocolTCP,
				},
			},
			Env: []corev1.EnvVar{
				corev1.EnvVar{
					Name: "sitewhere.config.k8s.name",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.name",
						},
					},
				},
				corev1.EnvVar{
					Name: "sitewhere.config.k8s.namespace",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.namespace",
						},
					},
				},
				corev1.EnvVar{
					Name: "sitewhere.config.k8s.pod.ip",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "status.podIP",
						},
					},
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.product.id",
					Value: instanceName,
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.keycloak.service.name",
					Value: "sitewhere-keycloak-http",
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.keycloak.api.port",
					Value: "80",
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.keycloak.realm",
					Value: "sitewhere",
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.keycloak.master.realm",
					Value: "master",
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.keycloak.master.username",
					Value: "sitewhere",
				},
				corev1.EnvVar{
					Name:  "sitewhere.config.keycloak.master.password",
					Value: "sitewhere",
				},
				corev1.EnvVar{
					Name: "sitewhere.config.keycloak.oidc.secret",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: instanceName,
							},
							Key: clientSecretKey,
						},
					},
				},
			},
		},
		SerivceSpec: &sitewhereiov1alpha4.MicroserviceServiceSpecification{
			Type: &clusterIPType,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "grpc-api",
					Port:       9000,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{IntVal: 9000},
				},
				corev1.ServicePort{
					Name:       "http-metrics",
					Port:       9090,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{IntVal: 9090},
				},
			},
		},
		Debug: &sitewhereiov1alpha4.MicroserviceDebugSpecification{
			Enabled:  false,
			JDWPPort: 8000 + debugOffset,
			JMXPort:  1100 + debugOffset,
		},
		Logging: &sitewhereiov1alpha4.MicroserviceLoggingSpecification{
			Overrides: []sitewhereiov1alpha4.MicroserviceLoggingEntry{
				sitewhereiov1alpha4.MicroserviceLoggingEntry{
					Logger: "com.sitewhere",
					Level:  "info",
				},
				sitewhereiov1alpha4.MicroserviceLoggingEntry{
					Logger: "com.sitewhere.grpc.client",
					Level:  "info",
				},
				sitewhereiov1alpha4.MicroserviceLoggingEntry{
					Logger: "com.sitewhere.microservice.grpc",
					Level:  "info",
				},
				sitewhereiov1alpha4.MicroserviceLoggingEntry{
					Logger: "com.sitewhere.microservice.kafka",
					Level:  "info",
				},
				sitewhereiov1alpha4.MicroserviceLoggingEntry{
					Logger: "org.redisson",
					Level:  "info",
				},
				sitewhereiov1alpha4.MicroserviceLoggingEntry{
					Logger: "com.sitewhere.asset",
					Level:  "info",
				},
			},
		},
	}
}

func renderAssetManagementMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"asset-management",
		true,
		"Asset Management",
		"Provides APIs for managing assets associated with device assignments",
		"devices_other",
		6)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.asset",
		Level:  "info",
	})
	return result
}

func renderBatchOperationsMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"batch-operations",
		true,
		"Batch Operations",
		"Handles processing of operations which affect a large number of devices",
		"view_module",
		11)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.batch",
		Level:  "info",
	})
	return result
}

func renderCommandDeliveryMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"command-delivery",
		true,
		"Command Delivery",
		"Manages delivery of commands in various formats based on invocation events",
		"call_made",
		12)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.commands",
		Level:  "info",
	})
	return result
}

func renderDeviceManagementMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"device-management",
		true,
		"Device Management",
		"Provides APIs for managing the device object model",
		"developer_board",
		4)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.device",
		Level:  "info",
	})
	return result
}

func renderDeviceRegistrationMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"device-registration",
		true,
		"Device Registration",
		"Handles registration of new devices with the system",
		"add_box",
		13)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.registration",
		Level:  "info",
	})
	return result
}

func renderDeviceStateMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"device-state",
		true,
		"Device State",
		"Provides device state management features such as device shadows",
		"warning",
		14)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.devicestate",
		Level:  "info",
	})
	return result
}

func renderEventManagementMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"event-management",
		true,
		"Event Management",
		"Provides APIs for persisting and accessing events generated by devices",
		"dynamic_feed",
		5)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.event",
		Level:  "info",
	})
	return result
}

func renderEventSourcesMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"event-sources",
		true,
		"Event Sources",
		"Handles inbound device data from various sources, protocols, and formats",
		"forward",
		8)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.sources",
		Level:  "info",
	})
	return result
}

func renderInboundProcessingMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"inbound-processing",
		true,
		"Inbound Processing",
		"Common processing logic applied to enrich and direct inbound events",
		"input",
		7)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.inbound",
		Level:  "info",
	})
	return result
}

func renderInstanceManagementMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var imConfiguration = &sitewhereiov1alpha4.InstanceMangementConfiguration{
		UserManagementConfiguration: &sitewhereiov1alpha4.UserManagementConfiguration{
			SyncopeHost:            "sitewhere-syncope.sitewhere-system.cluster.local",
			SyncopePort:            8080,
			JWTExpirationInMinutes: 60,
		},
	}
	marshalledBytes, err := json.Marshal(imConfiguration)
	if err != nil {
		return sitewhereiov1alpha4.SiteWhereMicroserviceSpec{}
	}
	var instanceManagementConfiguration = &runtime.RawExtension{}
	err = instanceManagementConfiguration.UnmarshalJSON(marshalledBytes)
	if err != nil {
		return sitewhereiov1alpha4.SiteWhereMicroserviceSpec{}
	}

	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"instance-management",
		false,
		"Instance Management",
		"Handles APIs for managing global aspects of an instance",
		"language",
		1)

	result.Configuration = instanceManagementConfiguration
	result.PodSpec.Ports = append(result.PodSpec.Ports, corev1.ContainerPort{
		ContainerPort: 8080,
		Protocol:      corev1.ProtocolTCP,
	})
	result.SerivceSpec.Ports = append(result.SerivceSpec.Ports, corev1.ServicePort{
		Name:       "http-rest",
		Port:       8080,
		Protocol:   corev1.ProtocolTCP,
		TargetPort: intstr.IntOrString{IntVal: 8080},
	})
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.instance",
		Level:  "info",
	})
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.web",
		Level:  "info",
	})
	return result
}

func renderLabelGenerationMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"label-generation",
		true,
		"Label Generation",
		"Supports generating labels such as bar codes and QR codes for devices",
		"label",
		9)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.labels",
		Level:  "info",
	})
	return result
}

func renderOutboundConnectorsMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"outbound-connectors",
		true,
		"Outbound Connectors",
		"Allows event streams to be delivered to external systems for additional processing",
		"label",
		16)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.connectors",
		Level:  "info",
	})
	return result
}

func renderScheduleManagementMicroservice(replicas int32, tag string, registry string, instanceName string) sitewhereiov1alpha4.SiteWhereMicroserviceSpec {
	var result = renderDefaultMicroserviceTemplate(replicas,
		tag,
		registry,
		instanceName,
		"schedule-management",
		true,
		"Schedule Management",
		"Supports scheduling of various system operations",
		"label",
		18)
	result.Logging.Overrides = append(result.Logging.Overrides, sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: "com.sitewhere.schedule",
		Level:  "info",
	})
	return result
}

// AddIstioVirtualService install Istio Virtual Service
func (i *CreateInstance) AddIstioVirtualService() error {
	restconfig, err := i.cfg.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return err
	}
	ic, err := versionedclient.NewForConfig(restconfig)
	if err != nil {
		return err
	}

	var vsName = fmt.Sprintf("%s-vs", i.InstanceName)
	var vsRouteHost = fmt.Sprintf("instance-management.%s.svc.cluster.local", i.Namespace)
	var prefixValue = fmt.Sprintf("/%s", i.InstanceName)

	var vs *v1alpha3.VirtualService = &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: i.Namespace,
			Name:      vsName,
		},
		Spec: networkingv1alpha3.VirtualService{
			Gateways: []string{
				sitewhereGatewayName,
			},
			Hosts: []string{
				"*",
			},
			Http: []*networkingv1alpha3.HTTPRoute{
				&networkingv1alpha3.HTTPRoute{
					Name: "instance-rest",
					Match: []*networkingv1alpha3.HTTPMatchRequest{
						&networkingv1alpha3.HTTPMatchRequest{
							Uri: &networkingv1alpha3.StringMatch{
								MatchType: &networkingv1alpha3.StringMatch_Prefix{
									Prefix: prefixValue,
								},
							},
						},
					},
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						&networkingv1alpha3.HTTPRouteDestination{
							Destination: &networkingv1alpha3.Destination{
								Host: vsRouteHost,
								Port: &networkingv1alpha3.PortSelector{
									Number: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = ic.NetworkingV1alpha3().VirtualServices(i.Namespace).Create(context.TODO(), vs, metav1.CreateOptions{})

	return nil
}
