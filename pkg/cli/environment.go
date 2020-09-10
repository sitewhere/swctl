/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cli

import (
	"os"
	"strconv"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	namespace string

	config *genericclioptions.ConfigFlags

	// KubeConfig is the path to the kubeconfig file
	KubeConfig string
	// KubeContext is the name of the kubeconfig context.
	KubeContext string
	// Bearer KubeToken used for authentication
	KubeToken string
	// Kubernetes API Server Endpoint for authentication
	KubeAPIServer string

	// Debug indicates whether or not Helm is running in Debug mode.
	Debug bool
}

// New creates a new settings
func New() *EnvSettings {
	env := &EnvSettings{
		namespace:     os.Getenv("SWCTL_NAMESPACE"),
		KubeContext:   os.Getenv("SWCTL_KUBECONTEXT"),
		KubeToken:     os.Getenv("SWCTL_KUBETOKEN"),
		KubeAPIServer: os.Getenv("SWCTL_KUBEAPISERVER"),
	}

	env.Debug, _ = strconv.ParseBool(os.Getenv("SWCTL_DEBUG"))

	// bind to kubernetes config flags
	env.config = &genericclioptions.ConfigFlags{
		Namespace:   &env.namespace,
		Context:     &env.KubeContext,
		BearerToken: &env.KubeToken,
		APIServer:   &env.KubeAPIServer,
		KubeConfig:  &env.KubeConfig,
	}
	return env
}

// Namespace gets the namespace from the configuration
func (s *EnvSettings) Namespace() string {
	if ns, _, err := s.config.ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return "default"
}

// RESTClientGetter gets the kubeconfig from EnvSettings
func (s *EnvSettings) RESTClientGetter() genericclioptions.RESTClientGetter {
	return s.config
}
