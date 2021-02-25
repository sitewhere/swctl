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

package config

import (
	"strings"
	"testing"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
)

func TestFromTemplate(t *testing.T) {
	t.Parallel()
	data := []struct {
		name            string
		templateContent string
		placeHolder     *PlaceHolder
		expected        *Configuration
		err             error
	}{
		{
			name: "full-template",
			templateContent: `microservices:
- functionalarea: some-area
  podspec:
    dockerspec:
      tag: "{{ .Tag }}"
`,
			placeHolder: &PlaceHolder{
				Tag: "some-tag",
			},
			expected: &Configuration{
				Microservices: []sitewhereiov1alpha4.SiteWhereMicroserviceSpec{
					{
						FunctionalArea: "some-area",
						PodSpec: &sitewhereiov1alpha4.MicroservicePodSpecification{
							DockerSpec: &sitewhereiov1alpha4.DockerSpec{
								Tag: "some-tag",
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, single := range data {
		t.Run(single.name, func(single struct {
			name            string
			templateContent string
			placeHolder     *PlaceHolder
			expected        *Configuration
			err             error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				result, err := FromTemplate(single.templateContent, single.placeHolder)
				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if len(result.Microservices) != len(single.expected.Microservices) {
						t.Fatalf("expected %v, got %v", single.expected, result)
					}
					for i, ms := range result.Microservices {
						expectedMs := single.expected.Microservices[i]
						if ms.FunctionalArea != expectedMs.FunctionalArea {
							t.Fatalf("expected Funtional Ares %s, got %s", expectedMs.FunctionalArea, ms.FunctionalArea)
						}
						if ms.PodSpec.DockerSpec.Tag != expectedMs.PodSpec.DockerSpec.Tag {
							t.Fatalf("expected Docker Tag %s, got %s", expectedMs.PodSpec.DockerSpec.Tag, ms.PodSpec.DockerSpec.Tag)
						}
					}
				}
			}
		}(single))
	}
}
