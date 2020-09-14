/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
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
