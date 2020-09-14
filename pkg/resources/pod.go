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

// CreatePodIfNotExists Create a Service Account if it does not exists.
func CreatePodIfNotExists(pod *v1.Pod, clientset kubernetes.Interface, namespace string) (*v1.Pod, error) {
	var err error
	var existingPod *v1.Pod

	existingPod, err = clientset.CoreV1().Pods(namespace).Get(
		context.TODO(),
		pod.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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
