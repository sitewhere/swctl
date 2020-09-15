/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"
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
	// Docker image tag
	Tag string
	// Use debug mode
	Debug bool
	// Configuration Template
	ConfigurationTemplate string
	// Dataset template
	DatasetTemplate string
}

// SiteWhere Docker Image default tag name
const dockerImageDefaultTag = "3.0.0.beta1"

// Default configuration Template
const defaultConfigurationTemplate = "default"

// Detaul Dataset template
const defaultDatasetTemplate = "default"

// NewCreateInstance constructs a new *Install
func NewCreateInstance(cfg *Configuration) *CreateInstance {
	return &CreateInstance{
		cfg:                   cfg,
		InstanceName:          "",
		Namespace:             "",
		Minimal:               false,
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
	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	dynamicClientset, err := i.cfg.KubernetesDynamicClientSet()
	if err != nil {
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

	instanceToCreate := alpha3.SiteWhereInstance{
		Name:                  i.InstanceName,
		Namespace:             i.Namespace,
		Tag:                   i.Tag,
		Debug:                 i.Debug,
		ConfigurationTemplate: i.ConfigurationTemplate,
		DatasetTemplate:       i.DatasetTemplate,
		Profile:               profile}

	err = createNamespaceAndResources(clientset, &instanceToCreate)

	if err != nil {
		return nil, err
	}

	err = createSiteWhereResources(dynamicClientset, &instanceToCreate, instanceToCreate.Namespace)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("SiteWhere Instance '%s' created\n", instance.Name)
	return &instance.CreateSiteWhereInstance{}, nil
}

// ExtractInstanceName returns the name of the instance that should be used.
func (i *CreateInstance) ExtractInstanceName(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}
	return args[0], nil
}

func createNamespaceAndResources(clientset kubernetes.Interface, instance *alpha3.SiteWhereInstance) error {
	var err error

	var ns *v1.Namespace
	ns, err = resources.CreateNamespaceIfNotExists(instance.Namespace, clientset)
	if err != nil {
		return err
	}
	var namespace = ns.ObjectMeta.Name
	var sa *v1.ServiceAccount
	sa, err = createServiceAccountIfNotExists(clientset, instance, namespace)
	if err != nil {
		return err
	}
	var role *rbacV1.ClusterRole
	role, err = createClusterRoleIfNotExists(clientset, instance)
	if err != nil {
		return err
	}
	_, err = createClusterRoleBindingIfNotExists(clientset, instance, sa, role)
	if err != nil {
		return err
	}
	_, err = createLoadBalancerServiceIfNotExists(clientset, instance, namespace)
	if err != nil {
		return err
	}
	return nil
}

func createSiteWhereResources(client dynamic.Interface, instance *alpha3.SiteWhereInstance, namespace string) error {
	var err error

	_, err = createCRSiteWhereInstaceIfNotExists(instance, namespace, client)
	if err != nil {
		// fmt.Printf("Error Creating CR SiteWhere Instace: %v\n", err)
		return err
	}

	var microservices = alpha3.GetSiteWhereMicroservicesList()
	for _, micrservice := range microservices {

		if micrservice.ID == "instance-management" {
			_, err = createCRSiteWhereInstanceManagementIfNotExists(instance, namespace, client)
			if err != nil {
				// fmt.Printf("Error Creating SiteWhere Instance Management Microservice: %v\n", err)
				return err
			}
		} else if instance.Profile == alpha3.Default || instance.Profile != micrservice.Profile {
			_, err = createCRSiteWhereMicroserviceIfNotExists(instance, namespace, &micrservice, client)
			if err != nil {
				// fmt.Printf("Error Creating SiteWhere %s Microservice: %v\n", micrservice.Name, err)
				return err
			}
		}
	}
	return nil
}

func createServiceAccountIfNotExists(clientset kubernetes.Interface, instance *alpha3.SiteWhereInstance, namespace string) (*v1.ServiceAccount, error) {
	var err error
	var sa *v1.ServiceAccount

	saName := fmt.Sprintf("sitewhere-instance-service-account-%s", namespace)

	sa, err = clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), saName, metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		sa = &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: saName,
				Labels: map[string]string{
					"app": instance.Name,
				},
			},
		}

		result, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(),
			sa,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return sa, nil
}

func createRoleIfNotExists(clientset kubernetes.Interface, instance *alpha3.SiteWhereInstance, namespace string) (*rbacV1.Role, error) {
	var err error
	var role *rbacV1.Role

	roleName := fmt.Sprintf("sitewhere-instance-role-%s", namespace)

	role, err = clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil && k8serror.IsNotFound(err) {
		role = &rbacV1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name: roleName,
				Labels: map[string]string{
					"app": instance.Name,
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

		result, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(),
			role,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return role, nil
}

func createRoleBindingIfNotExists(clientset kubernetes.Interface,
	instance *alpha3.SiteWhereInstance,
	namespace string,
	serviceAccount *v1.ServiceAccount,
	role *rbacV1.Role) (*rbacV1.RoleBinding, error) {
	var err error
	var roleBinding *rbacV1.RoleBinding

	roleBindingName := fmt.Sprintf("sitewhere-instance-role-binding-%s", namespace)

	roleBinding, err = clientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
	if err != nil && k8serror.IsNotFound(err) {
		roleBinding = &rbacV1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: roleBindingName,
				Labels: map[string]string{
					"app": instance.Name,
				},
			},
			Subjects: []rbacV1.Subject{
				{
					Kind:      "ServiceAccount",
					Namespace: namespace,
					Name:      serviceAccount.ObjectMeta.Name,
				},
			},
			RoleRef: rbacV1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     role.ObjectMeta.Name,
			},
		}

		result, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(),
			roleBinding,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return roleBinding, nil
}

