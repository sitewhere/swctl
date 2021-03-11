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

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	"helm.sh/helm/v3/pkg/action"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

// ListMicroservices is the action for listing SiteWhere Microservices
type ListMicroservices struct {
	cfg *action.Configuration

	// Name of the instance
	InstanceName string
}

// ListMicroservicesResult is result the action for listing SiteWhere Microservices
type ListMicroservicesResult struct {
	// Name of the instance
	Name string `json:"name"`

	// Microservices are the microservices of a instance
	Microservices []sitewhereiov1alpha4.SiteWhereMicroservice `json:"microservices"`
}

// NewListMicroservices constructs a new *ListMicroservices
func NewListMicroservices(cfg *action.Configuration) *ListMicroservices {
	return &ListMicroservices{
		cfg:          cfg,
		InstanceName: "",
	}
}

// Run executes the install command, returning the result of the installation
func (i *ListMicroservices) Run() (*ListMicroservicesResult, error) {
	var err error
	// check for kubernetes cluster
	if err = i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}
	client, err := ControllerClient(i.cfg)
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()

	var swInstanceCR sitewhereiov1alpha4.SiteWhereInstance
	err = client.Get(ctx, ctlcli.ObjectKey{Name: i.InstanceName}, &swInstanceCR)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("sitewhere instance '%s' not found", i.InstanceName)
		}
		return nil, err
	}

	var swMicroservoceList sitewhereiov1alpha4.SiteWhereMicroserviceList
	err = client.List(ctx, &swMicroservoceList, ctlcli.InNamespace(i.InstanceName))
	if err != nil {
		return nil, err
	}

	return &ListMicroservicesResult{
		Name:          i.InstanceName,
		Microservices: swMicroservoceList.Items,
	}, nil
}
