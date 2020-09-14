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

// CreateDeploymentIfNotExists Create a Service Account if it does not exists.
func CreateDeploymentIfNotExists(deploy *appsv1.Deployment, clientset kubernetes.Interface, namespace string) (*appsv1.Deployment, error) {
	var err error
	var existingDeploy *appsv1.Deployment

	existingDeploy, err = clientset.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploy.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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
