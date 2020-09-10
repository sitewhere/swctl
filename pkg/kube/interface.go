/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package kube

import (
	"io"
	"time"

	v1 "k8s.io/api/core/v1"
)

// Interface represents a client capable of communicating with the Kubernetes API.
//
// A KubernetesClient must be concurrency safe.
type Interface interface {
	// Create creates one or more resources.
	Create(resources ResourceList) (*Result, error)

	//Wait(resources ResourceList, timeout time.Duration) error

	// Delete destroys one or more resources.
	Delete(resources ResourceList) (*Result, []error)

	// Watch the resource in reader until it is "ready". This method
	//
	// For Jobs, "ready" means the Job ran to completion (exited without error).
	// For Pods, "ready" means the Pod phase is marked "succeeded".
	// For all other kinds, it means the kind was created or modified without
	// error.
	WatchUntilReady(resources ResourceList, timeout time.Duration) error

	// Update updates one or more resources or creates the resource
	// if it doesn't exist.
	Update(original, target ResourceList, force bool) (*Result, error)

	// Build creates a resource list from a Reader
	//
	// reader must contain a YAML stream (one or more YAML documents separated
	// by "\n---\n")
	//
	// Validates against OpenAPI schema if validate is true.
	Build(reader io.Reader, validate bool) (ResourceList, error)

	// WaitAndGetCompletedPodPhase waits up to a timeout until a pod enters a completed phase
	// and returns said phase (PodSucceeded or PodFailed qualify).
	WaitAndGetCompletedPodPhase(name string, timeout time.Duration) (v1.PodPhase, error)

	// isReachable checks whether the client is able to connect to the cluster
	IsReachable() error
}

var _ Interface = (*Client)(nil)
