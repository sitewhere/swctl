/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"context"

	rbacV1 "k8s.io/api/rbac/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

// CreateClusterRoleIfNotExists Create a ClusterRole if it does not exists.
func CreateClusterRoleIfNotExists(cr *rbacV1.ClusterRole, clientset kubernetes.Interface) (*rbacV1.ClusterRole, error) {
	var err error
	var existingCR *rbacV1.ClusterRole

	existingCR, err = clientset.RbacV1().ClusterRoles().Get(
		context.TODO(),
		cr.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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

	if err != nil && errors.IsNotFound(err) {
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

	if err != nil && errors.IsNotFound(err) {
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

	if err != nil && errors.IsNotFound(err) {
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
