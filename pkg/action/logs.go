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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
	ctlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

// Logs is the action for showing SiteWhere Microservoce Logs
type Logs struct {
	cfg *action.Configuration

	// Name of the Instance
	InstanceName string

	// Name of the Microservice in the instance
	MicroserviceName string

	// Follow if true, the follow the logs
	Follow bool
}

// NewLogs constructs a new *Logs
func NewLogs(cfg *action.Configuration) *Logs {
	return &Logs{
		cfg:              cfg,
		InstanceName:     "",
		MicroserviceName: "",
		Follow:           false,
	}
}

// Run executes the logs command, returning the result of the uninstallation
func (i *Logs) Run() error {
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
			return fmt.Errorf("sitewhere instance '%s' not found", i.InstanceName)
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

	var namespace = swInstanceCR.GetName()
	var deployName = swMicroserviceCR.Status.Deployment
	var containerName = swMicroserviceCR.GetName()
	var podName = ""

	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return err
	}

	deploy, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	podList, err := getPodsForDeployment(ctx, clientset, deploy)
	if err != nil {
		return err
	}

	if len(podList.Items) <= 0 {
		return fmt.Errorf("no Pods of Microservice %s for Instance %s were found", i.MicroserviceName, i.InstanceName)
	}

	// For now, let's choose the 1st Pod
	podName = podList.Items[0].GetName()

	return getPodLogs(clientset, namespace, podName, containerName, i.Follow, os.Stdout)
}

func getPodsForDeployment(ctx context.Context, clientset kubernetes.Interface, deploy *appsv1.Deployment) (*v1.PodList, error) {
	var set = labels.Set(deploy.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	return clientset.CoreV1().Pods(deploy.GetNamespace()).List(ctx, listOptions)
}

func getPodLogs(clientset kubernetes.Interface, namespace string, podName string, containerName string, follow bool, out io.Writer) error {
	podLogOptions := v1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
	}

	request := clientset.CoreV1().
		Pods(namespace).
		GetLogs(podName, &podLogOptions)

	readCloser, err := request.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer readCloser.Close()

	r := bufio.NewReader(readCloser)
	for {
		bytes, err := r.ReadBytes('\n')
		if _, err := out.Write(bytes); err != nil {
			return err
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}
