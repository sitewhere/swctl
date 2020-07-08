/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package alpha3 defines SiteWhere Structures
package alpha3

import "testing"

func TestGetSiteWhereMicroservicesList(t *testing.T) {

	result := GetSiteWhereMicroservicesList()

	if result == nil {
		t.Error("Expected no nil list")
	}
	if len(result) <= 0 {
		t.Error("Expected no empty list")
	}
}
