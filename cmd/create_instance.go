/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/sitewhere/swctl/internal"
	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// SiteWhere Docker Image default tag name
const dockerImageDefaultTag = "3.0.0.beta1"

// createInstanceCmd represents the instance command
var (
	namespace         = ""    // Namespace to use
	minimal           = false // Use minimal profile. Initialize only essential microservices.
	tag               = ""    // Docker image tag
	debug             = false // Use debug mode
	createInstanceCmd = &cobra.Command{
		Use:   "instance <name>",
		Short: "Create SiteWhere Instance",
		Long: `Use this command to create an Instance of SiteWhere.
For example, to create an instance with name "sitewhere" use:

  swctl create instance sitewhere

To create an instance with the minimal profile use:

  swctl create instance sitewhere -m
`,
		Args: commandCreateInstanceArgs,
		Run:  commandCreateInstanceRun,
	}
)

func init() {
	createInstanceCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the instance.")
	createInstanceCmd.Flags().BoolVarP(&minimal, "minimal", "m", false, "Minimal installation.")
	createInstanceCmd.Flags().StringVarP(&tag, "tag", "t", "", "Docker image tag.")
	createInstanceCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode.")
	createCmd.AddCommand(createInstanceCmd)
}

func commandCreateInstanceArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires one argument")
	}
	return nil
}

func commandCreateInstanceRun(cmd *cobra.Command, args []string) {
	name := args[0]
	var profile alpha3.SiteWhereProfile = alpha3.Default

	if namespace == "" {
		namespace = name
	}
	if tag == "" {
		tag = dockerImageDefaultTag
	}

	var configurationTemplate = "default"

	if minimal {
		profile = alpha3.Minimal
		configurationTemplate = "minimal"
	}

	instance := alpha3.SiteWhereInstance{
		Name:                  name,
		Namespace:             namespace,
		Tag:                   tag,
		Debug:                 debug,
		ConfigurationTemplate: configurationTemplate,
		DatasetTemplate:       "default",
		Profile:               profile}

	createSiteWhereInstance(&instance)
}

func createSiteWhereInstance(instance *alpha3.SiteWhereInstance) {

	config, err := createNamespaceAndResources(instance)

	if err != nil {
		fmt.Printf("Error Setting Namespace and Resources: %v\n", err)
		return
	}

	createSiteWhereResources(instance, instance.Namespace, config)

	fmt.Printf("SiteWhere Instance '%s' created\n", instance.Name)
}

func createNamespaceAndResources(instance *alpha3.SiteWhereInstance) (*rest.Config, error) {
	var err error

	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes Config: %v\n", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return nil, err
	}

	var ns *v1.Namespace
	ns, err = internal.CreateNamespaceIfNotExists(instance.Namespace, clientset)
	if err != nil {
		fmt.Printf("Error Creating Namespace: %s, %v", instance.Namespace, err)
		return nil, err
	}

	var namespace = ns.ObjectMeta.Name

	var sa *v1.ServiceAccount
	sa, err = createServiceAccountIfNotExists(instance, namespace, clientset)
	if err != nil {
		fmt.Printf("Error Creating Service Account: %s, %v", instance.Namespace, err)
		return nil, err
	}

	var role *rbacV1.Role
	role, err = createRoleIfNotExists(instance, namespace, clientset)
	if err != nil {
		fmt.Printf("Error Creating Role: %s, %v", instance.Namespace, err)
		return nil, err
	}

	_, err = createRoleBindingIfNotExists(instance, namespace, sa, role, clientset)
	if err != nil {
		fmt.Printf("Error Creating Role Binding: %s, %v", instance.Namespace, err)
		return nil, err
	}

	_, err = createLoadBalancerServiceIfNotExists(instance, namespace, clientset)
	if err != nil {
		fmt.Printf("Error Creating Loadbalancer Service: %s, %v", instance.Namespace, err)
		return nil, err
	}

	return config, nil
}

func createSiteWhereResources(instance *alpha3.SiteWhereInstance, namespace string, config *rest.Config) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return
	}

	_, err = createCRSiteWhereInstaceIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating CR SiteWhere Instace: %v\n", err)
		return
	}

	var microservices = alpha3.GetSiteWhereMicroservicesList()
	for _, micrservice := range microservices {

		if micrservice.ID == "instance-management" {
			_, err = createCRSiteWhereInstanceManagementIfNotExists(instance, namespace, client)
			if err != nil {
				fmt.Printf("Error Creating SiteWhere Instance Management Microservice: %v\n", err)
				return
			}
		} else if instance.Profile == alpha3.Default || instance.Profile != micrservice.Profile {
			_, err = createCRSiteWhereMicroserviceIfNotExists(instance, namespace, &micrservice, client)
			if err != nil {
				fmt.Printf("Error Creating SiteWhere %s Microservice: %v\n", micrservice.Name, err)
			}
		}
	}
}

func createServiceAccountIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, clientset *kubernetes.Clientset) (*v1.ServiceAccount, error) {
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

func createRoleIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, clientset *kubernetes.Clientset) (*rbacV1.Role, error) {
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

func createRoleBindingIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, serviceAccount *v1.ServiceAccount,
	role *rbacV1.Role, clientset *kubernetes.Clientset) (*rbacV1.RoleBinding, error) {
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

func createLoadBalancerServiceIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, clientset *kubernetes.Clientset) (*v1.Service, error) {
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
					"configuration": map[string]interface{}{
						"test": 123,
					},
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        microservice.Name,
					"description": microservice.Description,
					"icon":        microservice.Icon,
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
							}
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
