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

// createInstanceCmd represents the instance command
var (
	namespace         = ""
	createInstanceCmd = &cobra.Command{
		Use:   "instance",
		Short: "Create SiteWhere Instance",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires one argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if namespace == "" {
				namespace = name
			}

			instance := alpha3.SiteWhereInstance{
				Name:                  name,
				Namespace:             namespace,
				ConfigurationTemplate: "default",
				DatasetTemplate:       "default"}

			createSiteWhereInstance(&instance)
		},
	}
)

func init() {
	createInstanceCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the instance.")
	createCmd.AddCommand(createInstanceCmd)
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
	ns, err = createNamespaceIfNotExists(instance.Namespace, clientset)
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

	_, err = createCRSiteWhereAssetManagementIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Asset Management Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereBatchOperationsIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Batch Operations Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereCommandDeliveryIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Command Delivery Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereDeviceManagementIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Device Management Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereDeviceRegistrationIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Device Registration Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereDeviceStateIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Device State Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereEventManagementIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Event Management Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereEventSourcesIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Event Sources Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereInboundProcessingIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Inbound Processing Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereInstanceManagementIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Instance Management Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereLabelGenerationIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Label Generation Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereOutboundConnectorsIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Outbound Connectors Microservice: %v\n", err)
		return
	}

	_, err = createCRSiteWhereScheduleManagementIfNotExists(instance, namespace, client)
	if err != nil {
		fmt.Printf("Error Creating SiteWhere Schedule Management Microservice: %v\n", err)
		return
	}
}

func createNamespaceIfNotExists(namespace string, clientset *kubernetes.Clientset) (*v1.Namespace, error) {
	var err error
	var ns *v1.Namespace

	ns, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		ns = &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
				Labels: map[string]string{
					"app": namespace,
				},
			},
		}

		result, err := clientset.CoreV1().Namespaces().Create(context.TODO(),
			ns,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return ns, nil
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

func createCRSiteWhereAssetManagementIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	assetManagementMS, err := res.Get(context.TODO(), "asset-management-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		assetManagementMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "asset-management-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "asset-management",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Asset Management",
					"description": "Provides APIs for managing assets associated with device assignments",
					"icon":        "devices_other",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8006,
						"jmxPort":  1106,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), assetManagementMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return assetManagementMS, nil
}

func createCRSiteWhereBatchOperationsIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	batchOperationsMS, err := res.Get(context.TODO(), "batch-operations-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		batchOperationsMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "batch-operations-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "batch-operations",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Batch Operations",
					"description": "Handles processing of operations which affect a large number of devices",
					"icon":        "view_module",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8011,
						"jmxPort":  1111,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), batchOperationsMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return batchOperationsMS, nil
}

func createCRSiteWhereCommandDeliveryIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	commandDeliveryMS, err := res.Get(context.TODO(), "command-delivery-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		commandDeliveryMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "command-delivery-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "command-delivery",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Command Delivery",
					"description": "Manages delivery of commands in various formats based on invocation events",
					"icon":        "call_made",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8012,
						"jmxPort":  1112,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), commandDeliveryMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return commandDeliveryMS, nil
}

func createCRSiteWhereDeviceManagementIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	deviceManagementMS, err := res.Get(context.TODO(), "device-management-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		deviceManagementMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "device-management-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "device-management",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Device Management",
					"description": "Provides APIs for managing the device object model",
					"icon":        "developer_board",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8004,
						"jmxPort":  1104,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), deviceManagementMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return deviceManagementMS, nil
}

func createCRSiteWhereDeviceRegistrationIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	deviceRegistrationMS, err := res.Get(context.TODO(), "device-registration-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		deviceRegistrationMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "device-registration-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "device-registration",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Device Registration",
					"description": "Handles registration of new devices with the system",
					"icon":        "add_box",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8013,
						"jmxPort":  1113,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), deviceRegistrationMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return deviceRegistrationMS, nil
}

func createCRSiteWhereDeviceStateIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	deviceStateMS, err := res.Get(context.TODO(), "device-state-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		deviceStateMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "device-state-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "device-state",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Device State",
					"description": "Provides device state management features such as device shadows",
					"icon":        "warning",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8014,
						"jmxPort":  1114,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), deviceStateMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return deviceStateMS, nil
}

func createCRSiteWhereEventManagementIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	eventManagementMS, err := res.Get(context.TODO(), "event-management-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		eventManagementMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "event-management-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "event-management",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Event Management",
					"description": "Provides APIs for persisting and accessing events generated by devices",
					"icon":        "dynamic_feed",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8005,
						"jmxPort":  1105,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), eventManagementMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return eventManagementMS, nil
}

func createCRSiteWhereEventSourcesIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	eventSourcesMS, err := res.Get(context.TODO(), "event-sources-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		eventSourcesMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "event-sources-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "event-sources",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Event Sources",
					"description": "Handles inbound device data from various sources, protocols, and formats",
					"icon":        "forward",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8008,
						"jmxPort":  1108,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), eventSourcesMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return eventSourcesMS, nil
}

func createCRSiteWhereInboundProcessingIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {
	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	inboundProcessingMS, err := res.Get(context.TODO(), "inbound-processing-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		inboundProcessingMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "inbound-processing-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "inbound-processing",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Inbound Processing",
					"description": "Common processing logic applied to enrich and direct inbound events",
					"icon":        "input",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8007,
						"jmxPort":  1107,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), inboundProcessingMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return inboundProcessingMS, nil
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
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 8080,
							},
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       8080,
								"targetPort": 8080,
								"name":       "http-rest",
							},
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
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

func createCRSiteWhereLabelGenerationIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	labelGenerationMS, err := res.Get(context.TODO(), "label-generation-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		labelGenerationMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "label-generation-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "label-generation",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Label Generation",
					"description": "Supports generating labels such as bar codes and QR codes for devices",
					"icon":        "label",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8009,
						"jmxPort":  1109,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), labelGenerationMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return labelGenerationMS, nil
}

func createCRSiteWhereOutboundConnectorsIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	outboundConnectorsMS, err := res.Get(context.TODO(), "outbound-connectors-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		outboundConnectorsMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "outbound-connectors-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "outbound-connectors",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Outbound Connectors",
					"description": "Allows event streams to be delivered to external systems for additional processing",
					"icon":        "label",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8016,
						"jmxPort":  1116,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), outboundConnectorsMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return outboundConnectorsMS, nil
}

func createCRSiteWhereScheduleManagementIfNotExists(instance *alpha3.SiteWhereInstance, namespace string, client dynamic.Interface) (*unstructured.Unstructured, error) {

	res := client.Resource(sitewhereMicroserviceGVR).Namespace(namespace)

	outboundConnectorsMS, err := res.Get(context.TODO(), "schedule-management-microservice", metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		outboundConnectorsMS = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "SiteWhereMicroservice",
				"apiVersion": sitewhereMicroserviceGVR.Group + "/" + sitewhereMicroserviceGVR.Version,
				"metadata": map[string]interface{}{
					"name":      "schedule-management-microservice",
					"namespace": namespace,
					"labels": map[string]interface{}{
						"sitewhere.io/instance":        instance.Name,
						"sitewhere.io/functional-area": "schedule-management",
					},
				},
				"spec": map[string]interface{}{
					"replicas":    1, // TODO from parameter
					"multitenant": true,
					"name":        "Schedule Management",
					"description": "Supports scheduling of various system operations",
					"icon":        "label",
					"helm": map[string]interface{}{ // TODO Remove when operatior udpates to not using helm
						"chartName":      "sitewhere-0.3.0",
						"releaseName":    instance.Name,
						"releaseService": "Tiller",
					},
					"podSpec": map[string]interface{}{
						"imageRegistry":   "docker.io",
						"imageRepository": "sitewhere",
						"imageTag":        "3.0.0.beta1", // TODO from parameter
						"imagePullPolicy": "IfNotPresent",
						"ports": []map[string]interface{}{
							map[string]interface{}{
								"containerPort": 9000,
							},
							map[string]interface{}{
								"containerPort": 9090,
							},
						},
						"env": []map[string]interface{}{
							map[string]interface{}{
								"name": "sitewhere.config.k8s.name",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.name",
									},
								},
							},
							map[string]interface{}{
								"name": "sitewhere.config.k8s.namespace",
								"valueFrom": map[string]interface{}{
									"fieldRef": map[string]interface{}{
										"fieldPath": "metadata.namespace",
									},
								},
							},
							map[string]interface{}{
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
							map[string]interface{}{
								"port":       9000,
								"targetPort": 9000,
								"name":       "grpc-api",
							},
							map[string]interface{}{
								"port":       9090,
								"targetPort": 9090,
								"name":       "http-metrics",
							},
						},
					},
					"debug": map[string]interface{}{
						"enabled":  false,
						"jdwpPort": 8018,
						"jmxPort":  1118,
					},
				},
			},
		}

		result, err := res.Create(context.TODO(), outboundConnectorsMS, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		return result, err
	}
	return outboundConnectorsMS, nil
}
