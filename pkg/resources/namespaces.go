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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	kubernetes "k8s.io/client-go/kubernetes"
)

var (
	sitewhereSystemNamespace = "sitewhere-system"                                              // SiteWhere System Namespace
	decUnstructured          = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme) // Decoding Unstructed
)

// CreateNamespaceIfNotExists Create a Namespace in Kubernetes if it does not exists.
func CreateNamespaceIfNotExists(namespace string, clientset kubernetes.Interface) (*v1.Namespace, error) {
	var err error
	var ns *v1.Namespace

	ns, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if err != nil && errors.IsNotFound(err) {
		ns = &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
				Labels: map[string]string{
					"app": namespace,
				},
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
	return DeleteNamespaceIfExists(sitewhereSystemNamespace, clientset)
}
