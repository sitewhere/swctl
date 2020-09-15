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

// CreateSecretIfNotExists Create a Service Account if it does not exists.
func CreateSecretIfNotExists(sec *v1.Secret, clientset kubernetes.Interface, namespace string) (*v1.Secret, error) {
	var err error
	var existingSec *v1.Secret

	existingSec, err = clientset.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		sec.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
		result, err := clientset.CoreV1().Secrets(namespace).Create(
			context.TODO(),
			sec,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingSec, nil
}

// DeleteSecretIfExists Delete a Service Account if it exists.
func DeleteSecretIfExists(sec *v1.Secret, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().Secrets(namespace).Delete(
		context.TODO(),
		sec.ObjectMeta.Name,
		metav1.DeleteOptions{})
}
