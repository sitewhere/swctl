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
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsFake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
)

func TestWaitForPodContainersRunning(t *testing.T) {
	t.Parallel()
	data := []struct {
		name      string
		podName   string
		namespace string
		clientset kubernetes.Interface
		err       error
	}{
		// Pods exists and is ready, should return no error
		{
			name:      "pod-ready",
			podName:   "existing",
			namespace: "existing-ns",
			clientset: fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "existing",
					Namespace: "existing-ns",
				},
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						v1.ContainerStatus{
							Ready: true,
						},
					},
				},
			}),
			err: nil,
		},
	}
	for _, single := range data {
		t.Run(single.name, func(single struct {
			name      string
			podName   string
			namespace string
			clientset kubernetes.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := waitForPodContainersRunning(single.clientset, single.podName, single.namespace)
				if err != single.err {
					t.Fatalf("expected err: %s got err: %s", single.err, err)
				}
			}
		}(single))
	}
}

func TestWaitForDeploymentAvailable(t *testing.T) {
	t.Parallel()
	data := []struct {
		name       string
		deployName string
		namespace  string
		clientset  kubernetes.Interface
		err        error
	}{
		// Deployment exists and is available, should return no error
		{
			name:       "deploy-available",
			deployName: "existing",
			namespace:  "existing-ns",
			clientset: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "existing",
					Namespace: "existing-ns",
				},
				Status: appsv1.DeploymentStatus{
					ReadyReplicas:     1,
					AvailableReplicas: 1,
					Conditions: []appsv1.DeploymentCondition{
						appsv1.DeploymentCondition{
							Type: appsv1.DeploymentProgressing,
						},
					},
				},
			}),
			err: nil,
		},
	}
	for _, single := range data {
		t.Run(single.name, func(single struct {
			name       string
			deployName string
			namespace  string
			clientset  kubernetes.Interface
			err        error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := WaitForDeploymentAvailable(single.clientset, single.deployName, single.namespace)
				if err != single.err {
					t.Fatalf("expected err: %s got err: %s", single.err, err)
				}
			}
		}(single))
	}
}

func TestWaitForSecretExists(t *testing.T) {
	t.Parallel()
	data := []struct {
		name       string
		secretName string
		namespace  string
		clientset  kubernetes.Interface
		err        error
	}{
		// Secret exists, should return no error
		{
			name:       "secret-exists",
			secretName: "existing",
			namespace:  "existing-ns",
			clientset: fake.NewSimpleClientset(&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "existing",
					Namespace: "existing-ns",
				},
			}),
			err: nil,
		},
	}
	for _, single := range data {
		t.Run(single.name, func(single struct {
			name       string
			secretName string
			namespace  string
			clientset  kubernetes.Interface
			err        error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := WaitForSecretExists(single.clientset, single.secretName, single.namespace)
				if err != single.err {
					t.Fatalf("expected err: %s got err: %s", single.err, err)
				}
			}
		}(single))
	}
}

func TestWaitForCRDStablished(t *testing.T) {
	t.Parallel()
	data := []struct {
		name      string
		crdName   string
		clientset apiextensionsclientset.Interface
		err       error
	}{
		// CRD exists and is stablished, should return no error
		{
			name:    "crd-stablished",
			crdName: "existing",
			clientset: apiextensionsFake.NewSimpleClientset(&apiextv1beta1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: "existing",
				},
				Status: apiextv1beta1.CustomResourceDefinitionStatus{
					Conditions: []apiextv1beta1.CustomResourceDefinitionCondition{
						apiextv1beta1.CustomResourceDefinitionCondition{
							Type: apiextv1beta1.Established,
						},
					},
				},
			}),
			err: nil,
		},
	}
	for _, single := range data {
		t.Run(single.name, func(single struct {
			name      string
			crdName   string
			clientset apiextensionsclientset.Interface
			err       error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := WaitForCRDStablished(single.clientset, single.crdName)
				if err != single.err {
					t.Fatalf("expected err: %s got err: %s", single.err, err)
				}
			}
		}(single))
	}
}
