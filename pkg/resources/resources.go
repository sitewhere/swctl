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
func InstallResourceFromFile(file http.File,
	fileName string,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) (*metav1.ObjectMeta, error) {

	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
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
		pod, err := CreatePodIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &pod.ObjectMeta, err
	case *v1.ConfigMap:
		configMap, err := CreateConfigMapIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &configMap.ObjectMeta, err
	case *v1.Secret:
		secret, err := CreateSecretIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &secret.ObjectMeta, err
	case *v1.ServiceAccount:
		sa, err := CreateServiceAccountIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &sa.ObjectMeta, err
	case *v1.PersistentVolumeClaim:
		pvc, err := CreatePersistentVolumeClaimIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &pvc.ObjectMeta, err
	case *v1.Service:
		svc, err := CreateServiceIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &svc.ObjectMeta, err
	case *appsv1.Deployment:
		deploy, err := CreateDeploymentIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &deploy.ObjectMeta, err
	case *appsv1.StatefulSet:
		ss, err := CreateStatefulSetIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &ss.ObjectMeta, err
	case *v1.Namespace:
		ns, err := CreateNamespaceIfNotExists(SitewhereSystemNamespace(), clientset)
		if err != nil {
			return nil, err
		}
		return &ns.ObjectMeta, err
	case *rbacV1.ClusterRole:
		cr, err := CreateClusterRoleIfNotExists(o, clientset)
		if err != nil {
			return nil, err
		}
		return &cr.ObjectMeta, err
	case *rbacV1.ClusterRoleBinding:
		crb, err := CreateClusterRoleBindingIfNotExists(o, clientset)
		if err != nil {
			return nil, err
		}
		return &crb.ObjectMeta, err
	case *rbacV1.Role:
		role, err := CreateRoleIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &role.ObjectMeta, err
	case *rbacV1.RoleBinding:
		rb, err := CreateRoleBindingIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &rb.ObjectMeta, err
	case *policyV1beta1.PodDisruptionBudget:
		pol, err := CreatePodDisruptionBudgetIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &pol.ObjectMeta, err
	case *networkingv1.NetworkPolicy:
		np, err := CreateNetworkPolicyIfNotExists(o, clientset, SitewhereSystemNamespace())
		if err != nil {
			return nil, err
		}
		return &np.ObjectMeta, err
	case *apiextv1beta1.CustomResourceDefinition:
		crd, err := CreateCustomResourceDefinitionIfNotExists(o, apiextensionsClientset)
		if err != nil {
			return nil, err
		}
		return &crd.ObjectMeta, err
	default:
		fmt.Println(fmt.Sprintf("Resource with type %v not handled.", groupVersionKind))
		_ = o //o is unknown for us
		return nil, fmt.Errorf("Resources not handled %v", groupVersionKind)
	}
}

// UninstallResourceFromFile Uninstall a resource from a file name
func UninstallResourceFromFile(file http.File,
	fileName string,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {

	defer file.Close()
	contents, err := ioutil.ReadAll(file)
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
		err = DeletePodIfExists(o, clientset, SitewhereSystemNamespace())
	case *v1.ConfigMap:
		err = DeleteConfigMapIfExists(o, clientset, SitewhereSystemNamespace())
	case *v1.Secret:
		err = DeleteSecretIfExists(o, clientset, SitewhereSystemNamespace())
	case *v1.ServiceAccount:
		err = DeleteServiceAccountIfExists(o, clientset, SitewhereSystemNamespace())
	case *v1.PersistentVolumeClaim:
		err = DeletePersistentVolumeClaimIfExists(o, clientset, SitewhereSystemNamespace())
	case *v1.Service:
		err = DeleteServiceIfExists(o, clientset, SitewhereSystemNamespace())
	case *appsv1.Deployment:
		err = DeleteDeploymentIfExists(o, clientset, SitewhereSystemNamespace())
	case *appsv1.StatefulSet:
		err = DeleteStatefulSetIfExists(o, clientset, SitewhereSystemNamespace())
	case *rbacV1.ClusterRole:
		err = DeleteClusterRoleIfExists(o, clientset)
	case *rbacV1.ClusterRoleBinding:
		err = DeleteClusterRoleBindingIfExists(o, clientset)
	case *rbacV1.Role:
		err = DeleteRoleIfExists(o, clientset, SitewhereSystemNamespace())
	case *rbacV1.RoleBinding:
		err = DeleteRoleBindingIfExists(o, clientset, SitewhereSystemNamespace())
	case *policyV1beta1.PodDisruptionBudget:
		err = DeletePodDisruptionBudgetIfExists(o, clientset, SitewhereSystemNamespace())
	case *networkingv1.NetworkPolicy:
		err = DeleteNetworkPolicyIfExists(o, clientset, SitewhereSystemNamespace())
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

// WaitForDeploymentAvailable waits for a Deployment to became available
func WaitForDeploymentAvailable(clientset kubernetes.Interface, deploymentName string, namespace string) error {
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

// WaitForSecretExists waits for a Secret to exists
func WaitForSecretExists(clientset kubernetes.Interface, name string, namespace string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := secretExists(clientset, name, namespace)
		if running {
			return nil
		}

		if err != nil && !errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get Secret available")
		}
	}
	return nil
}

func secretExists(clientset kubernetes.Interface, name string, namespace string) (bool, error) {
	_, err := clientset.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	return true, nil
}

// WaitForCRDStablished waits for a CRD to be stablished
func WaitForCRDStablished(apiextensionsclientset apiextensionsclientset.Interface, name string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := crdStablished(apiextensionsclientset, name)
		if running {
			return nil
		}

		if err != nil && !errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get Secret available")
		}
	}
	return nil
}

func crdStablished(apiextensionsclientset apiextensionsclientset.Interface, name string) (bool, error) {
	crds := apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions()

	existingCRD, err := crds.Get(context.TODO(),
		name,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	for _, cond := range existingCRD.Status.Conditions {
		if cond.Type == apiextv1beta1.Established {
			return true, nil
		}
	}
	return false, nil
}
