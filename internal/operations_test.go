/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package alpha3 defines SiteWhere Structures
package internal

import (
	"fmt"
	"strings"
	"testing"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	extensionFake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
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
				result, err := CreateNamespaceIfNotExists(single.namespace, single.clientset)

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

func TestCreateServiceAccountIfNotExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		sa        v1.ServiceAccount
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Service Account exists, should return existing
		{
			sa: v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Service Account does not exist, should return created ns
		{
			sa: v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.sa.ObjectMeta.Name, func(single struct {
			sa        v1.ServiceAccount
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateServiceAccountIfNotExists(&single.sa, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.sa.ObjectMeta.Name {
						t.Fatalf("expected %s sa, got %s", single.sa.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteServiceAccountIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		sa        v1.ServiceAccount
		namespace string
		clientset kubernetes.Interface
		err       error
	}{

		// Service Account exists, should return existing
		{
			sa: v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Service Account does not exist, should return created ns
		{
			sa: v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("serviceaccounts \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.sa.ObjectMeta.Name, func(single struct {
			sa        v1.ServiceAccount
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteServiceAccountIfExists(&single.sa, single.clientset, single.namespace)

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

func TestCreatePodIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		pod       v1.Pod
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Pod exists, should return existing
		{
			pod: v1.Pod{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Pod does not exist, should return created ns
		{
			pod: v1.Pod{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.pod.ObjectMeta.Name, func(single struct {
			pod       v1.Pod
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreatePodIfNotExists(&single.pod, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.pod.ObjectMeta.Name {
						t.Fatalf("expected %s pod, got %s", single.pod.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeletePodIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		pod       v1.Pod
		namespace string
		clientset kubernetes.Interface
		err       error
	}{

		// Pod exists, should return existing
		{
			pod: v1.Pod{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Pod does not exist, should return created ns
		{
			pod: v1.Pod{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("pods \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.pod.ObjectMeta.Name, func(single struct {
			pod       v1.Pod
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeletePodIfExists(&single.pod, single.clientset, single.namespace)

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

func TestCreateConfigMapIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		cm        v1.ConfigMap
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// ConfigMap exists, should return existing
		{
			cm: v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// ConfigMap does not exist, should return created ns
		{
			cm: v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.cm.ObjectMeta.Name, func(single struct {
			cm        v1.ConfigMap
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateConfigMapIfNotExists(&single.cm, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.cm.ObjectMeta.Name {
						t.Fatalf("expected %s configmap, got %s", single.cm.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteConfigMapIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		cm        v1.ConfigMap
		namespace string
		clientset kubernetes.Interface
		err       error
	}{

		// ConfigMap exists, should return existing
		{
			cm: v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// ConfigMap does not exist, should return created ns
		{
			cm: v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("configmaps \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.cm.ObjectMeta.Name, func(single struct {
			cm        v1.ConfigMap
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteConfigMapIfExists(&single.cm, single.clientset, single.namespace)

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

func TestCreateSecretIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		sec       v1.Secret
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Secret exists, should return existing
		{
			sec: v1.Secret{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Secret does not exist, should return created ns
		{
			sec: v1.Secret{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.sec.ObjectMeta.Name, func(single struct {
			sec       v1.Secret
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateSecretIfNotExists(&single.sec, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.sec.ObjectMeta.Name {
						t.Fatalf("expected %s secrets, got %s", single.sec.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteSecretIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		sec       v1.Secret
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Secret exists, should return existing
		{
			sec: v1.Secret{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Secret does not exist, should return created ns
		{
			sec: v1.Secret{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("secrets \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.sec.ObjectMeta.Name, func(single struct {
			sec       v1.Secret
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteSecretIfExists(&single.sec, single.clientset, single.namespace)

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

func TestCreateServiceIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		svc       v1.Service
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Service exists, should return existing
		{
			svc: v1.Service{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Service does not exist, should return created ns
		{
			svc: v1.Service{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.svc.ObjectMeta.Name, func(single struct {
			svc       v1.Service
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateServiceIfNotExists(&single.svc, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.svc.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.svc.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteServiceIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		svc       v1.Service
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Service exists, should return existing
		{
			svc: v1.Service{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Service does not exist, should return created ns
		{
			svc: v1.Service{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("services \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.svc.ObjectMeta.Name, func(single struct {
			svc       v1.Service
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteServiceIfExists(&single.svc, single.clientset, single.namespace)

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

func TestCreateDeploymentIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		svc       appsv1.Deployment
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Deployment exists, should return existing
		{
			svc: appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Deployment does not exist, should return created ns
		{
			svc: appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.svc.ObjectMeta.Name, func(single struct {
			svc       appsv1.Deployment
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateDeploymentIfNotExists(&single.svc, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.svc.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.svc.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteDeploymentIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		depl      appsv1.Deployment
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Deployment exists, should return existing
		{
			depl: appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Deployment does not exist, should return created ns
		{
			depl: appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("deployments.apps \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.depl.ObjectMeta.Name, func(single struct {
			depl      appsv1.Deployment
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteDeploymentIfExists(&single.depl, single.clientset, single.namespace)

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

func TestCreateStatefulSetIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		svc       appsv1.StatefulSet
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// StatefulSet exists, should return existing
		{
			svc: appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// StatefulSet does not exist, should return created ns
		{
			svc: appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.svc.ObjectMeta.Name, func(single struct {
			svc       appsv1.StatefulSet
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateStatefulSetIfNotExists(&single.svc, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.svc.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.svc.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteStatefulSetIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		depl      appsv1.StatefulSet
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// StatefulSet exists, should return existing
		{
			depl: appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// StatefulSet does not exist, should return created ns
		{
			depl: appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("statefulsets.apps \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.depl.ObjectMeta.Name, func(single struct {
			depl      appsv1.StatefulSet
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteStatefulSetIfExists(&single.depl, single.clientset, single.namespace)

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

func TestCreateClusterRoleIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		cr        rbacV1.ClusterRole
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRole exists, should return existing
		{
			cr: rbacV1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacV1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRole does not exist, should return created ns
		{
			cr: rbacV1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.cr.ObjectMeta.Name, func(single struct {
			cr        rbacV1.ClusterRole
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateClusterRoleIfNotExists(&single.cr, single.clientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.cr.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.cr.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteClusterRoleIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		depl      rbacV1.ClusterRole
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRole exists, should return existing
		{
			depl: rbacV1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacV1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRole does not exist, should return created ns
		{
			depl: rbacV1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("clusterroles.rbac.authorization.k8s.io \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.depl.ObjectMeta.Name, func(single struct {
			depl      rbacV1.ClusterRole
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteClusterRoleIfExists(&single.depl, single.clientset)

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

func TestCreateClusterRoleBindingIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		crb       rbacV1.ClusterRoleBinding
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRoleBinding exists, should return existing
		{
			crb: rbacV1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacV1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRoleBinding does not exist, should return created ns
		{
			crb: rbacV1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.crb.ObjectMeta.Name, func(single struct {
			crb       rbacV1.ClusterRoleBinding
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateClusterRoleBindingIfNotExists(&single.crb, single.clientset)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.crb.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.crb.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteClusterRoleBindingIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		depl      rbacV1.ClusterRoleBinding
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRoleBinding exists, should return existing
		{
			depl: rbacV1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacV1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRoleBinding does not exist, should return created ns
		{
			depl: rbacV1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("clusterrolebindings.rbac.authorization.k8s.io \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.depl.ObjectMeta.Name, func(single struct {
			depl      rbacV1.ClusterRoleBinding
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteClusterRoleBindingIfExists(&single.depl, single.clientset)

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

func TestCreateRoleIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		role      rbacV1.Role
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Role exists, should return existing
		{
			role: rbacV1.Role{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacV1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Role does not exist, should return created ns
		{
			role: rbacV1.Role{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.role.ObjectMeta.Name, func(single struct {
			role      rbacV1.Role
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateRoleIfNotExists(&single.role, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.role.ObjectMeta.Name {
						t.Fatalf("expected %s role, got %s", single.role.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteRoleIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		role      rbacV1.Role
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Role exists, should return existing
		{
			role: rbacV1.Role{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacV1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Role does not exist, should return created ns
		{
			role: rbacV1.Role{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("roles.rbac.authorization.k8s.io \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.role.ObjectMeta.Name, func(single struct {
			role      rbacV1.Role
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteRoleIfExists(&single.role, single.clientset, single.namespace)

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

func TestCreateRoleBindingIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		rb        rbacV1.RoleBinding
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// RoleBinding exists, should return existing
		{
			rb: rbacV1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacV1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// RoleBinding does not exist, should return created ns
		{
			rb: rbacV1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.rb.ObjectMeta.Name, func(single struct {
			rb        rbacV1.RoleBinding
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := CreateRoleBindingIfNotExists(&single.rb, single.clientset, single.namespace)

				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if result.ObjectMeta.Name != single.rb.ObjectMeta.Name {
						t.Fatalf("expected %s service, got %s", single.rb.ObjectMeta.Name, result.ObjectMeta.Name)
					}
				}
			}
		}(single))
	}
}

func TestDeleteRoleBindingIfExists(t *testing.T) {

	t.Parallel()
	data := []struct {
		rb        rbacV1.RoleBinding
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// RoleBinding exists, should return existing
		{
			rb: rbacV1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacV1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// RoleBinding does not exist, should return created ns
		{
			rb: rbacV1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("rolebindings.rbac.authorization.k8s.io \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.rb.ObjectMeta.Name, func(single struct {
			rb        rbacV1.RoleBinding
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := DeleteRoleBindingIfExists(&single.rb, single.clientset, single.namespace)

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

//"github.com/wking/fakefs"
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
