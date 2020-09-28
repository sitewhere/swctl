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

// CreatePersistentVolumeClaimIfNotExists Create a Service Account if it does not exists.
func CreatePersistentVolumeClaimIfNotExists(pvc *v1.PersistentVolumeClaim, clientset kubernetes.Interface, namespace string) (*v1.PersistentVolumeClaim, error) {
	var err error
	var existingPVC *v1.PersistentVolumeClaim

	existingPVC, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Get(
		context.TODO(),
		pvc.ObjectMeta.Name,
		metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
		result, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(
			context.TODO(),
			pvc,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return existingPVC, nil
}

// DeletePersistentVolumeClaimIfExists Delete a Service Account if it exists.
func DeletePersistentVolumeClaimIfExists(pvc *v1.PersistentVolumeClaim, clientset kubernetes.Interface, namespace string) error {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(
		context.TODO(),
		pvc.ObjectMeta.Name,
		metav1.DeleteOptions{})
}
