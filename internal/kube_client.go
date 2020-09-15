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

// Package internal Implements swctl internal use only functions
package internal

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//var clientset *kubernetes.Clientset
//var apixClient *apixv1beta1client.ApiextensionsV1beta1Client

// GetKubeConfigFromKubeconfig Buid a Kubernetes Config from ~/.kube/config
func GetKubeConfigFromKubeconfig() (*rest.Config, error) {
	kubeconfig := filepath.Join(userHomeDir(), ".kube", "config")
	return GetKubeConfig(kubeconfig)
}

// GetKubeConfig Buid a Kubernetes Config from a filepath
func GetKubeConfig(pathToCfg string) (*rest.Config, error) {
	if pathToCfg == "" {
		// in cluster access
		return rest.InClusterConfig()
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", pathToCfg)

	if err != nil {
		log.Println("Using in cluster config")
		return GetKubeConfig("")
	}
	return cfg, nil
}

// GetKubernetesClient Returns Kubernetes Client
func GetKubernetesClient(pathToCfg string) (*kubernetes.Clientset, error) {
	config, err := GetKubeConfig(pathToCfg)

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// GetKubernetesClientV1Beta1 Returns Kubernetes Client v1 beta1 interface
func GetKubernetesClientV1Beta1(pathToCfg string) (*apixv1beta1client.ApiextensionsV1beta1Client, error) {
	config, err := GetKubeConfig(pathToCfg)

	if err != nil {
		return nil, err
	}

	return apixv1beta1client.NewForConfig(config)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}
