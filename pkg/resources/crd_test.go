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
	"strings"
	"testing"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	extensionFake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateCustomResourceDefinitionIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		crd                    apiextv1beta1.CustomResourceDefinition
		apiextensionsclientset apiextensionsclientset.Interface
		err                    error
	}{
		// CustomResourceDefinition exists, should return existing
		{
			crd: apiextv1beta1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
				Spec: apiextv1beta1.CustomResourceDefinitionSpec{
					Group:   "acme.com",
					Version: "v1alpha1",
					Names: apiextv1beta1.CustomResourceDefinitionNames{
						Kind: "Test",
					},
				},
			},
			apiextensionsclientset: extensionFake.NewSimpleClientset(&apiextv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}}),
		},
		// CustomResourceDefinition does not exist, should return created ns
		{
			crd: apiextv1beta1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Annotations: map[string]string{},
			}},
			apiextensionsclientset: extensionFake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.crd.ObjectMeta.Name, func(single struct {
			crd                    apiextv1beta1.CustomResourceDefinition
			apiextensionsclientset apiextensionsclientset.Interface
			err                    error
		}) func(t *testing.T) {
			return func(t *testing.T) {

				result, err := CreateCustomResourceDefinitionIfNotExists(&single.crd, single.apiextensionsclientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.crd.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.crd.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}
