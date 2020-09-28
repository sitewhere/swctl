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
