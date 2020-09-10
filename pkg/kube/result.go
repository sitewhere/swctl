/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package kube

// Result contains the information of created, updated, and deleted resources
// for various kube API calls along with helper methods for using those
// resources
type Result struct {
	Created ResourceList
	Updated ResourceList
	Deleted ResourceList
}

// If needed, we can add methods to the Result type for things like diffing
