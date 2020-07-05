/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"

	//	"k8s.io/client-go/dynamic"
	discovery "k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/rest"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"

	v1 "k8s.io/api/core/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/sitewhere/swctl/internal"
	_ "github.com/sitewhere/swctl/internal/statik" // User for statik
)

var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install SiteWhere CRD and Operator",
	Long: `Use this command to install SiteWhere 3.0 on a Kubernetes Cluster.
This command will install:
 - SiteWhere System Namespace: sitewhere-system (default)
 - SiteWhere Custom Resources Definition.
 - SiteWhere Templates.
 - SiteWhere Operator.
 - SiteWhere Infrastructure.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		config, err := internal.GetKubeConfigFromKubeconfig()
		if err != nil {
			fmt.Printf("Error getting Kubernetes Config: %v\n", err)
			return
		}

		statikFS, err := fs.New()
		if err != nil {
			fmt.Printf("Error Reading Resources: %v\n", err)
			return
		}

		// Install Custom Resource Definitions
		installSiteWhereCRDs(config, statikFS)
		// Install Templates
		installSiteWhereTemplates(config, statikFS)
		// Install Operator
		installSiteWhereOperator(config, statikFS)
		// Install Infrastructure
		installSiteWhereInfrastructure(config, statikFS)

		fmt.Printf("SiteWhere 3.0 Installed\n")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installSiteWhereCRDs(config *rest.Config, statikFS http.FileSystem) {

	for i := 1; i <= 14; i++ {
		var crdName = fmt.Sprintf("/crd/crd-%02d.yaml", i)
		installSiteWhereCRD(crdName, config, statikFS)
	}
}

func installSiteWhereTemplates(config *rest.Config, statikFS http.FileSystem) {

	for i := 1; i <= 39; i++ {
		var templateName = fmt.Sprintf("/templates/template-%02d.yaml", i)
		installSiteWhereTemplate(templateName, config, statikFS)
	}
}

func installSiteWhereOperator(config *rest.Config, statikFS http.FileSystem) {
	var err error

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return
	}

	_, err = internal.CreateNamespaceIfNotExists("sitewhere-system", clientset)

	for i := 1; i <= 23; i++ {
		var operatorResource = fmt.Sprintf("/operator/operator-%02d.yaml", i)
		installSiteWhereOperatorResource(operatorResource, config, statikFS)
	}
}

func installSiteWhereInfrastructure(config *rest.Config, statikFS http.FileSystem) {

	for i := 1; i <= 28; i++ {
		var infraResource = fmt.Sprintf("/infra-min/infra-min-%02d.yaml", i)
		installSiteWhereOperatorResource(infraResource, config, statikFS)
	}
}

func installSiteWhereCRD(crdName string, config *rest.Config, statikFS http.FileSystem) {

	r, err := statikFS.Open(crdName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crdName, err)
		return
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content of file%s: %v\n", crdName, err)
		return
	}

	sch := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(sch)
	_ = apiextv1beta1.AddToScheme(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, groupVersionKind, err := decode([]byte(contents), nil, nil)

	_ = groupVersionKind

	if err != nil {
		// If we can decode, try installing custom resource
		installSiteWhereTemplate(crdName, config, statikFS)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return
	}

	// now use switch over the type of the object
	// and match each type-case
	switch o := obj.(type) {
	case *v1.Pod:
		_, err = internal.CreatePodIfNotExists(o, clientset, "sitewhere-system")
	case *v1.ConfigMap:
		_, err = internal.CreateConfigMapIfNotExists(o, clientset, "sitewhere-system")
	case *v1.Secret:
		_, err = internal.CreateSecretIfNotExists(o, clientset, "sitewhere-system")
	case *v1.ServiceAccount:
		_, err = internal.CreateServiceAccountIfNotExists(o, clientset, "sitewhere-system")
	case *v1.PersistentVolumeClaim:
		_, err = internal.CreatePersistentVolumeClaimIfNotExists(o, clientset, "sitewhere-system")
	case *v1.Service:
		_, err = internal.CreateServiceIfNotExists(o, clientset, "sitewhere-system")
	case *appsv1.Deployment:
		_, err = internal.CreateDeploymentIfNotExists(o, clientset, "sitewhere-system")
	case *appsv1.StatefulSet:
		_, err = internal.CreateStatefulSetIfNotExists(o, clientset, "sitewhere-system")
	case *rbacV1.ClusterRole:
		_, err = internal.CreateClusterRoleIfNotExists(o, clientset)
	case *rbacV1.ClusterRoleBinding:
		_, err = internal.CreateClusterRoleBindingIfNotExists(o, clientset)
	case *rbacV1.Role:
		_, err = internal.CreateRoleIfNotExists(o, clientset, "sitewhere-system")
	case *rbacV1.RoleBinding:
		_, err = internal.CreateRoleBindingIfNotExists(o, clientset, "sitewhere-system")
	case *policyV1beta1.PodDisruptionBudget:
		_, err = internal.CreatePodDisruptionBudgetIfNotExists(o, clientset, "sitewhere-system")
	case *v1beta1.CustomResourceDefinition:
		createCustomResourceDefinition(o, config)

	default:
		fmt.Println(fmt.Sprintf("Resource with type %v not handled.", groupVersionKind))
		_ = o //o is unknown for us
	}

	if err != nil {
		fmt.Printf("Error Creating Resource: %v\n", err)
		return
	}

}

func installSiteWhereTemplate(crdName string, config *rest.Config, statikFS http.FileSystem) {
	r, err := statikFS.Open(crdName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", crdName, err)
		return
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading content: %v\n", err)
		return
	}

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewDiscoveryClientForConfig: %v\n", err)
		return
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting NewForConfig: %v\n", err)
		return
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(contents), nil, obj)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Printf("Error finding GRV: %v\n", err)
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
			fmt.Printf("Error creating resource from file %s of Kind: %s: %v", crdName, gvk.GroupKind().Kind, err)
		}
		return
	}
}

func installSiteWhereOperatorResource(opResourceName string, config *rest.Config, statikFS http.FileSystem) {
	installSiteWhereCRD(opResourceName, config, statikFS)
}

func createCustomResourceDefinition(crd *v1beta1.CustomResourceDefinition, config *rest.Config) {
	var err error

	apixClient, err := apixv1beta1client.NewForConfig(config)

	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return
	}

	crds := apixClient.CustomResourceDefinitions()

	_, err = crds.Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			fmt.Printf("Failed to create CRD %s: %v\n", crd.ObjectMeta.Name, err)
		}
	}
}
