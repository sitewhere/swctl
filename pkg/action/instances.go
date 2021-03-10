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
	"strings"

	"github.com/pkg/errors"
	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctlcli "sigs.k8s.io/controller-runtime/pkg/client"

	"helm.sh/helm/v3/pkg/action"

	"github.com/sitewhere/swctl/pkg/instance"
)

// Instances is the action for listing SiteWhere instances
type Instances struct {
	cfg *action.Configuration

	// Name of the instance
	InstanceName string
}

// NewInstances constructs a new *Instances
func NewInstances(cfg *action.Configuration) *Instances {
	return &Instances{
		cfg:          cfg,
		InstanceName: "",
	}
}

// Run executes the install command, returning the result of the installation
func (i *Instances) Run() (*instance.ListSiteWhereInstance, error) {

	instanceList := instance.ListSiteWhereInstance{}
	if i.InstanceName != "" {
		i.instanceDetail(&instanceList)
	}
	i.instances(&instanceList)
	return &instanceList, nil
}

// ExtractInstanceNameArg returns the name of the instance that should be used.
func (i *Instances) ExtractInstanceNameArg(args []string) (string, error) {
	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	} else if len(args) == 1 {
		return args[0], nil
	}
	return "", nil
}

func (i *Instances) instances(instanceList *instance.ListSiteWhereInstance) error {
	var client, err = ControllerClient(i.cfg)
	if err != nil {
		return err
	}

	ctx := context.TODO()
	var swInstancesList sitewhereiov1alpha4.SiteWhereInstanceList

	if i.InstanceName == "" {
		if err := client.List(ctx, &swInstancesList); err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
		}
		instanceList.Instances = swInstancesList.Items
	} else {
		var swInstanceCR sitewhereiov1alpha4.SiteWhereInstance
		err = client.Get(ctx, ctlcli.ObjectKey{Name: i.InstanceName}, &swInstanceCR)
		if err != nil {
			return err
		}
		instanceList.Instances = append(instanceList.Instances, swInstanceCR)
	}
	return nil
}

func (i *Instances) instanceDetail(instanceList *instance.ListSiteWhereInstance) error {

	var client, err = ControllerClient(i.cfg)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return err
	}

	var swMicroservoceList sitewhereiov1alpha4.SiteWhereMicroserviceList
	err = client.List(ctx, &swMicroservoceList, ctlcli.InNamespace(i.InstanceName))
	instanceList.SiteWhereMicroservice = swMicroservoceList.Items
	return nil
}
