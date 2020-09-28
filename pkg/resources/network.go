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

	kubernetes "k8s.io/client-go/kubernetes"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNetworkPolicyIfNotExists Create a CustomResourceDefinition if it does not exists.
func CreateNetworkPolicyIfNotExists(np *networkingv1.NetworkPolicy, clientset kubernetes.Interface, namespace string) (*networkingv1.NetworkPolicy, error) {
	var err error
	var existingNP *networkingv1.NetworkPolicy

	existingNP, err = clientset.NetworkingV1().NetworkPolicies(namespace).Get(
		context.TODO(),
		np.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
		result, err := clientset.NetworkingV1().NetworkPolicies(namespace).Create(
			context.TODO(),
			np,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingNP, nil
}

// DeleteNetworkPolicyIfExists Delete a NetworkPolicy if it exists
func DeleteNetworkPolicyIfExists(np *networkingv1.NetworkPolicy, clientset kubernetes.Interface, namespace string) error {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Delete(
		context.TODO(),
		np.ObjectMeta.Name,
		metav1.DeleteOptions{})
}
