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

package action

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rakyll/statik/fs"
	"k8s.io/apimachinery/pkg/api/errors"

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/sitewhere/swctl/pkg/install"
	"github.com/sitewhere/swctl/pkg/resources"
	"github.com/sitewhere/swctl/pkg/status"
)

// path for CRD manifests
const crdPath = "/crd/"

// path for template manifests
const templatePath = "/templates/"

// path for certificate manager
const certManagerPath = "/cm/"

// path to namespace objects
const namespacePath = "/namespace/"

// path for operator dependencies
const operatorDepsPath = "/operator-deps/"

// path for operator manifests
const operatorPath = "/operator/"

// path for infrastructure dependencies
const infraDepsPath = "/infra-deps/"

// path for operator infra
const infraPath = "/infra/"

// Install is the action for installing SiteWhere
type Install struct {
	cfg *Configuration

	StatikFS http.FileSystem
	// CRD indicates if we need to install SiteWhere Custom Resource Definitions
	CRD bool
	// Template indicates if we need to install SiteWhere templates
	Template bool
	// Operator indicates if we need to install SiteWhere Operator
	Operator bool
	// Infrastructure indicates if we need to install SiteWhere Infrastructure
	Infrastructure bool
	// Minimal installation only install escential SiteWhere components
	Minimal bool
	// Wait for components to be ready before return control.
	WaitReady bool
	// Use verbose mode
	Verbose bool
}

// NewInstall constructs a new *Install
func NewInstall(cfg *Configuration) *Install {
	statikFS, _ := fs.New()
	return &Install{
		cfg:            cfg,
		StatikFS:       statikFS,
		CRD:            true,
		Template:       true,
		Operator:       true,
		Infrastructure: true,
		Minimal:        false,
		WaitReady:      false,
		Verbose:        false,
	}
}

// Run executes the install command, returning the result of the installation
func (i *Install) Run() (*install.SiteWhereInstall, error) {
	var err error
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	var crdStatuses []status.SiteWhereStatus
	if i.CRD {
		// Install Custom Resource Definitions
		crdStatuses, err = i.InstallCRDs()
		if err != nil {
			return nil, err
		}
	}
	var templatesStatues []status.SiteWhereStatus
	if i.Template {
		// Install Templates
		templatesStatues, err = i.InstallTemplates()
		if err != nil {
			return nil, err
		}
	}
	var operatorStatuses []status.SiteWhereStatus
	if i.Operator {
		// Install Operator
		operatorStatuses, err = i.InstallOperator()
		if err != nil {
			return nil, err
		}
	}
	var infraStatuses []status.SiteWhereStatus
	if i.Infrastructure {
		// Install Infrastructure
		infraStatuses, err = i.InstallInfrastructure()
		if err != nil {
			return nil, err
		}
	}
	return &install.SiteWhereInstall{
		CDRStatuses:            crdStatuses,
		TemplatesStatues:       templatesStatues,
		OperatorStatuses:       operatorStatuses,
		InfrastructureStatuses: infraStatuses,
	}, nil
}

// InstallCRDs Install SiteWhere Custom Resource Definitions
func (i *Install) InstallCRDs() ([]status.SiteWhereStatus, error) {
	return i.installDirFiles(crdPath)
}

// InstallTemplates Install SiteWhere Templates CRD
func (i *Install) InstallTemplates() ([]status.SiteWhereStatus, error) {
	return i.installDirFiles(templatePath)
}

