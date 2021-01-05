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
	"github.com/sitewhere/swctl/pkg/status"
	"net/http"
	"os"

	"github.com/rakyll/statik/fs"

	"k8s.io/apimachinery/pkg/api/errors"

	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
	"github.com/sitewhere/swctl/pkg/uninstall"
)

// Uninstall is the action for installing SiteWhere
type Uninstall struct {
	cfg *Configuration

	// CRD indicates if we need to uninstall SiteWhere Custom Resource Definitions
	CRD bool
	// Infrastructure indicates if we need to install SiteWhere Infrastructure
	Infrastructure bool
	// Operator indicates if we need to install SiteWhere Operator
	Operator bool
	// Template indicates if we need to install SiteWhere templates
	Template bool
	StatikFS http.FileSystem

	// Minimal installation only install escential SiteWhere components.
	Minimal bool
	// Use verbose mode
	Verbose bool
	// Purge data
	Purge bool
}

// NewUninstall constructs a new *Uninstall
func NewUninstall(cfg *Configuration) *Uninstall {
	statikFS, _ := fs.New()
	return &Uninstall{
		cfg:            cfg,
		StatikFS:       statikFS,
		CRD:            true,
		Template:       true,
		Operator:       true,
		Infrastructure: true,
		Minimal:        false,
		Verbose:        false,
		Purge:          false,
	}
}

// Run executes the uninstall command, returning the result of the uninstallation
func (i *Uninstall) Run() (*uninstall.SiteWhereUninstall, error) {

	var err error
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}

	var infraStatuses []status.SiteWhereStatus
	if i.Infrastructure {
		// Uninstall Infrastructure
		infraStatuses, err = i.UninstallInfrastructure()
		if err != nil {
			fmt.Println(err)
		}
	}

	var operatorStatuses []status.SiteWhereStatus
	if i.Operator {
		// Uninstall Operator
		operatorStatuses, err = i.UninstallOperator()
		if err != nil {
			fmt.Println(err)
		}
	}

	var templatesStatues []status.SiteWhereStatus
	if i.Template {
		// Uninstall Templates
		templatesStatues, err = i.UninstallTemplates()
		if err != nil {
			fmt.Println(err)
		}
	}

	var crdDeletedStatuses []status.SiteWhereStatus

	if i.CRD {
		// Uninstall Custom Resource Definitions
		crdDeletedStatuses, err = i.UninstallCRDs()
		if err != nil {
			fmt.Println(err)
		}
	}

	return &uninstall.SiteWhereUninstall{CDRStatuses: crdDeletedStatuses, InfrastructureStatuses: infraStatuses, OperatorStatuses: operatorStatuses, TemplatesStatues: templatesStatues}, nil
}

// UninstallCRDs Uninstall SiteWhere Custom Resource Definitions
func (i *Uninstall) UninstallCRDs() ([]status.SiteWhereStatus, error) {
	return i.uninstallDirFiles(crdPath)
}

// UninstallTemplates Uninstall SiteWhere Templates CRD
func (i *Uninstall) UninstallTemplates() ([]status.SiteWhereStatus, error) {
	return i.uninstallDirFiles(templatePath)
}

// UninstallOperator Uninstall SiteWhere Operator resource file in the cluster
func (i *Uninstall) UninstallOperator() ([]status.SiteWhereStatus, error) {
	var result []status.SiteWhereStatus

	ns, err := i.uninstallDirFiles(namespacePath)
	if err != nil {
		return nil, err
	}
	result = append(result, ns...)

	certmager, err := i.uninstallDirFiles(certManagerPath)
	if err != nil {
		return nil, err
	}
	result = append(result, certmager...)

	operatorDeps, err := i.uninstallDirFiles(operatorDepsPath)
	if err != nil {
		return nil, err
	}
	result = append(result, operatorDeps...)

	operator, err := i.uninstallDirFiles(operatorPath)
	if err != nil {
		return nil, err
	}
	result = append(result, operator...)
	return result, nil
}

// UninstallInfrastructure Uninstall SiteWhere infrastructure
func (i *Uninstall) UninstallInfrastructure() ([]status.SiteWhereStatus, error) {
	var result []status.SiteWhereStatus

	infraDeps, err := i.uninstallDirFiles(infraDepsPath)
	if err != nil {
		return nil, err
	}
	result = append(result, infraDeps...)

	infra, err := i.uninstallDirFiles(infraPath)
	if err != nil {
		return nil, err
	}
	result = append(result, infra...)

	return result, nil
}

func (i *Uninstall) uninstallDirFiles(path string) ([]status.SiteWhereStatus, error) {
	r, err := i.StatikFS.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := r.Stat()
	if err != nil {
		return nil, err
	}
	return i.uninstallFiles("", fi)
}

func (i *Uninstall) uninstallFiles(parentPath string, fi os.FileInfo) ([]status.SiteWhereStatus, error) {

	var result []status.SiteWhereStatus

	if fi.IsDir() {
		dirName := parentPath + string(os.PathSeparator) + fi.Name()
		i.cfg.Log(fmt.Sprintf("Uninstalling Resources from %s", dirName))
		r, err := i.StatikFS.Open(dirName)
		if err != nil {
			return nil, err
		}
		files, err := r.Readdir(-1)
		if err != nil {
			return nil, err
		}
		for _, fileInfo := range files {
			unInstallResult, err := i.uninstallFiles(dirName, fileInfo)
			if err != nil && !errors.IsAlreadyExists(err) {
				return nil, err
			}
			result = append(result, unInstallResult...)
		}
	} else {
		var fileName = parentPath + string(os.PathSeparator) + fi.Name()
		i.cfg.Log(fmt.Sprintf("Uninstalling Resources %s", fileName))
		deployFile, err := i.StatikFS.Open(fileName)
		if err != nil {
			return nil, err
		}
		// Open the resource file
		res, err := i.cfg.KubeClient.Build(deployFile, false)
		if err != nil {
			return nil, err
		}
		if _, err := i.cfg.KubeClient.Delete(res); err != nil {
			var deleteStatus = status.SiteWhereStatus{
				Name:   fileName,
				Status: status.Unknown,
			}
			result = append(result, deleteStatus)
		} else {
			var deployStatus = status.SiteWhereStatus{
				Name:   fileName,
				Status: status.Uninstalled,
				//		ObjectMeta: createObject,
			}
			result = append(result, deployStatus)
		}
	}
	return result, nil
}
