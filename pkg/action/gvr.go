/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package action

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	sitewhereInstanceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha3",
		Resource: "instances",
	}
)

var (
	sitewhereMicroserviceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha3",
		Resource: "microservices",
	}
)
