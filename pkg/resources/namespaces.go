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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	kubernetes "k8s.io/client-go/kubernetes"
)

var (
	sitewhereSystemNamespace = "sitewhere-system"                                              // SiteWhere System Namespace
	decUnstructured          = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme) // Decoding Unstructed
)

// SitewhereSystemNamespace returns the namespace for SiteWhere
func SitewhereSystemNamespace() string {
	return sitewhereSystemNamespace
}

// CreateNamespaceIfNotExists Create a Namespace in Kubernetes if it does not exists.
func CreateNamespaceIfNotExists(namespace string, istioInject bool, clientset kubernetes.Interface) (*v1.Namespace, error) {
	var err error
	var ns *v1.Namespace

	ns, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
		var labels map[string]string
		if istioInject {
			labels = map[string]string{
				"app": namespace,
			}
		} else {
			labels = map[string]string{
				"app":             namespace,
				"istio-injection": "enabled",
			}
		}
		ns = &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   namespace,
				Labels: labels,
			},
		}

		result, err := clientset.CoreV1().Namespaces().Create(context.TODO(),
			ns,
			metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		return result, err
	}

	if err != nil {
		return nil, err
	}

	return ns, nil
}

// DeleteNamespaceIfExists Delete a Namespace in Kubernetes if it does exists.
func DeleteNamespaceIfExists(namespace string, clientset kubernetes.Interface) error {
	return clientset.CoreV1().Namespaces().Delete(context.TODO(),
		namespace,
		metav1.DeleteOptions{})
}

// DeleteSiteWhereNamespaceIfExists Delete a Namespace in Kubernetes if it does exists.
func DeleteSiteWhereNamespaceIfExists(clientset kubernetes.Interface) error {
	return DeleteNamespaceIfExists(SitewhereSystemNamespace(), clientset)
}
