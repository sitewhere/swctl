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

// CreateServiceIfNotExists Create a Service Account if it does not exists.
func CreateServiceIfNotExists(svc *v1.Service, clientset kubernetes.Interface, namespace string) (*v1.Service, error) {
	var err error
	var existingSVC *v1.Service

	existingSVC, err = clientset.CoreV1().Services(namespace).Get(
		context.TODO(),
		svc.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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
