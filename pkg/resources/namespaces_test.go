/**
 * Copyright © 2014-2021 The SiteWhere Authors
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
	"fmt"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateNamespaceIfNotExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Namespaces exists, should return existing
		{
			namespace: "existing",
			clientset: fake.NewSimpleClientset(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// Namespaces does not exist, should return created ns
		{
			namespace: "non-existing",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.namespace, func(single struct {
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateNamespaceIfNotExists(single.namespace, false, single.clientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.namespace {
						t.Fatalf("expected %s namespace, got %s", single.namespace, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteNamespaceIfExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Namespaces exists, should return existing
		{
			namespace: "existing",
			clientset: fake.NewSimpleClientset(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// Namespaces does not exist, should return created ns
		{
			namespace: "non-existing",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("namespaces \"non-existing\" not found"),
		},
	}
	for _, single := range data {
		t.Run(single.namespace, func(single struct {
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteNamespaceIfExists(single.namespace, single.clientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				}
			}
		}(single))
	}
}

func TestDeleteSiteWhereNamespaceIfExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Namespaces exists, should return existing
		{
			namespace: "sitewhere-system",
			clientset: fake.NewSimpleClientset(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "sitewhere-system",
					Annotations: map[string]string{},
				},
			}),
		},
		// Namespaces does not exist, should return created ns
		{
			namespace: "sitewhere-system",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("namespaces \"sitewhere-system\" not found"),
		},
	}
	for _, single := range data {
		t.Run(single.namespace, func(single struct {
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteSiteWhereNamespaceIfExists(single.clientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				}
			}
		}(single))
	}
}

func TestCheckIfExistsNamespace(t *testing.T) {

	t.Parallel()
	data := []struct {
		namespace string
		clientset kubernetes.Interface
		result    bool
		err       error
	}{
		// Namespaces exists, should return existing
		{
			namespace: "existing",
			clientset: fake.NewSimpleClientset(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
			result: true,
			err:    nil,
		},
		// Namespaces does not exist, should return created ns
		{
			namespace: "non-existing",
			clientset: fake.NewSimpleClientset(),
			result:    false,
			err:       nil,
		},
	}
	for _, single := range data {
		t.Run(single.namespace, func(single struct {
			namespace string
			clientset kubernetes.Interface
			result    bool
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CheckIfExistsNamespace(single.namespace, single.clientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result != single.result {
						t.Fatalf("expected %t, got %t", single.result, result)
					}
				}
			}
		}(single))
	}
}
