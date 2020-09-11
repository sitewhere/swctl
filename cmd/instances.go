/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sitewhere/swctl/internal"
	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var (
	sitewhereInstanceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha3",
		Resource: "instances",
	}
)

var (
	sitewhereMicroserviceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha3",
		Resource: "microservices",
	}
)

const (
	frmtAttr                  = "%-35s: %-32s\n"
	firstLevelTemplateString  = "    %-31s: %-32s\n"
	secondLevelTemplateFloat  = "      %-29s: %-6.2f\n"
	secondLevelTemplateInt    = "      %-29s: %-d\n"
	secondLevelTemplateBool   = "      %-29s: %-t\n"
	secondLevelTemplateString = "      %-29s: %-32s\n"
	thirdLevelTemplateString  = "        %-27s: %-32s\n"
	thirdLevelTemplateInt     = "        %-27s: %-d\n"
)

// instancesCmd represents the instances command
var (
	instancesOutput = ""
	instancesCmd    = &cobra.Command{
		Use:   "instances [OPTIONS] [instance]",
		Short: "Manage SiteWhere Instance",
		Long:  `Manage SiteWhere Instance.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 2 {
				return errors.New("requires one or zero arguments")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				handleListInstances()
			} else {
				name := args[0]
				handleInstance(name)
			}
		},
	}
)

func init() {
	instancesCmd.Flags().StringVarP(&instancesOutput, "output", "o", "", "Output format. One of 'yaml' or 'json'.")
	// rootCmd.AddCommand(instancesCmd)
}

func handleListInstances() {
	var err error

	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes Config: %v\n", err)
		return
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return
	}

	res := client.Resource(sitewhereInstanceGVR)
	options := metav1.ListOptions{}
	sitewhereInstaces, err := res.List(context.TODO(), options)

	if err != nil {
		fmt.Printf("Error reading SiteWhere Instances: %v\n", err)
		return
	}

	template := "%-20s%-20s%-20s%-20s%-20s%-20s\n"
	fmt.Printf(template, "NAME", "NAMESPACE", "CONFIG TMPL", "DATESET TMPL", "TM STATUS", "UM STATUS")

	for _, instance := range sitewhereInstaces.Items {
		sitewhereInstace := extractFromResource(&instance)
		fmt.Printf(
			template,
			sitewhereInstace.Name,
			sitewhereInstace.Namespace,
			sitewhereInstace.ConfigurationTemplate,
			sitewhereInstace.DatasetTemplate,
			sitewhereInstace.Status.TenantManagementStatus,
			sitewhereInstace.Status.UserManagementStatus,
		)
	}
}

func handleInstance(instanceName string) {
	var err error

	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes Config: %v\n", err)
		return
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return
	}

	res := client.Resource(sitewhereInstanceGVR)
	options := metav1.GetOptions{}
	sitewhereInstace, err := res.Get(context.TODO(), instanceName, options)

	if err != nil {
		fmt.Printf(
			"SiteWhere Instace %s NOT FOUND.\n",
			instanceName)
		return
	}

	printSiteWhereInstance(sitewhereInstace)
}

func printSiteWhereInstance(crSiteWhereInstace *unstructured.Unstructured) {
	sitewhereInstace := extractFromResource(crSiteWhereInstace)

	if strings.TrimSpace(instancesOutput) == "json" {
		printJSONSiteWhereInstance(sitewhereInstace)
	} else if strings.TrimSpace(instancesOutput) == "yaml" {
		printYAMLSiteWhereInstance(sitewhereInstace)
	} else {
		printStandardSiteWhereInstance(sitewhereInstace)
	}
}

func printJSONSiteWhereInstance(sitewhereInstace *alpha3.SiteWhereInstance) {
	e, err := json.Marshal(sitewhereInstace)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}

func printYAMLSiteWhereInstance(sitewhereInstace *alpha3.SiteWhereInstance) {
	e, err := yaml.Marshal(sitewhereInstace)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}

func printStandardSiteWhereInstance(sitewhereInstace *alpha3.SiteWhereInstance) {
	fmt.Printf(frmtAttr, "Instance Name", sitewhereInstace.Name)
	fmt.Printf(frmtAttr, "Instance Namespace", sitewhereInstace.Namespace)
	fmt.Printf(frmtAttr, "Configuration Template", sitewhereInstace.ConfigurationTemplate)
	fmt.Printf(frmtAttr, "Dataset Template", sitewhereInstace.DatasetTemplate)
	fmt.Printf(frmtAttr, "Tenant Management Status", sitewhereInstace.Status.TenantManagementStatus)
	fmt.Printf(frmtAttr, "User Management Status", sitewhereInstace.Status.UserManagementStatus)
	printSiteWhereInstanceConfiguration(sitewhereInstace.Configuration)
	printSiteWhereMicroservicesStatuses(sitewhereInstace.Microservices)
}

func printSiteWhereInstanceConfiguration(config *alpha3.SiteWhereInstanceConfiguration) {
	fmt.Printf("Configuration:\n")
	printSiteWhereInstanceConfigurationInfrastructure(config.Infrastructure)
	printSiteWhereInstanceConfigurationPersistence(config.Persistence)
}

func printSiteWhereInstanceConfigurationInfrastructure(config *alpha3.SiteWhereInstanceInfrastructureConfiguration) {
	fmt.Printf("  Infrastructure:\n")

	if config != nil {
		fmt.Printf(firstLevelTemplateString, "Namespace", config.Namespace)
		printSiteWhereInstanceConfigurationInfrastructureGRPC(config.GRPC)
		printSiteWhereInstanceConfigurationInfrastructureKafka(config.Kafka)
		printSiteWhereInstanceConfigurationInfrastructureMetrics(config.Metrics)
		printSiteWhereInstanceConfigurationInfrastructureRedis(config.Redis)
	}
}

func printSiteWhereInstanceConfigurationInfrastructureGRPC(config *alpha3.SiteWhereInstanceInfrastructureGRPCConfiguration) {
	fmt.Printf("    gRPC:\n")
	if config != nil {
		fmt.Printf(secondLevelTemplateFloat, "Backoff Multiplier", config.BackoffMultiplier)
		fmt.Printf(secondLevelTemplateInt, "Initial Backoff (sec)", config.InitialBackoffSeconds)
		fmt.Printf(secondLevelTemplateInt, "Max Backoff (sec)", config.MaxBackoffSeconds)
		fmt.Printf(secondLevelTemplateInt, "Max Retry", config.MaxRetryCount)
		fmt.Printf(secondLevelTemplateBool, "Resolve FQDN", config.ResolveFQDN)
	}
}

func printSiteWhereInstanceConfigurationInfrastructureKafka(config *alpha3.SiteWhereInstanceInfrastructureKafkaConfiguration) {
	fmt.Printf("    Kafka:\n")
	if config != nil {
		fmt.Printf(secondLevelTemplateString, "Hostname", config.Hostname)
		fmt.Printf(secondLevelTemplateInt, "Port", config.Port)
		fmt.Printf(secondLevelTemplateInt, "Def Topic Partitions", config.DefaultTopicPartitions)
		fmt.Printf(secondLevelTemplateInt, "Def Topic Replication Factor", config.DefaultTopicReplicationFactor)
	}
}

func printSiteWhereInstanceConfigurationInfrastructureMetrics(config *alpha3.SiteWhereInstanceInfrastructureMetricsConfiguration) {
	fmt.Printf("    Metrics:\n")
	if config != nil {
		fmt.Printf(secondLevelTemplateBool, "Enabled", config.Enabled)
		fmt.Printf(secondLevelTemplateInt, "HTTP Port", config.HTTPPort)
	}
}

func printSiteWhereInstanceConfigurationInfrastructureRedis(config *alpha3.SiteWhereInstanceInfrastructureRedisConfiguration) {
	fmt.Printf("    Redis:\n")
	if config != nil {
		fmt.Printf(secondLevelTemplateString, "Hostname", config.Hostname)
		fmt.Printf(secondLevelTemplateInt, "Port", config.Port)
		fmt.Printf(secondLevelTemplateInt, "Node Count", config.NodeCount)
		fmt.Printf(secondLevelTemplateString, "Master Group Name", config.MasterGroupName)
	}
}

func printSiteWhereInstanceConfigurationPersistence(config *alpha3.SiteWhereInstancePersistenceConfiguration) {
	fmt.Printf("  Persistence:\n")
	if config != nil {
		printSiteWhereInstanceConfigurationCassandraPersistence(config.CassandraConfigurations)
		printSiteWhereInstanceConfigurationInfluxDBPersistence(config.InfluxDBConfigurations)
		printSiteWhereInstanceConfigurationRDBPersistence(config.RDBConfigurations)
	}
}

func printSiteWhereInstanceConfigurationCassandraPersistence(config map[string]alpha3.SiteWhereInstancePersistenceCassandraConfiguration) {
	fmt.Printf("    Cassandra:\n")
	if config != nil {
		for key, value := range config {
			fmt.Printf(secondLevelTemplateString, "Entry", key)
			fmt.Printf(thirdLevelTemplateString, "Contact Points", value.ContactPoints)
			fmt.Printf(thirdLevelTemplateString, "Keyspace", value.Keyspace)
		}
	}
}

func printSiteWhereInstanceConfigurationInfluxDBPersistence(config map[string]alpha3.SiteWhereInstancePersistenceInfluxDBConfiguration) {
	fmt.Printf("    InfluxDB:\n")
	if config != nil {
		for key, value := range config {
			fmt.Printf(secondLevelTemplateString, "Entry", key)
			fmt.Printf(thirdLevelTemplateString, "Hostname", value.Hostname)
			fmt.Printf(thirdLevelTemplateInt, "Port", value.Port)
			fmt.Printf(thirdLevelTemplateString, "Database Name", value.DatabaseName)
		}
	}
}

func printSiteWhereInstanceConfigurationRDBPersistence(config map[string]alpha3.SiteWhereInstancePersistenceRDBConfiguration) {
	fmt.Printf("    RDB:\n")
}

func printSiteWhereMicroservicesStatuses(microservices []alpha3.SiteWhereMicroserviceStatus) {
	fmt.Printf("  Microservices:\n")
	for _, micrservice := range microservices {
		fmt.Printf(secondLevelTemplateString, micrservice.Name, micrservice.Status)
	}
}

func extractFromResource(crSiteWhereInstace *unstructured.Unstructured) *alpha3.SiteWhereInstance {
	var result = alpha3.SiteWhereInstance{}

	metadata, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "metadata")
	if err != nil {
		fmt.Printf("Error reading metadata for %s: %v\n", crSiteWhereInstace, err)
		return nil
	}
	if !exists {
		fmt.Printf("Metadata not found for for SiteWhere Instance: %s", crSiteWhereInstace)
	} else {
		extractSiteWhereInstanceMetadata(metadata, &result)
	}

	spec, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "spec")
	if err != nil {
		fmt.Printf("Error reading spec for %s: %v\n", result.Name, err)
		return nil
	}
	if !exists {
		fmt.Printf("Spec not found for for SiteWhere Instance: %s", result.Name)
	} else {
		extractSiteWhereInstanceSpec(spec, &result)
	}

	status, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "status")
	if err != nil {
		fmt.Printf("Error reading status for %s: %v\n", result.Name, err)
		return nil
	}
	if !exists {
		result.Status = &alpha3.SiteWhereInstanceStatus{
			TenantManagementStatus: "Unknown",
			UserManagementStatus:   "Unknown",
		}
	} else {
		extractSiteWhereInstanceStatus(status, &result)
	}

	microservices, err := queryMicroservices(result.Name)

	if err != nil {
		fmt.Printf("Error reading microservices statuses for %s: %v\n", result.Name, err)
		return nil
	}

	result.Microservices = microservices

	return &result
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
	instanceNamespace := spec["instanceNamespace"]
	configurationTemplate := spec["configurationTemplate"]
	datasetTemplate := spec["datasetTemplate"]
	configuration := spec["configuration"]
	sitewhereConfiguration := extractSiteWhereInstanceConfiguration(configuration)

	instance.Namespace = fmt.Sprintf("%v", instanceNamespace)
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

func queryMicroservices(instanceName string) ([]alpha3.SiteWhereMicroserviceStatus, error) {
	var microservices = alpha3.GetSiteWhereMicroservicesList()
	var result = []alpha3.SiteWhereMicroserviceStatus{}

	config, err := internal.GetKubeConfigFromKubeconfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes Config: %v\n", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting Kubernetes Client: %v\n", err)
		return nil, err
	}

	for _, micrservice := range microservices {

		microserviceStatus, err := queryMicroserviceStatus(instanceName, &micrservice, clientset)

		if err != nil {
			return nil, err
		}

		result = append(result, microserviceStatus)
	}

	return result, nil
}

func queryMicroserviceStatus(instanceName string, microservice *alpha3.SiteWhereMicroservice, clientset *kubernetes.Clientset) (alpha3.SiteWhereMicroserviceStatus, error) {
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
