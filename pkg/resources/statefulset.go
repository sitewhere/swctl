/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

// CreateStatefulSetIfNotExists Create a Service Account if it does not exists.
func CreateStatefulSetIfNotExists(ss *appsv1.StatefulSet, clientset kubernetes.Interface, namespace string) (*appsv1.StatefulSet, error) {
	var err error
	var existingSS *appsv1.StatefulSet

	existingSS, err = clientset.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		ss.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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
