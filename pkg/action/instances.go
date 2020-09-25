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
	"github.com/sitewhere/swctl/pkg/instance"
	"github.com/sitewhere/swctl/pkg/resources"
)

// Instances is the action for listing SiteWhere instances
type Instances struct {
	cfg *Configuration
}

// NewInstances constructs a new *Instances
func NewInstances(cfg *Configuration) *Instances {
	return &Instances{
		cfg: cfg,
	}
}

// Run executes the install command, returning the result of the installation
func (i *Instances) Run() (*instance.ListSiteWhereInstance, error) {
	var err error
	dynamicClientset, err := i.cfg.KubernetesDynamicClientSet()
	if err != nil {
		return nil, err
	}
	clientset, err := i.cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	sitewhereInstances, err := resources.ListSitewhereInstacesCR(dynamicClientset, clientset)
	if err != nil {
		return nil, err
	}
	return &instance.ListSiteWhereInstance{
		Instances: sitewhereInstances,
	}, nil
}

// func handleListInstances() {
// 	// var err error

// 	// config, err := internal.GetKubeConfigFromKubeconfig()
// 	// if err != nil {
// 	// 	fmt.Printf("Error getting Kubernetes Config: %v\n", err)
// 	// 	return
// 	// }

// 	// client, err := dynamic.NewForConfig(config)
// 	// if err != nil {
// 	// 	fmt.Printf("Error getting Kubernetes Client: %v\n", err)
// 	// 	return
// 	// }

// 	// if err != nil {
// 	// 	fmt.Printf("Error reading SiteWhere Instances: %v\n", err)
// 	// 	return
// 	// }

// 	// template := "%-20s%-20s%-20s%-20s%-20s%-20s\n"
// 	// fmt.Printf(template, "NAME", "NAMESPACE", "CONFIG TMPL", "DATESET TMPL", "TM STATUS", "UM STATUS")

// 	// for _, instance := range sitewhereInstaces.Items {
// 	// 	sitewhereInstace := extractFromResource(&instance)
// 	// 	fmt.Printf(
// 	// 		template,
// 	// 		sitewhereInstace.Name,
// 	// 		sitewhereInstace.Namespace,
// 	// 		sitewhereInstace.ConfigurationTemplate,
// 	// 		sitewhereInstace.DatasetTemplate,
// 	// 		sitewhereInstace.Status.TenantManagementStatus,
// 	// 		sitewhereInstace.Status.UserManagementStatus,
// 	// 	)
// 	// }
// }

// func handleInstance(instanceName string) {
// 	// var err error

// 	// config, err := internal.GetKubeConfigFromKubeconfig()
// 	// if err != nil {
// 	// 	fmt.Printf("Error getting Kubernetes Config: %v\n", err)
// 	// 	return
// 	// }

// 	// client, err := dynamic.NewForConfig(config)
// 	// if err != nil {
// 	// 	fmt.Printf("Error getting Kubernetes Client: %v\n", err)
// 	// 	return
// 	// }

// 	// res := client.Resource(sitewhereInstanceGVR)
// 	// options := metav1.GetOptions{}
// 	// sitewhereInstace, err := res.Get(context.TODO(), instanceName, options)

// 	// if err != nil {
// 	// 	fmt.Printf(
// 	// 		"SiteWhere Instace %s NOT FOUND.\n",
// 	// 		instanceName)
// 	// 	return
// 	// }

// 	// printSiteWhereInstance(sitewhereInstace)
// }
