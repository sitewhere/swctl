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

package instance

// CreateSiteWhereInstance destribe the creating of a SiteWhere Instance.
type CreateSiteWhereInstance struct {
	// Name of the instance
	InstanceName string `json:"instanceName"`
	// Namespace to use
	Namespace string `json:"namespace"`
	// Docker Image Tag
	Tag string `json:"tag"`
	// Number of replicas
	Replicas int64
	// Use debug mode
	Debug bool
	// Configuration Template
	ConfigurationTemplate string
	// Dataset template
	DatasetTemplate string
	// Service Account created
	ServiceAccountName string `json:"serviceAccountName"`
	// Custer Role created
	ClusterRoleName string `json:"clusterRoleName"`
	// Cluster Role Binding created
	ClusterRoleBindingName string `json:"clusterRoleBindingName"`
	// LoadBalancer Service created
	LoadBalanceServiceName string `json:"loadBalanceServiceName"`
	// Instance Custom Resources Name
	InstanceCustomResourceName string `json:"instanceCustomResourceName"`
}
