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
