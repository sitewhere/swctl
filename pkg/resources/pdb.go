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
