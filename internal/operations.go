/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package internal Implements swctl internal use only functions
package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	discovery "k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached/memory"
	dynamic "k8s.io/client-go/dynamic"
	kubernetes "k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
	restmapper "k8s.io/client-go/restmapper"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	sitewhereSystemNamespace = "sitewhere-system"                                              // SiteWhere System Namespace
	decUnstructured          = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme) // Decoding Unstructed
)

const (
	deployRunningThreshold     = time.Minute * 10 // Max wait time
	deployRunningCheckInterval = time.Second * 2
)

// InstallResourceFromFile Install a resource from a file name
func InstallResourceFromFile(fileName string, config SiteWhereConfiguration) error {
	r, err := config.GetStatikFS().Open(fileName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", fileName, err)
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content of file%s: %v\n", fileName, err)
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
		return CreateCustomResourceFromFile(fileName, config.GetConfig(), config.GetStatikFS())
	}

	clientset := config.GetClientset()

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
	case *apiextv1beta1.CustomResourceDefinition:
		apiextensionsClient, err := apiextensionsclientset.NewForConfig(config.GetConfig())
		if err != nil {
			fmt.Printf("Error getting Kubernetes API Extension Client: %v\n", err)
			return err
		}
		_, err = CreateCustomResourceDefinitionIfNotExists(o, apiextensionsClient)

	default:
		fmt.Println(fmt.Sprintf("Resource with type %v not handled.", groupVersionKind))
		_ = o //o is unknown for us
	}

	if err != nil {
		fmt.Printf("Error Creating Resource: %v\n", err)
		return err
	}

	return nil
}

