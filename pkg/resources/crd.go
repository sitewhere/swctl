/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
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
func CreateCustomResourceFromFile(crName string, statikFS http.FileSystem, config *rest.Config) error {
	r, err := statikFS.Open(crName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crName, err)
		return err
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content of %s: %v\n", crName, err)
		return err
	}

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewDiscoveryClientForConfig for %s: %v\n", crName, err)
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewForConfig for %s: %v\n", crName, err)
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(contents), nil, obj)
	if err != nil {
		fmt.Printf("Error decoding for %s: %v\n", crName, err)
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Printf("Error finding GRV for %s: %v\n", crName, err)
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

	_, err = dr.Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("Error creating resource from file %s of Kind: %s: %v", crName, gvk.GroupKind().Kind, err)
		}
		return err
	}
	return nil
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
