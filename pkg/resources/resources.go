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

package resources

import (
	"context"
	"fmt"

	"time"

	kubernetes "k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/api/errors"

	appsv1 "k8s.io/api/apps/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	deployRunningThreshold     = time.Minute * 10 // Max wait time
	deployRunningCheckInterval = time.Second * 2
)

func waitForPodContainersRunning(clientset kubernetes.Interface, podName string, namespace string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := podContainersRunning(clientset, podName, namespace)
		if running {
			return nil
		}

		if err != nil && errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for running pods: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get all running containers")
		}
	}
	return nil
}

func podContainersRunning(clientset kubernetes.Interface, podName string, namespace string) (bool, error) {
	existingPod, err := clientset.CoreV1().Pods(namespace).Get(
		context.TODO(),
		podName,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	for _, status := range existingPod.Status.ContainerStatuses {
		if !status.Ready {
			return false, nil
		}
	}

	return true, nil
}

// WaitForDeploymentAvailable waits for a Deployment to became available
func WaitForDeploymentAvailable(clientset kubernetes.Interface, deploymentName string, namespace string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := deploymentAvailable(clientset, deploymentName, namespace)
		if running {
			return nil
		}

		if err != nil && !errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get deployment available")
		}
	}
	return nil
}

func deploymentAvailable(clientset kubernetes.Interface, deploymentName string, namespace string) (bool, error) {
	existingDeploy, err := clientset.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	if existingDeploy.Status.ReadyReplicas < existingDeploy.Status.AvailableReplicas {
		return false, nil
	}
	for _, cond := range existingDeploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentProgressing {
			return true, nil
		}
	}
	return false, nil
}

// WaitForSecretExists waits for a Secret to exists
func WaitForSecretExists(clientset kubernetes.Interface, name string, namespace string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := secretExists(clientset, name, namespace)
		if running {
			return nil
		}

		if err != nil && !errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get Secret available")
		}
	}
	return nil
}

func secretExists(clientset kubernetes.Interface, name string, namespace string) (bool, error) {
	_, err := clientset.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	return true, nil
}

// WaitForCRDStablished waits for a CRD to be stablished
func WaitForCRDStablished(apiextensionsclientset apiextensionsclientset.Interface, name string) error {
	end := time.Now().Add(deployRunningThreshold)

	for true {
		<-time.NewTimer(deployRunningCheckInterval).C

		var err error
		running, err := crdStablished(apiextensionsclientset, name)
		if running {
			return nil
		}

		if err != nil && !errors.IsNotFound(err) {
			fmt.Printf(fmt.Sprintf("Encountered an error checking for deployment available: %s\n", err))
		}

		if time.Now().After(end) {
			return fmt.Errorf("Failed to get Secret available")
		}
	}
	return nil
}

func crdStablished(apiextensionsclientset apiextensionsclientset.Interface, name string) (bool, error) {
	crds := apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions()

	existingCRD, err := crds.Get(context.TODO(),
		name,
		metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	for _, cond := range existingCRD.Status.Conditions {
		if cond.Type == apiextv1beta1.Established {
			return true, nil
		}
	}
	return false, nil
}
