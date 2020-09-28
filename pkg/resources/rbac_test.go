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

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateClusterRoleIfNotExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		cr        rbacv1.ClusterRole
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRole exists, should return existing
		{
			cr: rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRole does not exist, should return created ns
		{
			cr: rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.cr.ObjectMeta.Name, func(single struct {
			cr        rbacv1.ClusterRole
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
		depl      rbacv1.ClusterRole
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRole exists, should return existing
		{
			depl: rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRole does not exist, should return created ns
		{
			depl: rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("clusterroles.rbac.authorization.k8s.io \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.depl.ObjectMeta.Name, func(single struct {
			depl      rbacv1.ClusterRole
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
		crb       rbacv1.ClusterRoleBinding
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRoleBinding exists, should return existing
		{
			crb: rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRoleBinding does not exist, should return created ns
		{
			crb: rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
		},
	}
	for _, single := range data {
		t.Run(single.crb.ObjectMeta.Name, func(single struct {
			crb       rbacv1.ClusterRoleBinding
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
		depl      rbacv1.ClusterRoleBinding
		clientset kubernetes.Interface
		err       error
	}{
		// ClusterRoleBinding exists, should return existing
		{
			depl: rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(&rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Annotations: map[string]string{},
				},
			}),
		},
		// ClusterRoleBinding does not exist, should return created ns
		{
			depl: rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "non-existing",
				Annotations: map[string]string{},
			}},
			clientset: fake.NewSimpleClientset(),
			err:       fmt.Errorf("clusterrolebindings.rbac.authorization.k8s.io \"non-existing\" not found"),
		},
	}

	for _, single := range data {
		t.Run(single.depl.ObjectMeta.Name, func(single struct {
			depl      rbacv1.ClusterRoleBinding
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
		role      rbacv1.Role
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Role exists, should return existing
		{
			role: rbacv1.Role{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacv1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Role does not exist, should return created ns
		{
			role: rbacv1.Role{ObjectMeta: metav1.ObjectMeta{
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
			role      rbacv1.Role
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
		role      rbacv1.Role
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Role exists, should return existing
		{
			role: rbacv1.Role{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacv1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// Role does not exist, should return created ns
		{
			role: rbacv1.Role{ObjectMeta: metav1.ObjectMeta{
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
			role      rbacv1.Role
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
		rb        rbacv1.RoleBinding
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// RoleBinding exists, should return existing
		{
			rb: rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// RoleBinding does not exist, should return created ns
		{
			rb: rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
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
			rb        rbacv1.RoleBinding
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
		rb        rbacv1.RoleBinding
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// RoleBinding exists, should return existing
		{
			rb: rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
				Name:        "existing",
				Namespace:   "ns",
				Annotations: map[string]string{},
			}},
			namespace: "ns",
			clientset: fake.NewSimpleClientset(&rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "existing",
					Namespace:   "ns",
					Annotations: map[string]string{},
				},
			}),
		},
		// RoleBinding does not exist, should return created ns
		{
			rb: rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
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
			rb        rbacv1.RoleBinding
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
