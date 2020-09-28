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

	rbacv1 "k8s.io/api/rbac/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

// CreateClusterRoleIfNotExists Create a ClusterRole if it does not exists.
func CreateClusterRoleIfNotExists(cr *rbacv1.ClusterRole, clientset kubernetes.Interface) (*rbacv1.ClusterRole, error) {
	var err error
	var existingCR *rbacv1.ClusterRole

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
func DeleteClusterRoleIfExists(cr *rbacv1.ClusterRole, clientset kubernetes.Interface) error {
	return clientset.RbacV1().ClusterRoles().Delete(
		context.TODO(),
		cr.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateClusterRoleBindingIfNotExists Create a ClusterRoleBinding if it does not exists.
func CreateClusterRoleBindingIfNotExists(crb *rbacv1.ClusterRoleBinding, clientset kubernetes.Interface) (*rbacv1.ClusterRoleBinding, error) {
	var err error
	var existingCRB *rbacv1.ClusterRoleBinding

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
func DeleteClusterRoleBindingIfExists(crb *rbacv1.ClusterRoleBinding, clientset kubernetes.Interface) error {
	return clientset.RbacV1().ClusterRoleBindings().Delete(
		context.TODO(),
		crb.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateRoleIfNotExists Create a Role if it does not exists.
func CreateRoleIfNotExists(role *rbacv1.Role, clientset kubernetes.Interface, namespace string) (*rbacv1.Role, error) {
	var err error
	var existingRole *rbacv1.Role

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
func DeleteRoleIfExists(role *rbacv1.Role, clientset kubernetes.Interface, namespace string) error {
	return clientset.RbacV1().Roles(namespace).Delete(
		context.TODO(),
		role.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateRoleBindingIfNotExists Create a RoleBinding if it does not exists.
func CreateRoleBindingIfNotExists(rb *rbacv1.RoleBinding, clientset kubernetes.Interface, namespace string) (*rbacv1.RoleBinding, error) {
	var err error
	var existingRoleBinding *rbacv1.RoleBinding

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
func DeleteRoleBindingIfExists(rb *rbacv1.RoleBinding, clientset kubernetes.Interface, namespace string) error {
	return clientset.RbacV1().RoleBindings(namespace).Delete(
		context.TODO(),
		rb.ObjectMeta.Name,
		metav1.DeleteOptions{})
}
