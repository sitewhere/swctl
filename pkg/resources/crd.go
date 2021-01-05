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
	"io/ioutil"
	"net/http"

	discovery "k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached/memory"
	dynamic "k8s.io/client-go/dynamic"

	rest "k8s.io/client-go/rest"
	restmapper "k8s.io/client-go/restmapper"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateCustomResourceDefinitionIfNotExists Create a CustomResourceDefinition if it does not exists.
func CreateCustomResourceDefinitionIfNotExists(crd *apiextv1beta1.CustomResourceDefinition, apiextensionsclientset apiextensionsclientset.Interface) (*apiextv1beta1.CustomResourceDefinition, error) {
	var err error

	crds := apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions()

	_, err = crds.Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	return crd, nil
}

// DeleteCustomResourceDefinitionIfExists Delete a CustomResourceDefinition if it exists
func DeleteCustomResourceDefinitionIfExists(crd *apiextv1beta1.CustomResourceDefinition, apiextensionsclientset apiextensionsclientset.Interface) error {
	return apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(
		context.TODO(),
		crd.ObjectMeta.Name,
		metav1.DeleteOptions{})
}

// CreateCustomResourceFromFile Reads a File from statik and creates a CustomResource from it.
func CreateCustomResourceFromFile(crName string, statikFS http.FileSystem, config *rest.Config) (*metav1.ObjectMeta, error) {
	r, err := statikFS.Open(crName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crName, err)
		return nil, err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content of %s: %v\n", crName, err)
		return nil, err
	}

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewDiscoveryClientForConfig for %s: %v\n", crName, err)
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewForConfig for %s: %v\n", crName, err)
		return nil, err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(contents), nil, obj)
	if err != nil {
		fmt.Printf("Error decoding for %s: %v\n", crName, err)
		return nil, err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Printf("Error finding GRV for %s: %v\n", crName, err)
		return nil, err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	created, err := dr.Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("Error creating resource from file %s of Kind: %s: %v", crName, gvk.GroupKind().Kind, err)
		}
		return nil, err
	}
	return &metav1.ObjectMeta{
		Name: created.GetName(),
	}, nil
}

// DeleteCustomResourceFromFile Reads a File from statik and deletes a CustomResource from it.
func DeleteCustomResourceFromFile(crName string, config *rest.Config, statikFS http.FileSystem) error {
	r, err := statikFS.Open(crName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crName, err)
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content: %v\n", err)
		return err
	}

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewDiscoveryClientForConfig: %v\n", err)
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewForConfig: %v\n", err)
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(contents), nil, obj)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	err = dr.Delete(context.TODO(), obj.GetName(), metav1.DeleteOptions{})

	if err != nil && !errors.IsNotFound(err) {
		fmt.Printf("Error deleting resource from file %s of Kind: %s: %v", crName, gvk.GroupKind().Kind, err)
		return err
	}
	return nil
}
