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
	"fmt"
	"strings"
	"testing"

	policyV1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreatePodDisruptionBudgetIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		pdb       policyV1beta1.PodDisruptionBudget
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// PodDisruptionBudget exists, should return existing
		{
			pdb: policyV1beta1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&policyV1beta1.PodDisruptionBudget{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// PodDisruptionBudget does not exist, should return created ns
		{
			pdb: policyV1beta1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.pdb.ObjectMeta.Name, func(single struct {
			pdb       policyV1beta1.PodDisruptionBudget
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreatePodDisruptionBudgetIfNotExists(&single.pdb, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.pdb.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.pdb.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeletePodDisruptionBudgetIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		pdb       policyV1beta1.PodDisruptionBudget
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// PodDisruptionBudget exists, should return existing
		{
			pdb: policyV1beta1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&policyV1beta1.PodDisruptionBudget{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// PodDisruptionBudget does not exist, should return created ns
		{
			pdb: policyV1beta1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("poddisruptionbudgets.policy \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.pdb.ObjectMeta.Name, func(single struct {
			pdb       policyV1beta1.PodDisruptionBudget
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeletePodDisruptionBudgetIfExists(&single.pdb, single.clientset, single.namespace)

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
