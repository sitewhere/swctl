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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateNamespaceIfNotExists Create a Namespace in Kubernetes if it does not exists.
func CreateNamespaceIfNotExists(namespace string, clientset *kubernetes.Clientset) (*v1.Namespace, error) {
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

// CreateServiceAccountIfNotExists Create a Service Account if it does not exists.
func CreateServiceAccountIfNotExists(sa *v1.ServiceAccount, clientset *kubernetes.Clientset, namespace string) (*v1.ServiceAccount, error) {
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

// CreatePodIfNotExists Create a Service Account if it does not exists.
func CreatePodIfNotExists(pod *v1.Pod, clientset *kubernetes.Clientset, namespace string) (*v1.Pod, error) {
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

// CreateConfigMapIfNotExists Create a Service Account if it does not exists.
func CreateConfigMapIfNotExists(cm *v1.ConfigMap, clientset *kubernetes.Clientset, namespace string) (*v1.ConfigMap, error) {
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

// CreateSecretIfNotExists Create a Service Account if it does not exists.
func CreateSecretIfNotExists(sec *v1.Secret, clientset *kubernetes.Clientset, namespace string) (*v1.Secret, error) {
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

// CreatePersistentVolumeClaimIfNotExists Create a Service Account if it does not exists.
func CreatePersistentVolumeClaimIfNotExists(pvc *v1.PersistentVolumeClaim, clientset *kubernetes.Clientset, namespace string) (*v1.PersistentVolumeClaim, error) {
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

// CreateServiceIfNotExists Create a Service Account if it does not exists.
func CreateServiceIfNotExists(svc *v1.Service, clientset *kubernetes.Clientset, namespace string) (*v1.Service, error) {
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

// CreateDeploymentIfNotExists Create a Service Account if it does not exists.
func CreateDeploymentIfNotExists(deploy *appsv1.Deployment, clientset *kubernetes.Clientset, namespace string) (*appsv1.Deployment, error) {
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

// CreateStatefulSetIfNotExists Create a Service Account if it does not exists.
func CreateStatefulSetIfNotExists(ss *appsv1.StatefulSet, clientset *kubernetes.Clientset, namespace string) (*appsv1.StatefulSet, error) {
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

// CreateClusterRoleIfNotExists Create a ClusterRole if it does not exists.
func CreateClusterRoleIfNotExists(cr *rbacV1.ClusterRole, clientset *kubernetes.Clientset) (*rbacV1.ClusterRole, error) {
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

// CreateClusterRoleBindingIfNotExists Create a ClusterRoleBinding if it does not exists.
func CreateClusterRoleBindingIfNotExists(crb *rbacV1.ClusterRoleBinding, clientset *kubernetes.Clientset) (*rbacV1.ClusterRoleBinding, error) {
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

// CreateRoleIfNotExists Create a Role if it does not exists.
func CreateRoleIfNotExists(role *rbacV1.Role, clientset *kubernetes.Clientset, namespace string) (*rbacV1.Role, error) {
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

// CreateRoleBindingIfNotExists Create a RoleBinding if it does not exists.
func CreateRoleBindingIfNotExists(rb *rbacV1.RoleBinding, clientset *kubernetes.Clientset, namespace string) (*rbacV1.RoleBinding, error) {
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

// CreatePodDisruptionBudgetIfNotExists Create a PodDisruptionBudget if it does not exists.
func CreatePodDisruptionBudgetIfNotExists(rb *policyV1beta1.PodDisruptionBudget, clientset *kubernetes.Clientset, namespace string) (*policyV1beta1.PodDisruptionBudget, error) {
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
