/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package alpha3 defines SiteWhere Structures
package cmd

import "testing"

func TestInit(t *testing.T) {
	result := rootCmd.Commands()
	if result == nil {
		t.Error("Expected no nil list")
	}

	var found = false

	for _, cmd := range result {
		if "create" == cmd.Name() {
			found = true
			break
		}
	}

	if !found {
		t.Error("Command create not found")
	}
}
