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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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
		// Namespaces exists, should return existing
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