// InstallOperator Install SiteWhere Operator resource file in the cluster
func (i *Install) InstallOperator() ([]status.SiteWhereStatus, error) {
	var result []status.SiteWhereStatus
	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	ns, err := i.installDirFiles(namespacePath)
	if err != nil {
		return nil, err
	}
	result = append(result, ns...)

	certmager, err := i.installDirFiles(certManagerPath)
	if err != nil {
		return nil, err
	}
	result = append(result, certmager...)

	err = resources.WaitForDeploymentAvailable(clientset, "cert-manager-cainjector", "cert-manager")
	if err != nil {
		return nil, err
	}
	err = resources.WaitForDeploymentAvailable(clientset, "cert-manager", "cert-manager")
	if err != nil {
		return nil, err
	}
	err = resources.WaitForDeploymentAvailable(clientset, "cert-manager-webhook", "cert-manager")
	if err != nil {
		return nil, err
	}
	err = resources.WaitForSecretExists(clientset, "cert-manager-webhook-ca", "cert-manager")
	if err != nil {
		return nil, err
	}
	operatorDeps, err := i.installDirFiles(operatorDepsPath)
	if err != nil {
		return nil, err
	}
	result = append(result, operatorDeps...)

	err = resources.WaitForSecretExists(clientset, "webhook-server-cert", "sitewhere-system")
	if err != nil {
		return nil, err
	}

	operator, err := i.installDirFiles(operatorPath)
	if err != nil {
		return nil, err
	}
	result = append(result, operator...)
	return result, nil
}

// InstallInfrastructure Install SiteWhere infrastructure
func (i *Install) InstallInfrastructure() ([]status.SiteWhereStatus, error) {
	var result []status.SiteWhereStatus
	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	apiextensionsclientset, err := i.cfg.KubernetesAPIExtensionClientSet()
	if err != nil {
		return nil, err
	}

	infraDeps, err := i.installDirFiles(infraDepsPath)
	if err != nil {
		return nil, err
	}
	result = append(result, infraDeps...)

	err = resources.WaitForCRDStablished(apiextensionsclientset, "kafkas.kafka.strimzi.io")
	if err != nil {
		return nil, err
	}

	err = resources.WaitForDeploymentAvailable(clientset, "strimzi-cluster-operator", "sitewhere-system")
	if err != nil {
		return nil, err
	}

	infra, err := i.installDirFiles(infraPath)
	if err != nil {
		return nil, err
	}
	result = append(result, infra...)

	return result, nil
}

func (i *Install) installDirFiles(path string) ([]status.SiteWhereStatus, error) {
	r, err := i.StatikFS.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := r.Stat()
	if err != nil {
		return nil, err
	}
	return i.installFiles("", fi)
}

func (i *Install) installFiles(parentPath string, fi os.FileInfo) ([]status.SiteWhereStatus, error) {

	var result []status.SiteWhereStatus

	if fi.IsDir() {
		dirName := parentPath + string(os.PathSeparator) + fi.Name()
		i.cfg.Log(fmt.Sprintf("Installing Resources from %s", dirName))
		r, err := i.StatikFS.Open(dirName)
		if err != nil {
			return nil, err
		}
		files, err := r.Readdir(-1)
		if err != nil {
			return nil, err
		}
		for _, fileInfo := range files {
			installResult, err := i.installFiles(dirName, fileInfo)
			if err != nil && !errors.IsAlreadyExists(err) {
				return nil, err
			}
			result = append(result, installResult...)
		}
	} else {
		var fileName = parentPath + string(os.PathSeparator) + fi.Name()
		i.cfg.Log(fmt.Sprintf("Installing Resources %s", fileName))
		deployFile, err := i.StatikFS.Open(fileName)
		if err != nil {
			return nil, err
		}
		// Open the resource file
		res, err := i.cfg.KubeClient.Build(deployFile, false)
		if err != nil {
			return nil, err
		}
		if _, err := i.cfg.KubeClient.Create(res); err != nil {
			// If the error is Resource already exists, continue.
			if errors.IsAlreadyExists(err) {
				i.cfg.Log(fmt.Sprintf("Resource %s is already present. Skipping.", fileName))
			} else {
				var deployStatus = status.SiteWhereStatus{
					Name:   fileName,
					Status: status.Unknown,
				}
				result = append(result, deployStatus)
			}
		} else {
			var deployStatus = status.SiteWhereStatus{
				Name:   fileName,
				Status: status.Installed,
				//		ObjectMeta: createObject,
			}
			result = append(result, deployStatus)
		}
	}
	return result, nil
}