func createClusterRoleIfNotExists(
	clientset kubernetes.Interface,
	instance *alpha3.SiteWhereInstance) (*rbacV1.ClusterRole, error) {
	var err error
	var clusterRole *rbacV1.ClusterRole

	roleName := fmt.Sprintf("sitewhere-instance-clusterrole-%s", instance.Name)

	clusterRole, err = clientset.RbacV1().ClusterRoles().Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil && k8serror.IsNotFound(err) {
		clusterRole = &rbacV1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: roleName,
				Labels: map[string]string{
					"app": instance.Name,
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

		result, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(),
			clusterRole,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return clusterRole, nil
}

func createClusterRoleBindingIfNotExists(
	clientset kubernetes.Interface,
	instance *alpha3.SiteWhereInstance, serviceAccount *v1.ServiceAccount,
	clusterRole *rbacV1.ClusterRole) (*rbacV1.ClusterRoleBinding, error) {
	var err error
	var clusterRoleBinding *rbacV1.ClusterRoleBinding

	roleBindingName := fmt.Sprintf("sitewhere-instance-clusterrole-binding-%s", instance.Name)

	clusterRoleBinding, err = clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), roleBindingName, metav1.GetOptions{})
	if err != nil && k8serror.IsNotFound(err) {
		clusterRoleBinding = &rbacV1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: roleBindingName,
				Labels: map[string]string{
					"app": instance.Name,
				},
			},
			Subjects: []rbacV1.Subject{
				{
					Kind:      "ServiceAccount",
					Namespace: instance.Namespace,
					Name:      serviceAccount.ObjectMeta.Name,
				},
			},
			RoleRef: rbacV1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     clusterRole.ObjectMeta.Name,
			},
		}

		result, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(),
			clusterRoleBinding,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return clusterRoleBinding, nil
}

func createLoadBalancerServiceIfNotExists(clientset kubernetes.Interface,
	instance *alpha3.SiteWhereInstance,
	namespace string) (*v1.Service, error) {
	var err error
	var service *v1.Service

	service, err = clientset.CoreV1().Services(namespace).Get(context.TODO(), "sitewhere-rest-http", metav1.GetOptions{})
	if err != nil && k8serror.IsNotFound(err) {

		service = &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "sitewhere-rest-http",
				Labels: map[string]string{
					"app": instance.Name,
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
					"app.kubernetes.io/instance": instance.Name,
					"sitewhere.io/name":          "instance-management",
				},
			},
		}

		result, err := clientset.CoreV1().Services(namespace).Create(context.TODO(),
			service,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}
	return service, nil
}

func createCRSiteWhereInstaceIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereInstanceGVR)

	sitewhereInstaces, err := res.Get(context.TODO(), instance.Name, metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		sitewhereInstaces = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereInstance",
				"apiVersion": sitewhereInstanceGVR.Group + "/" + sitewhereInstanceGVR.Version,
				"metadata": map[string]interface{}{
					"name": instance.Name,
				},
				"spec": map[string]interface{}{
					"instanceNamespace":     instance.Namespace,
					"configurationTemplate": instance.ConfigurationTemplate,
					"datasetTemplate":       instance.DatasetTemplate,
				},
			},
		}

		result, err := res.Create(context.TODO(), sitewhereInstaces, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}

	return sitewhereInstaces, nil
}

func createCRSiteWhereInstanceManagementIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	instanceManagementMS, err := res.Get(context.TODO(), "instance-management-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		instanceManagementMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "instance-management-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "instance-management",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
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
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        instance.Tag,
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
						"enabled":  instance.Debug,
						"jdwpPort": 8001,
						"jmxPort":  1101,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), instanceManagementMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return instanceManagementMS, nil
}

func createCRSiteWhereMicroserviceIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, microservice *alpha3.SiteWhereMicroservice, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	msName := fmt.Sprintf("%s-microservice", microservice.ID)

	sitewhereMicroservice, err := res.Get(context.TODO(), msName, metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		sitewhereMicroservice = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      msName,
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": microservice.ID,
					},
				},
				"spec": map[string]interface{}{
					"configuration": map[string]interface{}{},
					"replicas":      1, // TODO from parameter
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
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        instance.Tag,
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
						"enabled":  instance.Debug,
						"jdwpPort": 8000 + microservice.PortOffset,
						"jmxPort":  1100 + microservice.PortOffset,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), sitewhereMicroservice, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return sitewhereMicroservice, nil
}
