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

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	"helm.sh/helm/v3/pkg/action"
	ctlcli "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sitewhere/swctl/pkg/logs"
)

// LogLevel is the action for changing SiteWhere Microservoce Log Levels
type LogLevel struct {
	cfg *action.Configuration

	// Name of the Instance
	InstanceName string

	// Name of the Microservice in the instance
	MicroserviceName string

	// Level is the new log level to set
	Level logs.Level

	// Logger are the logger to change the level to the new value
	Logger []string
}

// NewLogLevel constructs a new *Logs
func NewLogLevel(cfg *action.Configuration) *LogLevel {
	return &LogLevel{
		cfg:              cfg,
		InstanceName:     "",
		MicroserviceName: "",
		Level:            "",
	}
}

// Run executes the log-level command
func (i *LogLevel) Run() error {
	var err error
	// check for kubernetes cluster
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return err
	}

	controllerClient, err := ControllerClient(i.cfg)
	if err != nil {
		return err
	}

	var ctx = context.TODO()

	// Find the SiteWhere Instance
	var swInstanceCR sitewhereiov1alpha4.SiteWhereInstance
	err = controllerClient.Get(ctx, ctlcli.ObjectKey{Name: i.InstanceName}, &swInstanceCR)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("the Instance %s does not exists", i.InstanceName)
		}
		return err
	}

	// Find the SiteWhere Microservice
	var swMicroserviceCR sitewhereiov1alpha4.SiteWhereMicroservice
	var objectKey ctlcli.ObjectKey = ctlcli.ObjectKey{
		Namespace: i.InstanceName,
		Name:      i.MicroserviceName,
	}

	err = controllerClient.Get(ctx, objectKey, &swMicroserviceCR)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("the Microservice %s for Instance %s does not exists", i.MicroserviceName, i.InstanceName)
		}
		return err
	}

	newOverrides, err := i.generateLoggingOverrides(swMicroserviceCR.Spec.Logging.Overrides)
	if err != nil {
		return err
	}
	swMicroserviceCR.Spec.Logging.Overrides = newOverrides

	err = controllerClient.Update(ctx, &swMicroserviceCR, &ctlcli.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GenerateLoggingOverrides Transforms the logging override
func (i *LogLevel) generateLoggingOverrides(entries []sitewhereiov1alpha4.MicroserviceLoggingEntry) ([]sitewhereiov1alpha4.MicroserviceLoggingEntry, error) {
	var result []sitewhereiov1alpha4.MicroserviceLoggingEntry
	for _, entry := range entries {
		transformedEntry := transformEntry(entry, i.Level, i.Logger)
		result = append(result, transformedEntry)
	}
	return result, nil
}

func transformEntry(entry sitewhereiov1alpha4.MicroserviceLoggingEntry, level logs.Level, logger []string) sitewhereiov1alpha4.MicroserviceLoggingEntry {
	var shouldTransform = false
	if len(logger) > 0 {
		for _, l := range logger {
			if l == entry.Logger {
				shouldTransform = true
				break
			}
		}
	} else {
		shouldTransform = true
	}
	if !shouldTransform {
		return entry
	}
	return sitewhereiov1alpha4.MicroserviceLoggingEntry{
		Logger: entry.Logger,
		Level:  string(level),
	}
}
