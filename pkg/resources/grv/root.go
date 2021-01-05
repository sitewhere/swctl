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

package grv

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	sitewhereInstanceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha4",
		Resource: "instances",
	}
)

var (
	sitewhereMicroserviceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha4",
		Resource: "microservices",
	}
)

var (
	sitewhereTenantGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha4",
		Resource: "tenants",
	}
)

// SiteWhereInstanceGRV SiteWhere Intance Group Version Resource
func SiteWhereInstanceGRV() schema.GroupVersionResource {
	return sitewhereInstanceGVR
}

// SiteWhereMicroserviceGRV SiteWhere Micrservice Group Version Resource
func SiteWhereMicroserviceGRV() schema.GroupVersionResource {
	return sitewhereMicroserviceGVR
}

// SiteWhereTenantGRV SiteWhere Tenant Group Version Resource
func SiteWhereTenantGRV() schema.GroupVersionResource {
	return sitewhereTenantGVR
}
