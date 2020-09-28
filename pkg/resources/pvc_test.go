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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreatePersistentVolumeClaimIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		pvc       v1.PersistentVolumeClaim
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// PersistentVolumeClaim exists, should return existing
		{
			pvc: v1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// PersistentVolumeClaim does not exist, should return created ns
		{
			pvc: v1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.pvc.ObjectMeta.Name, func(single struct {
			pvc       v1.PersistentVolumeClaim
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreatePersistentVolumeClaimIfNotExists(&single.pvc, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.pvc.ObjectMeta.Name {
						t.Fatalf("expected %s persistentvolumeclaim, got %s", single.pvc.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeletePersistentVolumeClaimIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		pvc       v1.PersistentVolumeClaim
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// PersistentVolumeClaim exists, should return existing
		{
			pvc: v1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// PersistentVolumeClaim does not exist, should return created ns
		{
			pvc: v1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("persistentvolumeclaims \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.pvc.ObjectMeta.Name, func(single struct {
			pvc       v1.PersistentVolumeClaim
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeletePersistentVolumeClaimIfExists(&single.pvc, single.clientset, single.namespace)

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
