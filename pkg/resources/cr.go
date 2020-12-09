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

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/resources/grv"

	sitewhereiov1alpha4 "github.com/sitewhere/sitewhere-k8s-operator/apis/sitewhere.io/v1alpha4"
)

// CreateSiteWhereInstanceCR Creates a Custom Resource for a SiteWhere Instance
func CreateSiteWhereInstanceCR(cr *unstructured.Unstructured, dynamicClient dynamic.Interface) (*unstructured.Unstructured, error) {
	sitewhereInstanceGVR := grv.SiteWhereInstanceGRV()
	res := dynamicClient.Resource(sitewhereInstanceGVR)
	return res.Create(context.TODO(), cr, metav1.CreateOptions{})
}

// ListSitewhereInstacesCR List SiteWhere Instance CR installed
func ListSitewhereInstacesCR(dynamicClient dynamic.Interface, clientset kubernetes.Interface) ([]*sitewhereiov1alpha4.SiteWhereInstance, error) {
	sitewhereInstanceGVR := grv.SiteWhereInstanceGRV()
	res := dynamicClient.Resource(sitewhereInstanceGVR)
	sitewhereInstaces, err := res.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var result []*sitewhereiov1alpha4.SiteWhereInstance
	for _, instance := range sitewhereInstaces.Items {
		sitewhereInstace, err := extractFromResource(&instance, clientset)
		if err != nil {
			return nil, err
		}
		result = append(result, sitewhereInstace)
	}
	return result, nil
}

// CreateSiteWhereMicroserviceCR Creates a Custom Resource for a SiteWhere Microservice
func CreateSiteWhereMicroserviceCR(cr *unstructured.Unstructured, namespace string, dynamicClient dynamic.Interface) (*unstructured.Unstructured, error) {
	sitewhereMicroserviceGVR := grv.SiteWhereMicroserviceGRV()
	res := dynamicClient.Resource(sitewhereMicroserviceGVR).Namespace(namespace)
	return res.Create(context.TODO(), cr, metav1.CreateOptions{})
}

// CreateSiteWhereTenantCR Creates a Custom Resource for a SiteWhere Tenant
func CreateSiteWhereTenantCR(cr *unstructured.Unstructured, namespace string, dynamicClient dynamic.Interface) (*unstructured.Unstructured, error) {
	sitewhereTenantGVR := grv.SiteWhereTenantGRV()
	res := dynamicClient.Resource(sitewhereTenantGVR).Namespace(namespace)
	return res.Create(context.TODO(), cr, metav1.CreateOptions{})
}

func extractFromResource(crSiteWhereInstace *unstructured.Unstructured, clientset kubernetes.Interface) (*sitewhereiov1alpha4.SiteWhereInstance, error) {
	var result = sitewhereiov1alpha4.SiteWhereInstance{}
	/*
		metadata, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "metadata")
		if err != nil {
			return nil, fmt.Errorf("Error reading metadata for %s: %v", crSiteWhereInstace, err)
		}
		if !exists {
			fmt.Printf("Metadata not found for for SiteWhere Instance: %s", crSiteWhereInstace)
		} else {
			extractSiteWhereInstanceMetadata(metadata, &result)
		}
		spec, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "spec")
		if err != nil {
			return nil, fmt.Errorf("Error reading spec for %s: %v", result.Name, err)
		}
		if !exists {
			fmt.Printf("Spec not found for for SiteWhere Instance: %s", result.Name)
		} else {
			extractSiteWhereInstanceSpec(spec, &result)
		}
		status, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "status")
		if err != nil {
			return nil, fmt.Errorf("Error reading status for %s: %v", result.Name, err)
		}
		if !exists {
			result.Status = &alpha3.SiteWhereInstanceStatus{
				TenantManagementStatus: "Unknown",
				UserManagementStatus:   "Unknown",
			}
		} else {
			extractSiteWhereInstanceStatus(status, &result)
		}
		microservices, err := queryMicroservices(result.Name, clientset)
		if err != nil {
			return nil, fmt.Errorf("Error reading microservices statuses for %s: %v", result.Name, err)
		}
		result.Microservices = microservices
	*/
	return &result, nil
}

