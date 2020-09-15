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