// UninstallResourceFromFile Uninstall a resource from a file name
func UninstallResourceFromFile(fileName string, config *rest.Config, statikFS http.FileSystem) error {
	r, err := statikFS.Open(fileName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", fileName, err)
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content of file%s: %v\n", fileName, err)
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
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

// CreateNamespaceIfNotExists Create a Namespace in Kubernetes if it does not exists.
func CreateNamespaceIfNotExists(namespace string, clientset kubernetes.Interface) (*v1.Namespace, error) {
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

// DeleteNamespaceIfExists Delete a Namespace in Kubernetes if it does exists.
func DeleteNamespaceIfExists(namespace string, clientset kubernetes.Interface) error {
	return clientset.CoreV1().Namespaces().Delete(context.TODO(),
		namespace,
		metav1.DeleteOptions{})
}

// DeleteSiteWhereNamespaceIfExists Delete a Namespace in Kubernetes if it does exists.
func DeleteSiteWhereNamespaceIfExists(clientset kubernetes.Interface) error {
	return DeleteNamespaceIfExists(sitewhereSystemNamespace, clientset)
}

// CreateServiceAccountIfNotExists Create a Service Account if it does not exists.
func CreateServiceAccountIfNotExists(sa *v1.ServiceAccount, clientset kubernetes.Interface, namespace string) (*v1.ServiceAccount, error) {
	var err error
	var existingSA *v1.ServiceAccount

	existingSA, err = clientset.CoreV1().ServiceAccounts(namespace).Get(
		context.TODO(),
		sa.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.CoreV1().ServiceAccounts(namespace).Create(
			context.TODO(),
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

	return existingSA, nil
}

// DeleteServiceAccountIfExists Delete a Service Account if it exists.
func DeleteServiceAccountIfExists(sa *v1.ServiceAccount, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().ServiceAccounts(namespace).Delete(
		context.TODO(),
		sa.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreatePodIfNotExists Create a Service Account if it does not exists.
func CreatePodIfNotExists(pod *v1.Pod, clientset kubernetes.Interface, namespace string) (*v1.Pod, error) {
	var err error
	var existingPod *v1.Pod

	existingPod, err = clientset.CoreV1().Pods(namespace).Get(
		context.TODO(),
		pod.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.CoreV1().Pods(namespace).Create(
			context.TODO(),
			pod,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingPod, nil
}

// DeletePodIfExists Delete a Service Account if it exists.
func DeletePodIfExists(pod *v1.Pod, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().Pods(namespace).Delete(
		context.TODO(),
		pod.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateConfigMapIfNotExists Create a Service Account if it does not exists.
func CreateConfigMapIfNotExists(cm *v1.ConfigMap, clientset kubernetes.Interface, namespace string) (*v1.ConfigMap, error) {
	var err error
	var existingCM *v1.ConfigMap

	existingCM, err = clientset.CoreV1().ConfigMaps(namespace).Get(
		context.TODO(),
		cm.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.CoreV1().ConfigMaps(namespace).Create(
			context.TODO(),
			cm,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingCM, nil
}

// DeleteConfigMapIfExists Delete a Service Account if it exists.
func DeleteConfigMapIfExists(cm *v1.ConfigMap, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().ConfigMaps(namespace).Delete(
		context.TODO(),
		cm.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateSecretIfNotExists Create a Service Account if it does not exists.
func CreateSecretIfNotExists(sec *v1.Secret, clientset kubernetes.Interface, namespace string) (*v1.Secret, error) {
	var err error
	var existingSec *v1.Secret

	existingSec, err = clientset.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		sec.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.CoreV1().Secrets(namespace).Create(
			context.TODO(),
			sec,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingSec, nil
}

// DeleteSecretIfExists Delete a Service Account if it exists.
func DeleteSecretIfExists(sec *v1.Secret, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().Secrets(namespace).Delete(
		context.TODO(),
		sec.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreatePersistentVolumeClaimIfNotExists Create a Service Account if it does not exists.
func CreatePersistentVolumeClaimIfNotExists(pvc *v1.PersistentVolumeClaim, clientset kubernetes.Interface, namespace string) (*v1.PersistentVolumeClaim, error) {
	var err error
	var existingPVC *v1.PersistentVolumeClaim

	existingPVC, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Get(
		context.TODO(),
		pvc.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(
			context.TODO(),
			pvc,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingPVC, nil
}

// DeletePersistentVolumeClaimIfExists Delete a Service Account if it exists.
func DeletePersistentVolumeClaimIfExists(pvc *v1.PersistentVolumeClaim, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(
		context.TODO(),
		pvc.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateServiceIfNotExists Create a Service Account if it does not exists.
func CreateServiceIfNotExists(svc *v1.Service, clientset kubernetes.Interface, namespace string) (*v1.Service, error) {
	var err error
	var existingSVC *v1.Service

	existingSVC, err = clientset.CoreV1().Services(namespace).Get(
		context.TODO(),
		svc.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.CoreV1().Services(namespace).Create(
			context.TODO(),
			svc,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingSVC, nil
}

// DeleteServiceIfExists Delete a Service Account if it exists.
func DeleteServiceIfExists(svc *v1.Service, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().Services(namespace).Delete(
		context.TODO(),
		svc.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateDeploymentIfNotExists Create a Service Account if it does not exists.
func CreateDeploymentIfNotExists(deploy *appsv1.Deployment, clientset kubernetes.Interface, namespace string) (*appsv1.Deployment, error) {
	var err error
	var existingDeploy *appsv1.Deployment

	existingDeploy, err = clientset.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploy.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.AppsV1().Deployments(namespace).Create(
			context.TODO(),
			deploy,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingDeploy, nil
}

// DeleteDeploymentIfExists Delete a Service Account if it does not exists.
func DeleteDeploymentIfExists(deploy *appsv1.Deployment, clientset kubernetes.Interface, namespace string) error {
	return clientset.AppsV1().Deployments(namespace).Delete(
		context.TODO(),
		deploy.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateStatefulSetIfNotExists Create a Service Account if it does not exists.
func CreateStatefulSetIfNotExists(ss *appsv1.StatefulSet, clientset kubernetes.Interface, namespace string) (*appsv1.StatefulSet, error) {
	var err error
	var existingSS *appsv1.StatefulSet

	existingSS, err = clientset.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		ss.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.AppsV1().StatefulSets(namespace).Create(
			context.TODO(),
			ss,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingSS, nil
}

// DeleteStatefulSetIfExists Delete a Service Account if it exists.
func DeleteStatefulSetIfExists(ss *appsv1.StatefulSet, clientset kubernetes.Interface, namespace string) error {
	return clientset.AppsV1().StatefulSets(namespace).Delete(
		context.TODO(),
		ss.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateClusterRoleIfNotExists Create a ClusterRole if it does not exists.
func CreateClusterRoleIfNotExists(cr *rbacV1.ClusterRole, clientset kubernetes.Interface) (*rbacV1.ClusterRole, error) {
	var err error
	var existingCR *rbacV1.ClusterRole

	existingCR, err = clientset.RbacV1().ClusterRoles().Get(
		context.TODO(),
		cr.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.RbacV1().ClusterRoles().Create(
			context.TODO(),
			cr,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingCR, nil
}

// DeleteClusterRoleIfExists Delete a ClusterRole if it exists.
func DeleteClusterRoleIfExists(cr *rbacV1.ClusterRole, clientset kubernetes.Interface) error {
	return clientset.RbacV1().ClusterRoles().Delete(
		context.TODO(),
		cr.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateClusterRoleBindingIfNotExists Create a ClusterRoleBinding if it does not exists.
func CreateClusterRoleBindingIfNotExists(crb *rbacV1.ClusterRoleBinding, clientset kubernetes.Interface) (*rbacV1.ClusterRoleBinding, error) {
	var err error
	var existingCRB *rbacV1.ClusterRoleBinding

	existingCRB, err = clientset.RbacV1().ClusterRoleBindings().Get(
		context.TODO(),
		crb.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.RbacV1().ClusterRoleBindings().Create(
			context.TODO(),
			crb,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingCRB, nil
}

// DeleteClusterRoleBindingIfExists Delete a ClusterRoleBinding if it exists.
func DeleteClusterRoleBindingIfExists(crb *rbacV1.ClusterRoleBinding, clientset kubernetes.Interface) error {
	return clientset.RbacV1().ClusterRoleBindings().Delete(
		context.TODO(),
		crb.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateRoleIfNotExists Create a Role if it does not exists.
func CreateRoleIfNotExists(role *rbacV1.Role, clientset kubernetes.Interface, namespace string) (*rbacV1.Role, error) {
	var err error
	var existingRole *rbacV1.Role

	existingRole, err = clientset.RbacV1().Roles(namespace).Get(
		context.TODO(),
		role.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.RbacV1().Roles(namespace).Create(
			context.TODO(),
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

	return existingRole, nil
}

// DeleteRoleIfExists Delete a Role if it does not exists.
func DeleteRoleIfExists(role *rbacV1.Role, clientset kubernetes.Interface, namespace string) error {
	return clientset.RbacV1().Roles(namespace).Delete(
		context.TODO(),
		role.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateRoleBindingIfNotExists Create a RoleBinding if it does not exists.
func CreateRoleBindingIfNotExists(rb *rbacV1.RoleBinding, clientset kubernetes.Interface, namespace string) (*rbacV1.RoleBinding, error) {
	var err error
	var existingRoleBinding *rbacV1.RoleBinding

	existingRoleBinding, err = clientset.RbacV1().RoleBindings(namespace).Get(
		context.TODO(),
		rb.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.RbacV1().RoleBindings(namespace).Create(
			context.TODO(),
			rb,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingRoleBinding, nil
}

// DeleteRoleBindingIfExists Delete a RoleBinding if it exists.
func DeleteRoleBindingIfExists(rb *rbacV1.RoleBinding, clientset kubernetes.Interface, namespace string) error {
	return clientset.RbacV1().RoleBindings(namespace).Delete(
		context.TODO(),
		rb.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreatePodDisruptionBudgetIfNotExists Create a PodDisruptionBudget if it does not exists.
func CreatePodDisruptionBudgetIfNotExists(rb *policyV1beta1.PodDisruptionBudget, clientset kubernetes.Interface, namespace string) (*policyV1beta1.PodDisruptionBudget, error) {
	var err error
	var existingPodDisruptionBudget *policyV1beta1.PodDisruptionBudget

	existingPodDisruptionBudget, err = clientset.PolicyV1beta1().PodDisruptionBudgets(namespace).Get(
		context.TODO(),
		rb.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && k8serror.IsNotFound(err) {
		result, err := clientset.PolicyV1beta1().PodDisruptionBudgets(namespace).Create(
			context.TODO(),
			rb,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingPodDisruptionBudget, nil
}

// DeletePodDisruptionBudgetIfExists Delete a PodDisruptionBudget if it exists.
func DeletePodDisruptionBudgetIfExists(rb *policyV1beta1.PodDisruptionBudget, clientset kubernetes.Interface, namespace string) error {
	return clientset.PolicyV1beta1().PodDisruptionBudgets(namespace).Delete(
		context.TODO(),
		rb.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateCustomResourceDefinitionIfNotExists Create a CustomResourceDefinition if it does not exists.
func CreateCustomResourceDefinitionIfNotExists(crd *apiextv1beta1.CustomResourceDefinition, apiextensionsclientset apiextensionsclientset.Interface) (*apiextv1beta1.CustomResourceDefinition, error) {
	var err error

	crds := apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions()

	_, err = crds.Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	return crd, nil
}

// DeleteCustomResourceDefinitionIfExists Delete a CustomResourceDefinition if it exists
func DeleteCustomResourceDefinitionIfExists(crd *apiextv1beta1.CustomResourceDefinition, apiextensionsclientset apiextensionsclientset.Interface) error {
	return apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(
		context.TODO(),
		crd.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateCustomResourceFromFile Reads a File from statik and creates a CustomResource from it.
func CreateCustomResourceFromFile(crName string, config *rest.Config, statikFS http.FileSystem) error {
	r, err := statikFS.Open(crName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crName, err)
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content of %s: %v\n", crName, err)
		return err
	}

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewDiscoveryClientForConfig for %s: %v\n", crName, err)
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewForConfig for %s: %v\n", crName, err)
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(contents), nil, obj)
	if err != nil {
		fmt.Printf("Error decoding for %s: %v\n", crName, err)
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Printf("Error finding GRV for %s: %v\n", crName, err)
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	_, err = dr.Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("Error creating resource from file %s of Kind: %s: %v", crName, gvk.GroupKind().Kind, err)
		}
		return err
	}
	return nil
}

// DeleteCustomResourceFromFile Reads a File from statik and deletes a CustomResource from it.
func DeleteCustomResourceFromFile(crName string, config *rest.Config, statikFS http.FileSystem) error {
	r, err := statikFS.Open(crName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crName, err)
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content: %v\n", err)
		return err
	}

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewDiscoveryClientForConfig: %v\n", err)
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewForConfig: %v\n", err)
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(contents), nil, obj)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	err = dr.Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})

	if err != nil && !errors.IsNotFound(err) {
		fmt.Printf("Error deleting resource from file %s of Kind: %s: %v", crName, gvk.GroupKind().Kind, err)
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

		if err != nil && k8serror.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for running pods: %s", err))
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

		if err != nil && k8serror.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s", err))
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
