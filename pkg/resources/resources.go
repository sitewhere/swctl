/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	kubernetes "k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/api/errors"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	deployRunningThreshold     = time.Minute * 10 // Max wait time
	deployRunningCheckInterval = time.Second * 2
)

// InstallResourceFromFile Install a resource from a file name
func InstallResourceFromFile(fileName string,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	r, err := statikFS.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	sch := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(sch)
	_ = apiextv1beta1.AddToScheme(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, groupVersionKind, err := decode([]byte(contents), nil, nil)

	_ = groupVersionKind

	if err != nil {
		// If we can decode, try installing custom resource
		return CreateCustomResourceFromFile(fileName, statikFS, config)
	}

	// now use switch over the type of the object
	// and match each type-case
	switch o := obj.(type) {
	case *v1.Pod:
		_, err = CreatePodIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *v1.ConfigMap:
		_, err = CreateConfigMapIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *v1.Secret:
		_, err = CreateSecretIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *v1.ServiceAccount:
		_, err = CreateServiceAccountIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *v1.PersistentVolumeClaim:
		_, err = CreatePersistentVolumeClaimIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *v1.Service:
		_, err = CreateServiceIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *appsv1.Deployment:
		_, err = CreateDeploymentIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *appsv1.StatefulSet:
		_, err = CreateStatefulSetIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *rbacV1.ClusterRole:
		_, err = CreateClusterRoleIfNotExists(o, clientset)
	case *rbacV1.ClusterRoleBinding:
		_, err = CreateClusterRoleBindingIfNotExists(o, clientset)
	case *rbacV1.Role:
		_, err = CreateRoleIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *rbacV1.RoleBinding:
		_, err = CreateRoleBindingIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *policyV1beta1.PodDisruptionBudget:
		_, err = CreatePodDisruptionBudgetIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *networkingv1.NetworkPolicy:
		_, err = CreateNetworkPolicyIfNotExists(o, clientset, sitewhereSystemNamespace)
	case *apiextv1beta1.CustomResourceDefinition:
		_, err = CreateCustomResourceDefinitionIfNotExists(o, apiextensionsClientset)
	default:
		fmt.Println(fmt.Sprintf("Resource with type %v not handled.", groupVersionKind))
		_ = o //o is unknown for us
	}

	if err != nil {
		return err
	}

	return nil
}

// UninstallResourceFromFile Uninstall a resource from a file name
func UninstallResourceFromFile(fileName string,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	r, err := statikFS.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	sch := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(sch)
	_ = apiextv1beta1.AddToScheme(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, groupVersionKind, err := decode([]byte(contents), nil, nil)

	_ = groupVersionKind

	if err != nil {
		// If we can decode, try uninstalling custom resource
		return DeleteCustomResourceFromFile(fileName, config, statikFS)
	}

	// now use switch over the type of the object
	// and match each type-case
	switch o := obj.(type) {
	case *v1.Pod:
		err = DeletePodIfExists(o, clientset, sitewhereSystemNamespace)
	case *v1.ConfigMap:
		err = DeleteConfigMapIfExists(o, clientset, sitewhereSystemNamespace)
	case *v1.Secret:
		err = DeleteSecretIfExists(o, clientset, sitewhereSystemNamespace)
	case *v1.ServiceAccount:
		err = DeleteServiceAccountIfExists(o, clientset, sitewhereSystemNamespace)
	case *v1.PersistentVolumeClaim:
		err = DeletePersistentVolumeClaimIfExists(o, clientset, sitewhereSystemNamespace)
	case *v1.Service:
		err = DeleteServiceIfExists(o, clientset, sitewhereSystemNamespace)
	case *appsv1.Deployment:
		err = DeleteDeploymentIfExists(o, clientset, sitewhereSystemNamespace)
	case *appsv1.StatefulSet:
		err = DeleteStatefulSetIfExists(o, clientset, sitewhereSystemNamespace)
	case *rbacV1.ClusterRole:
		err = DeleteClusterRoleIfExists(o, clientset)
	case *rbacV1.ClusterRoleBinding:
		err = DeleteClusterRoleBindingIfExists(o, clientset)
	case *rbacV1.Role:
		err = DeleteRoleIfExists(o, clientset, sitewhereSystemNamespace)
	case *rbacV1.RoleBinding:
		err = DeleteRoleBindingIfExists(o, clientset, sitewhereSystemNamespace)
	case *policyV1beta1.PodDisruptionBudget:
		err = DeletePodDisruptionBudgetIfExists(o, clientset, sitewhereSystemNamespace)
	case *networkingv1.NetworkPolicy:
		err = DeleteNetworkPolicyIfExists(o, clientset, sitewhereSystemNamespace)
	case *apiextv1beta1.CustomResourceDefinition:
		apiextensionsClient, err := apiextensionsclientset.NewForConfig(config)
		if err != nil {
			fmt.Printf("Error getting Kubernetes API Extension Client: %v\n", err)
			return err
		}
		err = DeleteCustomResourceDefinitionIfExists(o, apiextensionsClient)
	default:
		fmt.Println(fmt.Sprintf("Resource with type %v not handled.", groupVersionKind))
		_ = o //o is unknown for us
	}

	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func waitForPodContainersRunning(clientset kubernetes.Interface, podName string, namespace string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := podContainersRunning(clientset, podName, namespace)
		if running {
			return nil
		}

		if err != nil && errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for running pods: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get all running containers")
		}
	}
	return nil
}

func podContainersRunning(clientset kubernetes.Interface, podName string, namespace string) (bool, error) {
	existingPod, err := clientset.CoreV1().Pods(namespace).Get(
		context.TODO(),
		podName,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	for _, status := range existingPod.Status.ContainerStatuses {
		if !status.Ready {
			return false, nil
		}
	}

	return true, nil
}

func waitForDeploymentAvailable(clientset kubernetes.Interface, deploymentName string, namespace string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := deploymentAvailable(clientset, deploymentName, namespace)
		if running {
			return nil
		}

		if err != nil && !errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get deployment available")
		}
	}
	return nil
}

func deploymentAvailable(clientset kubernetes.Interface, deploymentName string, namespace string) (bool, error) {
	existingDeploy, err := clientset.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	if existingDeploy.Status.ReadyReplicas < existingDeploy.Status.AvailableReplicas {
		return false, nil
	}

	return true, nil
}
