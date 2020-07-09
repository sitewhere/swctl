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
