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

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateCustomResourceDefinitionIfNotExists Create a CustomResourceDefinition if it does not exists.
func CreateCustomResourceDefinitionIfNotExists(crd *apiextv1beta1.CustomResourceDefinition, apiextensionsclientset apiextensionsclientset.Interface) (*apiextv1beta1.CustomResourceDefinition, error) {
	var err error

	crds := apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions()

	_, err = crds.Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	return crd, nil
}

// DeleteCustomResourceDefinitionIfExists Delete a CustomResourceDefinition if it exists
func DeleteCustomResourceDefinitionIfExists(crd *apiextv1beta1.CustomResourceDefinition, apiextensionsclientset apiextensionsclientset.Interface) error {
	return apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(
		context.TODO(),
		crd.ObjectMeta.Name,
		metav1.DeleteOptions{})
}