func extractSiteWhereInstanceMetadata(metadata map[string]interface{}, instance *alpha3.SiteWhereInstance) {
	name, exists, err := unstructured.NestedString(metadata, "name")
	if err != nil {
		fmt.Printf("Error Name from Metadata: %v\n", err)
	} else if !exists {
		fmt.Printf("Name from Metadata")
	} else {
		instance.Name = name
	}
}

func extractSiteWhereInstanceSpec(spec map[string]interface{}, instance *alpha3.SiteWhereInstance) {
	configurationTemplate := spec["configurationTemplate"]
	datasetTemplate := spec["datasetTemplate"]
	configuration := spec["configuration"]
	sitewhereConfiguration := extractSiteWhereInstanceConfiguration(configuration)
	instance.ConfigurationTemplate = fmt.Sprintf("%v", configurationTemplate)
	instance.DatasetTemplate = fmt.Sprintf("%v", datasetTemplate)
	instance.Configuration = sitewhereConfiguration
}

func extractSiteWhereInstanceStatus(status map[string]interface{}, instance *alpha3.SiteWhereInstance) {
	tenantManagementBootstrapState, exists, err := unstructured.NestedString(status, "tenantManagementBootstrapState")
	if err != nil || !exists {
		tenantManagementBootstrapState = "Unknown"
	}
	userManagementBootstrapState, exists, err := unstructured.NestedString(status, "userManagementBootstrapState")
	if err != nil || !exists {
		userManagementBootstrapState = "Unknown"
	}
	instance.Status = &alpha3.SiteWhereInstanceStatus{
		TenantManagementStatus: tenantManagementBootstrapState,
		UserManagementStatus:   userManagementBootstrapState,
	}
}

