/**
 * Copyright © 2014-2021 The SiteWhere Authors
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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"

	"github.com/sitewhere/swctl/pkg/install"
	"github.com/sitewhere/swctl/pkg/resources"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

const defaultHelmChartVersion string = "0.1.10"

// Install is the action for installing SiteWhere
type Install struct {
	cfg *action.Configuration

	settings *cli.EnvSettings

	// HelmChartVersion is the version of SiteWhere Infrastructure Helm Chart
	HelmChartVersion string

	// SkipCRD indicates if we need to install SiteWhere Custom Resource Definitions
	SkipCRD bool
	// SkipTemplate indicates if we need to install SiteWhere templates
	SkipTemplate bool
	// SkipOperator indicates if we need to install SiteWhere Operator
	SkipOperator bool
	// SkipInfrastructure indicates if we need to install SiteWhere Infrastructure
	SkipInfrastructure bool
	// Wait for components to be ready before return control.
	WaitReady bool
	// Use verbose mode
	Verbose bool
	// Minimal if true, deploy minimal infrastucure
	Minimal bool
	// StorageClass is the name of the storage class for the infrastructure
	StorageClass string
	// KafkaPVCStorageSize is the size of Kafka PVC Storage Size
	KafkaPVCStorageSize string
}

// NewInstall constructs a new *Install
func NewInstall(cfg *action.Configuration, settings *cli.EnvSettings) *Install {
	return &Install{
		cfg:                 cfg,
		settings:            settings,
		HelmChartVersion:    defaultHelmChartVersion,
		SkipCRD:             false,
		SkipTemplate:        false,
		SkipOperator:        false,
		SkipInfrastructure:  false,
		WaitReady:           false,
		Verbose:             false,
		Minimal:             false,
		StorageClass:        "",
		KafkaPVCStorageSize: "",
	}
}

// Run executes the install command, returning the result of the installation
func (i *Install) Run() (*install.SiteWhereInstall, error) {
	var err error
	err = i.CheckInstallPrerequisites()
	if err != nil {
		return nil, err
	}
	err = i.addSiteWhereRepository()
	if err != nil {
		return nil, err
	}
	err = i.updateSiteWhereRepository()
	if err != nil {
		return nil, err
	}
	return i.installRelease()
}

// CheckInstallPrerequisites checks for SiteWhere Install Prerequisites
func (i *Install) CheckInstallPrerequisites() error {
	var err error
	// check for kubernetes cluster
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return err
	}
	clientSet, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return err
	}
	// check for Istio installed on the cluster
	ok, err := resources.CheckIfExistsNamespace("istio-system", clientSet)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf(ErrIstioNotInstalled)
	}
	return nil
}

func (i *Install) addSiteWhereRepository() error {
	repoFile := i.settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(sitewhereRepoName) {
		fmt.Printf("repository name (%s) already exists\n", sitewhereRepoName)
		return nil
	}

	c := repo.Entry{
		Name: sitewhereRepoName,
		URL:  sitewhereRepoURL,
	}

	r, err := repo.NewChartRepository(&c, getter.All(i.settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", sitewhereRepoURL)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	return nil
}

func (i *Install) updateSiteWhereRepository() error {
	repoFile := i.settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(i.settings))
		if err != nil {
			return err
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()

	return nil
}

func (i *Install) installRelease() (*install.SiteWhereInstall, error) {
	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	var logConf action.DebugLog
	if i.Verbose {
		logConf = log.Printf
	} else {
		logConf = Discardf
	}

	if err := actionConfig.Init(i.settings.RESTClientGetter(), sitewhereSystemNamespace, os.Getenv("HELM_DRIVER"), logConf); err != nil {
		return nil, err
	}

	installAction := action.NewInstall(actionConfig)
	if installAction.Version == "" && installAction.Devel {
		installAction.Version = ">0.0.0-0"
	}
	installAction.Namespace = sitewhereSystemNamespace
	installAction.ReleaseName = sitewhereReleaseName
	installAction.CreateNamespace = true
	installAction.SkipCRDs = i.SkipCRD
	installAction.Wait = i.WaitReady
	installAction.Version = i.HelmChartVersion

	cp, err := installAction.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", sitewhereRepoName, sitewhereChartName), i.settings)

	p := getter.All(i.settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Skip operator
	vals["operator"] = map[string]interface{}{
		"enabled": !i.SkipOperator,
	}

	// Skip templates
	vals["templates"] = map[string]interface{}{
		"enabled": !i.SkipTemplate,
	}

	// Skip infrastructure
	vals["tags"] = map[string]interface{}{
		"infrastructure": !i.SkipInfrastructure,
	}
	vals["postgresql"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	vals["influxdb"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	vals["redis"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	vals["nifi"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	vals["mosquitto"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	vals["strimzi"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	vals["keycloak"] = map[string]interface{}{
		"enabled": !i.SkipInfrastructure,
	}
	if i.Minimal {
		vals["strimzi"] = map[string]interface{}{
			"replicas": 1,
			"isr":      1,
		}
	}

	// set storage class
	if i.StorageClass != "" {
		// InfluxDB storage class
		vals["influxdb"] = map[string]interface{}{
			"persistence": map[string]interface{}{
				"storageClass": i.StorageClass,
			},
		}
		// Redis
		vals["redis"] = map[string]interface{}{
			"global": map[string]interface{}{
				"storageClass": i.StorageClass,
			},
			"master": map[string]interface{}{
				"persistence": map[string]interface{}{
					"storageClass": i.StorageClass,
				},
			},
			"slave": map[string]interface{}{
				"persistence": map[string]interface{}{
					"storageClass": i.StorageClass,
				},
			},
		}
		// Nifi
		vals["nifi"] = map[string]interface{}{
			"persistence": map[string]interface{}{
				"storageClass": i.StorageClass,
			},
			"zookeeper": map[string]interface{}{
				"global": map[string]interface{}{
					"storageClass": i.StorageClass,
				},
				"persistence": map[string]interface{}{
					"storageClass": i.StorageClass,
				},
			},
		}
		// PostgreSQL
		vals["postgresql"] = map[string]interface{}{
			"global": map[string]interface{}{
				"storageClass": i.StorageClass,
			},
			"persistence": map[string]interface{}{
				"storageClass": i.StorageClass,
			},
		}
		// Keycloak
		vals["keycloak"] = map[string]interface{}{
			"postgresql": map[string]interface{}{
				"persistence": map[string]interface{}{
					"storageClass": i.StorageClass,
				},
			},
		}
		// Strimzi
		vals["strimzi"] = map[string]interface{}{
			"storage": map[string]interface{}{
				"type":         "persistent-claim",
				"storageClass": i.StorageClass,
			},
		}
	}

	// KafkaPVCStorageSize
	if i.KafkaPVCStorageSize != "" {
		vals["strimzi"] = map[string]interface{}{
			"storage": map[string]interface{}{
				"size": i.KafkaPVCStorageSize,
			},
		}
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if installAction.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          installAction.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: i.settings.RepositoryConfig,
					RepositoryCache:  i.settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	res, err := installAction.Run(chartRequested, vals)

	if err != nil {
		return nil, err
	}

	return &install.SiteWhereInstall{
		Release:   res.Name,
		Namespace: res.Namespace,
	}, nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
