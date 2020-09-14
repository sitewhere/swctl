/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"context"

	policyV1beta1 "k8s.io/api/policy/v1beta1"
	kubernetes "k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreatePodDisruptionBudgetIfNotExists Create a PodDisruptionBudget if it does not exists.
func CreatePodDisruptionBudgetIfNotExists(rb *policyV1beta1.PodDisruptionBudget, clientset kubernetes.Interface, namespace string) (*policyV1beta1.PodDisruptionBudget, error) {
	var err error
	var existingPodDisruptionBudget *policyV1beta1.PodDisruptionBudget

	existingPodDisruptionBudget, err = clientset.PolicyV1beta1().PodDisruptionBudgets(namespace).Get(
		context.TODO(),
		rb.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
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