func extractSiteWhereInstanceConfiguration(config interface{}) *alpha3.SiteWhereInstanceConfiguration {
	var result = alpha3.SiteWhereInstanceConfiguration{}
	if configMap, ok := config.(map[string]interface{}); ok {
		infrastructure := configMap["infrastructure"]
		if infrastructure != nil {
			result.Infrastructure = extractSiteWhereInstanceConfigurationInfrastructure(infrastructure)
		}
		persistenceConfigurations := configMap["persistenceConfigurations"]
		if persistenceConfigurations != nil {
			result.Persistence = extractSiteWhereInstanceConfigurationPersistenceConfiguration(persistenceConfigurations)
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructure(infrastructureConfig interface{}) *alpha3.SiteWhereInstanceInfrastructureConfiguration {
	var result = alpha3.SiteWhereInstanceInfrastructureConfiguration{}
	if configMap, ok := infrastructureConfig.(map[string]interface{}); ok {
		namespace, exists, err := unstructured.NestedString(configMap, "namespace")
		if err != nil {
			fmt.Printf("Error reading Infrastructure Namespace: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Infrastructure Namespace not found")
		} else {
			result.Namespace = namespace
		}
		gRPC := configMap["grpc"]
		if gRPC != nil {
			result.GRPC = extractSiteWhereInstanceConfigurationInfrastructureGRPC(gRPC)
		}
		kafka := configMap["kafka"]
		if kafka != nil {
			result.Kafka = extractSiteWhereInstanceConfigurationInfrastructureKafka(kafka)
		}
		metrics := configMap["metrics"]
		if kafka != nil {
			result.Metrics = extractSiteWhereInstanceConfigurationInfrastructureMetrics(metrics)
		}
		redis := configMap["redis"]
		if kafka != nil {
			result.Redis = extractSiteWhereInstanceConfigurationInfrastructureRedis(redis)
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureGRPC(gRPCConfig interface{}) *alpha3.SiteWhereInstanceInfrastructureGRPCConfiguration {
	var result = alpha3.SiteWhereInstanceInfrastructureGRPCConfiguration{}
	if configMap, ok := gRPCConfig.(map[string]interface{}); ok {
		backoffMultiplier, exists, err := unstructured.NestedFloat64(configMap, "backoffMultiplier")
		if err != nil {
			fmt.Printf("Error reading backoffMultiplier: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("backoffMultiplier not found")
		} else {
			result.BackoffMultiplier = backoffMultiplier
		}
		initialBackoffSeconds, exists, err := unstructured.NestedInt64(configMap, "initialBackoffSeconds")
		if err != nil {
			fmt.Printf("Error reading initialBackoffSeconds: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("initialBackoffSeconds not found")
		} else {
			result.InitialBackoffSeconds = initialBackoffSeconds
		}
		maxBackoffSeconds, exists, err := unstructured.NestedInt64(configMap, "maxBackoffSeconds")
		if err != nil {
			fmt.Printf("Error reading maxBackoffSeconds: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("maxBackoffSeconds not found")
		} else {
			result.MaxBackoffSeconds = maxBackoffSeconds
		}
		maxRetryCount, exists, err := unstructured.NestedInt64(configMap, "maxRetryCount")
		if err != nil {
			fmt.Printf("Error reading maxRetryCount: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("maxRetryCount not found")
		} else {
			result.MaxRetryCount = maxRetryCount
		}
		resolveFQDN, exists, err := unstructured.NestedBool(configMap, "resolveFQDN")
		if err != nil {
			fmt.Printf("Error reading resolveFQDN: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("resolveFQDN not found")
		} else {
			result.ResolveFQDN = resolveFQDN
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureKafka(kafkaConfig interface{}) *alpha3.SiteWhereInstanceInfrastructureKafkaConfiguration {
	var result = alpha3.SiteWhereInstanceInfrastructureKafkaConfiguration{}
	if configMap, ok := kafkaConfig.(map[string]interface{}); ok {
		port, exists, err := unstructured.NestedInt64(configMap, "port")
		if err != nil {
			fmt.Printf("Error reading Kafka Port: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Kafka Port not found")
		} else {
			result.Port = port
		}
		hostname, exists, err := unstructured.NestedString(configMap, "hostname")
		if err != nil {
			fmt.Printf("Error reading Kafka Hostname: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Kafka Hostname not found")
		} else {
			result.Hostname = hostname
		}
		defaultTopicPartitions, exists, err := unstructured.NestedInt64(configMap, "defaultTopicPartitions")
		if err != nil {
			fmt.Printf("Error reading Kafka defaultTopicPartitions: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Kafka defaultTopicPartitions not found")
		} else {
			result.DefaultTopicPartitions = defaultTopicPartitions
		}
		defaultTopicReplicationFactor, exists, err := unstructured.NestedInt64(configMap, "defaultTopicReplicationFactor")
		if err != nil {
			fmt.Printf("Error reading Kafka defaultTopicReplicationFactor: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Kafka defaultTopicReplicationFactor not found")
		} else {
			result.DefaultTopicReplicationFactor = defaultTopicReplicationFactor
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureMetrics(metricsConfig interface{}) *alpha3.SiteWhereInstanceInfrastructureMetricsConfiguration {
	var result = alpha3.SiteWhereInstanceInfrastructureMetricsConfiguration{}
	if configMap, ok := metricsConfig.(map[string]interface{}); ok {
		enabled, exists, err := unstructured.NestedBool(configMap, "enabled")
		if err != nil {
			fmt.Printf("Error reading Metrics Enabled: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Metrics Enabled not found")
		} else {
			result.Enabled = enabled
		}
		httpPort, exists, err := unstructured.NestedInt64(configMap, "httpPort")
		if err != nil {
			fmt.Printf("Error reading Metrics HTTP Port: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Metrics HTTP Port not found")
		} else {
			result.HTTPPort = httpPort
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureRedis(redisConfig interface{}) *alpha3.SiteWhereInstanceInfrastructureRedisConfiguration {
	var result = alpha3.SiteWhereInstanceInfrastructureRedisConfiguration{}
	if configMap, ok := redisConfig.(map[string]interface{}); ok {
		hostname, exists, err := unstructured.NestedString(configMap, "hostname")
		if err != nil {
			fmt.Printf("Error reading Redis Hostname: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Redis Hostname not found")
		} else {
			result.Hostname = hostname
		}
		port, exists, err := unstructured.NestedInt64(configMap, "port")
		if err != nil {
			fmt.Printf("Error reading Redis Port: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Redis Port not found")
		} else {
			result.Port = port
		}
		nodeCount, exists, err := unstructured.NestedInt64(configMap, "nodeCount")
		if err != nil {
			fmt.Printf("Error reading Redis Node Count: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Redis Node Count not found")
		} else {
			result.NodeCount = nodeCount
		}
		masterGroupName, exists, err := unstructured.NestedString(configMap, "masterGroupName")
		if err != nil {
			fmt.Printf("Error reading Redis Master Group Name: %v\n", err)
			return nil
		}
		if !exists {
			fmt.Printf("Redis Master Group Name not found")
		} else {
			result.MasterGroupName = masterGroupName
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationPersistenceConfiguration(persistenceConfig interface{}) *alpha3.SiteWhereInstancePersistenceConfiguration {
	var result = alpha3.SiteWhereInstancePersistenceConfiguration{}
	if configMap, ok := persistenceConfig.(map[string]interface{}); ok {
		cassandraConfigurations := configMap["cassandraConfigurations"]
		if cassandraConfigurations != nil {
			result.CassandraConfigurations = extractSiteWhereInstanceConfigurationPersistenceCassandraConfigurations(cassandraConfigurations)
		}
		influxDbConfigurations := configMap["influxDbConfigurations"]
		if influxDbConfigurations != nil {
			result.InfluxDBConfigurations = extractSiteWhereInstanceConfigurationPersistenceInfluxDBConfigurations(influxDbConfigurations)
		}
		rdbConfigurations := configMap["rdbConfigurations"]
		if rdbConfigurations != nil {
			result.RDBConfigurations = extractSiteWhereInstanceConfigurationPersistenceRDBConfigurations(rdbConfigurations)
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationPersistenceCassandraConfigurations(cassandraConfig interface{}) map[string]alpha3.SiteWhereInstancePersistenceCassandraConfiguration {
	if configMap, ok := cassandraConfig.(map[string]interface{}); ok {
		result := make(map[string]alpha3.SiteWhereInstancePersistenceCassandraConfiguration)
		for key, value := range configMap {
			var configuration = extractSiteWhereInstanceConfigurationPersistenceCassandraConfiguration(value)
			result[key] = configuration
		}
		return result
	}
	return nil
}

func extractSiteWhereInstanceConfigurationPersistenceCassandraConfiguration(cassandraConfig interface{}) alpha3.SiteWhereInstancePersistenceCassandraConfiguration {
	var result = alpha3.SiteWhereInstancePersistenceCassandraConfiguration{}
	if configMap, ok := cassandraConfig.(map[string]interface{}); ok {
		contactPoints, exists, err := unstructured.NestedString(configMap, "contactPoints")
		if err != nil {
			fmt.Printf("Error reading Cassandra Contact Points: %v\n", err)
		} else if !exists {
			fmt.Printf("Cassandra Contact Points not found")
		} else {
			result.ContactPoints = contactPoints
		}
		keyspace, exists, err := unstructured.NestedString(configMap, "keyspace")
		if err != nil {
			fmt.Printf("Error reading Cassandra Keyspace: %v\n", err)
		} else if !exists {
			fmt.Printf("Cassandra Keyspace not found")
		} else {
			result.Keyspace = keyspace
		}
	}
	return result
}

func extractSiteWhereInstanceConfigurationPersistenceInfluxDBConfigurations(influxDBConfig interface{}) map[string]alpha3.SiteWhereInstancePersistenceInfluxDBConfiguration {
	if configMap, ok := influxDBConfig.(map[string]interface{}); ok {
		result := make(map[string]alpha3.SiteWhereInstancePersistenceInfluxDBConfiguration)
		for key, value := range configMap {
			var configuration = extractSiteWhereInstanceConfigurationPersistenceInfluxDBConfiguration(value)
			result[key] = configuration
		}
		return result
	}
	return nil
}

func extractSiteWhereInstanceConfigurationPersistenceInfluxDBConfiguration(influxDBConfig interface{}) alpha3.SiteWhereInstancePersistenceInfluxDBConfiguration {
	var result = alpha3.SiteWhereInstancePersistenceInfluxDBConfiguration{}
	if configMap, ok := influxDBConfig.(map[string]interface{}); ok {
		port, exists, err := unstructured.NestedInt64(configMap, "port")
		if err != nil {
			fmt.Printf("Error reading InfluxDB Port: %v\n", err)
		} else if !exists {
			fmt.Printf("InfluxDB Port not found")
		} else {
			result.Port = port
		}
		hostname, exists, err := unstructured.NestedString(configMap, "hostname")
		if err != nil {
			fmt.Printf("Error reading InfluxDB Hostname: %v\n", err)
		} else if !exists {
			fmt.Printf("InfluxDB Hostname not found")
		} else {
			result.Hostname = hostname
		}
		databaseName, exists, err := unstructured.NestedString(configMap, "databaseName")
		if err != nil {
			fmt.Printf("Error reading InfluxDB DatabaseName: %v\n", err)
		} else if !exists {
			fmt.Printf("InfluxDB DatabaseName not found")
		} else {
			result.DatabaseName = databaseName
		}
	}
	return result
}

func extractSiteWhereInstanceConfigurationPersistenceRDBConfigurations(rdbConfig interface{}) map[string]alpha3.SiteWhereInstancePersistenceRDBConfiguration {
	if configMap, ok := rdbConfig.(map[string]interface{}); ok {
		result := make(map[string]alpha3.SiteWhereInstancePersistenceRDBConfiguration)
		for key, value := range configMap {
			var configuration = extractSiteWhereInstanceConfigurationPersistenceRDBConfiguration(value)
			result[key] = configuration
		}
		return result
	}
	return nil
}

func extractSiteWhereInstanceConfigurationPersistenceRDBConfiguration(rdbConfig interface{}) alpha3.SiteWhereInstancePersistenceRDBConfiguration {
	var result = alpha3.SiteWhereInstancePersistenceRDBConfiguration{}
	return result
}

func queryMicroservices(instanceName string, clientset kubernetes.Interface) ([]alpha3.SiteWhereMicroserviceStatus, error) {
	var microservices = alpha3.GetSiteWhereMicroservicesList()
	var result = []alpha3.SiteWhereMicroserviceStatus{}
	for _, micrservice := range microservices {
		microserviceStatus, err := queryMicroserviceStatus(instanceName, &micrservice, clientset)
		if err != nil {
			return nil, err
		}
		result = append(result, microserviceStatus)
	}
	return result, nil
}

func queryMicroserviceStatus(instanceName string, microservice *alpha3.SiteWhereMicroservice, clientset kubernetes.Interface) (alpha3.SiteWhereMicroserviceStatus, error) {
	var status = "Unknown"
	deploymentName := fmt.Sprintf("%s-%s", instanceName, microservice.ID)
	deployment, err := clientset.AppsV1().Deployments(instanceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil && k8serror.IsNotFound(err) {
		return alpha3.SiteWhereMicroserviceStatus{
			Name:   microservice.ID,
			Status: "NotFound",
		}, nil
	}
	if err != nil {
		return alpha3.SiteWhereMicroserviceStatus{
			Name:   microservice.ID,
			Status: "Error",
		}, err
	}
	if deployment.Status.ReadyReplicas > 0 {
		status = "Ready"
	}
	return alpha3.SiteWhereMicroserviceStatus{
		Name:   microservice.ID,
		Status: status,
	}, nil
}
