/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"context"

	v1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

// CreateServiceAccountIfNotExists Create a Service Account if it does not exists.
func CreateServiceAccountIfNotExists(sa *v1.ServiceAccount, clientset kubernetes.Interface, namespace string) (*v1.ServiceAccount, error) {
	var err error
	var existingSA *v1.ServiceAccount

	existingSA, err = clientset.CoreV1().ServiceAccounts(namespace).Get(
		context.TODO(),
		sa.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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
