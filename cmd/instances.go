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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// SiteWhereInstanceInfrastructureGRPCConfiguration SiteWhere Instance Infrastructure gRPC configurations
type SiteWhereInstanceInfrastructureGRPCConfiguration struct {
	BackoffMultiplier     float64 `json:"backoffMultiplier"`
	InitialBackoffSeconds int64   `json:"initialBackoffSeconds"`
	MaxBackoffSeconds     int64   `json:"maxBackoffSeconds"`
	MaxRetryCount         int64   `json:"maxRetryCount"`
	ResolveFQDN           bool    `json:"resolveFQDN"`
}

// SiteWhereInstanceInfrastructureKafkaConfiguration SiteWhere Instance Infrastrucre Kafka configurations
type SiteWhereInstanceInfrastructureKafkaConfiguration struct {
	Hostname                      string `json:"hostname"`
	Port                          int64  `json:"port"`
	DefaultTopicPartitions        int64  `json:"defaultTopicPartitions"`
	DefaultTopicReplicationFactor int64  `json:"defaultTopicReplicationFactor"`
}

// SiteWhereInstanceInfrastructureMetricsConfiguration SiteWhere Instance Infrastrucre Metrics configurations
type SiteWhereInstanceInfrastructureMetricsConfiguration struct {
	Enabled  bool  `json:"enabled"`
	HTTPPort int64 `json:"httpPort"`
}

// SiteWhereInstanceInfrastructureRedisConfiguration SiteWhere Instance Infrastrucre Redis configurations
type SiteWhereInstanceInfrastructureRedisConfiguration struct {
	Hostname        string `json:"hostname"`
	Port            int64  `json:"port"`
	NodeCount       int64  `json:"nodeCount"`
	MasterGroupName string `json:"masterGroupName"`
}

// SiteWhereInstanceInfrastructureConfiguration SiteWhere Instance Infrastructure configurations
type SiteWhereInstanceInfrastructureConfiguration struct {
	Namespace string                                               `json:"namespace"`
	GRPC      *SiteWhereInstanceInfrastructureGRPCConfiguration    `json:"grpc"`
	Kafka     *SiteWhereInstanceInfrastructureKafkaConfiguration   `json:"kafka"`
	Metrics   *SiteWhereInstanceInfrastructureMetricsConfiguration `json:"metrics"`
	Redis     *SiteWhereInstanceInfrastructureRedisConfiguration   `json:"redis"`
}

// SiteWhereInstancePersistenceConfiguration SiteWhere Instance Persistence configurations
type SiteWhereInstancePersistenceConfiguration struct {
}

// SiteWhereInstanceConfiguration SiteWhere Instance configurations
type SiteWhereInstanceConfiguration struct {
	Infrastructure *SiteWhereInstanceInfrastructureConfiguration `json:"infrastructure"`
	Persistence    *SiteWhereInstancePersistenceConfiguration    `json:"persistenceConfigurations"`
}

// SiteWhereInstance represents an Instacen in SiteWhere
type SiteWhereInstance struct {
	Name                  string                          `json:"name"`
	Namespace             string                          `json:"namespace"`
	ConfigurationTemplate string                          `json:"configurationTemaplate"`
	DatasetTemplate       string                          `json:"datasetTemplate"`
	Configuration         *SiteWhereInstanceConfiguration `json:"configuration"`
}

var clientset *kubernetes.Clientset
var apixClient *apixv1beta1client.ApiextensionsV1beta1Client
var (
	sitewhereInstanceGVR = schema.GroupVersionResource{
		Group:    "sitewhere.io",
		Version:  "v1alpha3",
		Resource: "instances",
	}
)

const frmtAttr = "%-35s: %-32s\n"

const firstLevelTemplateString = "    %-31s: %-32s\n"

const secondLevelTemplateFloat = "      %-29s: %-6.2f\n"
const secondLevelTemplateInt = "      %-29s: %-d\n"
const secondLevelTemplateBool = "      %-29s: %-t\n"
const secondLevelTemplateString = "      %-29s: %-32s\n"

// instancesCmd represents the instances command
var (
	output       = ""
	instancesCmd = &cobra.Command{
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
	instancesCmd.Flags().StringVarP(&output, "output", "o", "", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(instancesCmd)
}

func handleListInstances() {
	var err error

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := getKubeConfig(kubeconfig)
	if err != nil {
		log.Printf("Error getting Kubernetes Config: %v", err)
		return
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("Error getting Kubernetes Client: %v", err)
		return
	}

	res := client.Resource(sitewhereInstanceGVR)
	options := metav1.ListOptions{}
	sitewhereInstaces, err := res.List(context.TODO(), options)

	if err != nil {
		log.Printf("Error reading SiteWhere Instances: %v", err)
		return
	}

	template := "%-20s%-20s%-20s%-20s\n"
	fmt.Printf(template, "NAME", "NAMESPACE", "CONFIG TMPL", "DATESET TMPL")

	for _, instance := range sitewhereInstaces.Items {
		sitewhereInstace := extractFromResource(&instance)
		fmt.Printf(
			template,
			sitewhereInstace.Name,
			sitewhereInstace.Namespace,
			sitewhereInstace.ConfigurationTemplate,
			sitewhereInstace.DatasetTemplate)
	}
}

func handleInstance(instanceName string) {
	var err error

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := getKubeConfig(kubeconfig)
	if err != nil {
		log.Printf("Error getting Kubernetes Config: %v", err)
		return
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("Error getting Kubernetes Client: %v", err)
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

	if strings.TrimSpace(output) == "json" {
		printJSONSiteWhereInstance(sitewhereInstace)
	} else if strings.TrimSpace(output) == "yaml" {
		printYAMLSiteWhereInstance(sitewhereInstace)
	} else {
		printStandardSiteWhereInstance(sitewhereInstace)
	}
}

func printJSONSiteWhereInstance(sitewhereInstace *SiteWhereInstance) {
	e, err := json.Marshal(sitewhereInstace)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}

func printYAMLSiteWhereInstance(sitewhereInstace *SiteWhereInstance) {
	e, err := yaml.Marshal(sitewhereInstace)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}

func printStandardSiteWhereInstance(sitewhereInstace *SiteWhereInstance) {
	fmt.Printf(frmtAttr, "Instance Name", sitewhereInstace.Name)
	fmt.Printf(frmtAttr, "Instance Namespace", sitewhereInstace.Namespace)
	fmt.Printf(frmtAttr, "Configuration Template", sitewhereInstace.ConfigurationTemplate)
	fmt.Printf(frmtAttr, "Dataset Template", sitewhereInstace.DatasetTemplate)
	printSiteWhereInstanceConfiguration(sitewhereInstace.Configuration)
}

func printSiteWhereInstanceConfiguration(config *SiteWhereInstanceConfiguration) {
	fmt.Printf("Configuration:\n")
	printSiteWhereInstanceConfigurationInfrastructure(config.Infrastructure)
	printSiteWhereInstanceConfigurationPersistence(config.Persistence)
}

func printSiteWhereInstanceConfigurationInfrastructure(config *SiteWhereInstanceInfrastructureConfiguration) {
	fmt.Printf("  Infrastructure:\n")

	fmt.Printf(firstLevelTemplateString, "Namespace", config.Namespace)
	printSiteWhereInstanceConfigurationInfrastructureGRPC(config.GRPC)
	printSiteWhereInstanceConfigurationInfrastructureKafka(config.Kafka)
	printSiteWhereInstanceConfigurationInfrastructureMetrics(config.Metrics)
	printSiteWhereInstanceConfigurationInfrastructureRedis(config.Redis)
}

func printSiteWhereInstanceConfigurationInfrastructureGRPC(config *SiteWhereInstanceInfrastructureGRPCConfiguration) {
	fmt.Printf("    gRPC:\n")
	fmt.Printf(secondLevelTemplateFloat, "Backoff Multiplier", config.BackoffMultiplier)
	fmt.Printf(secondLevelTemplateInt, "Initial Backoff (sec)", config.InitialBackoffSeconds)
	fmt.Printf(secondLevelTemplateInt, "Max Backoff (sec)", config.MaxBackoffSeconds)
	fmt.Printf(secondLevelTemplateInt, "Max Retry", config.MaxRetryCount)
	fmt.Printf(secondLevelTemplateBool, "Resolve FQDN", config.ResolveFQDN)
}

func printSiteWhereInstanceConfigurationInfrastructureKafka(config *SiteWhereInstanceInfrastructureKafkaConfiguration) {
	fmt.Printf("    Kafka:\n")
	fmt.Printf(secondLevelTemplateString, "Hostname", config.Hostname)
	fmt.Printf(secondLevelTemplateInt, "Port", config.Port)
	fmt.Printf(secondLevelTemplateInt, "Def Topic Partitions", config.DefaultTopicPartitions)
	fmt.Printf(secondLevelTemplateInt, "Def Topic Replication Factor", config.DefaultTopicReplicationFactor)
}

func printSiteWhereInstanceConfigurationInfrastructureMetrics(config *SiteWhereInstanceInfrastructureMetricsConfiguration) {
	fmt.Printf("    Metrics:\n")
	fmt.Printf(secondLevelTemplateBool, "Enabled", config.Enabled)
	fmt.Printf(secondLevelTemplateInt, "HTTP Port", config.HTTPPort)
}

func printSiteWhereInstanceConfigurationInfrastructureRedis(config *SiteWhereInstanceInfrastructureRedisConfiguration) {
	fmt.Printf("    Redis:\n")
	fmt.Printf(secondLevelTemplateString, "Hostname", config.Hostname)
	fmt.Printf(secondLevelTemplateInt, "Port", config.Port)
	fmt.Printf(secondLevelTemplateInt, "Node Count", config.NodeCount)
	fmt.Printf(secondLevelTemplateString, "Master Group Name", config.MasterGroupName)
}

func printSiteWhereInstanceConfigurationPersistence(config *SiteWhereInstancePersistenceConfiguration) {
	fmt.Printf("  Persistence:\n")
}

func extractFromResource(crSiteWhereInstace *unstructured.Unstructured) *SiteWhereInstance {
	metadata, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "metadata")
	var result = SiteWhereInstance{}

	if err != nil {
		log.Printf("Error reading metadata for %s: %v", crSiteWhereInstace, err)
		return nil
	}
	if !exists {
		log.Printf("Metadata not found for for SiteWhere Instance: %s", crSiteWhereInstace)
	} else {
		extractSiteWhereInstanceMetadata(metadata, &result)
	}
	spec, exists, err := unstructured.NestedMap(crSiteWhereInstace.Object, "spec")
	if err != nil {
		log.Printf("Error reading spec for %s: %v", crSiteWhereInstace, err)
		return nil
	}
	if !exists {
		log.Printf("Spec not found for for SiteWhere Instance: %s", crSiteWhereInstace)
	} else {
		extractSiteWhereInstanceSpec(spec, &result)
	}

	return &result
}

func extractSiteWhereInstanceMetadata(metadata map[string]interface{}, instance *SiteWhereInstance) {
	name := extractSiteWhereInstanceName(metadata)
	instance.Name = fmt.Sprintf("%v", name)
}

func extractSiteWhereInstanceName(metadata map[string]interface{}) interface{} {
	return metadata["name"]
}

func extractSiteWhereInstanceSpec(spec map[string]interface{}, instance *SiteWhereInstance) {
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

func extractSiteWhereInstanceConfiguration(config interface{}) *SiteWhereInstanceConfiguration {
	var result = SiteWhereInstanceConfiguration{}

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

func extractSiteWhereInstanceConfigurationInfrastructure(infrastructureConfig interface{}) *SiteWhereInstanceInfrastructureConfiguration {
	var result = SiteWhereInstanceInfrastructureConfiguration{}

	if configMap, ok := infrastructureConfig.(map[string]interface{}); ok {
		namespace, exists, err := unstructured.NestedString(configMap, "namespace")
		if err != nil {
			log.Printf("Error reading Infrastructure Namespace: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Infrastructure Namespace not found")
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

func extractSiteWhereInstanceConfigurationInfrastructureGRPC(gRPCConfig interface{}) *SiteWhereInstanceInfrastructureGRPCConfiguration {
	var result = SiteWhereInstanceInfrastructureGRPCConfiguration{}

	if configMap, ok := gRPCConfig.(map[string]interface{}); ok {
		backoffMultiplier, exists, err := unstructured.NestedFloat64(configMap, "backoffMultiplier")
		if err != nil {
			log.Printf("Error reading backoffMultiplier: %v", err)
			return nil
		}
		if !exists {
			log.Printf("backoffMultiplier not found")
		} else {
			result.BackoffMultiplier = backoffMultiplier
		}

		initialBackoffSeconds, exists, err := unstructured.NestedInt64(configMap, "initialBackoffSeconds")
		if err != nil {
			log.Printf("Error reading initialBackoffSeconds: %v", err)
			return nil
		}
		if !exists {
			log.Printf("initialBackoffSeconds not found")
		} else {
			result.InitialBackoffSeconds = initialBackoffSeconds
		}

		maxBackoffSeconds, exists, err := unstructured.NestedInt64(configMap, "maxBackoffSeconds")
		if err != nil {
			log.Printf("Error reading maxBackoffSeconds: %v", err)
			return nil
		}
		if !exists {
			log.Printf("maxBackoffSeconds not found")
		} else {
			result.MaxBackoffSeconds = maxBackoffSeconds
		}

		maxRetryCount, exists, err := unstructured.NestedInt64(configMap, "maxRetryCount")
		if err != nil {
			log.Printf("Error reading maxRetryCount: %v", err)
			return nil
		}
		if !exists {
			log.Printf("maxRetryCount not found")
		} else {
			result.MaxRetryCount = maxRetryCount
		}

		resolveFQDN, exists, err := unstructured.NestedBool(configMap, "resolveFQDN")
		if err != nil {
			log.Printf("Error reading resolveFQDN: %v", err)
			return nil
		}
		if !exists {
			log.Printf("resolveFQDN not found")
		} else {
			result.ResolveFQDN = resolveFQDN
		}
	}

	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureKafka(kafkaConfig interface{}) *SiteWhereInstanceInfrastructureKafkaConfiguration {
	var result = SiteWhereInstanceInfrastructureKafkaConfiguration{}

	if configMap, ok := kafkaConfig.(map[string]interface{}); ok {
		port, exists, err := unstructured.NestedInt64(configMap, "port")
		if err != nil {
			log.Printf("Error reading Kafka Port: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Kafka Port not found")
		} else {
			result.Port = port
		}

		hostname, exists, err := unstructured.NestedString(configMap, "hostname")
		if err != nil {
			log.Printf("Error reading Kafka Hostname: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Kafka Hostname not found")
		} else {
			result.Hostname = hostname
		}

		defaultTopicPartitions, exists, err := unstructured.NestedInt64(configMap, "defaultTopicPartitions")
		if err != nil {
			log.Printf("Error reading Kafka defaultTopicPartitions: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Kafka defaultTopicPartitions not found")
		} else {
			result.DefaultTopicPartitions = defaultTopicPartitions
		}

		defaultTopicReplicationFactor, exists, err := unstructured.NestedInt64(configMap, "defaultTopicReplicationFactor")
		if err != nil {
			log.Printf("Error reading Kafka defaultTopicReplicationFactor: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Kafka defaultTopicReplicationFactor not found")
		} else {
			result.DefaultTopicReplicationFactor = defaultTopicReplicationFactor
		}
	}

	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureMetrics(metricsConfig interface{}) *SiteWhereInstanceInfrastructureMetricsConfiguration {
	var result = SiteWhereInstanceInfrastructureMetricsConfiguration{}

	if configMap, ok := metricsConfig.(map[string]interface{}); ok {
		enabled, exists, err := unstructured.NestedBool(configMap, "enabled")
		if err != nil {
			log.Printf("Error reading Metrics Enabled: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Metrics Enabled not found")
		} else {
			result.Enabled = enabled
		}

		httpPort, exists, err := unstructured.NestedInt64(configMap, "httpPort")
		if err != nil {
			log.Printf("Error reading Metrics HTTP Port: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Metrics HTTP Port not found")
		} else {
			result.HTTPPort = httpPort
		}

	}
	return &result
}

func extractSiteWhereInstanceConfigurationInfrastructureRedis(redisConfig interface{}) *SiteWhereInstanceInfrastructureRedisConfiguration {
	var result = SiteWhereInstanceInfrastructureRedisConfiguration{}

	if configMap, ok := redisConfig.(map[string]interface{}); ok {
		hostname, exists, err := unstructured.NestedString(configMap, "hostname")
		if err != nil {
			log.Printf("Error reading Redis Hostname: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Redis Hostname not found")
		} else {
			result.Hostname = hostname
		}

		port, exists, err := unstructured.NestedInt64(configMap, "port")
		if err != nil {
			log.Printf("Error reading Redis Port: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Redis Port not found")
		} else {
			result.Port = port
		}

		nodeCount, exists, err := unstructured.NestedInt64(configMap, "nodeCount")
		if err != nil {
			log.Printf("Error reading Redis Node Count: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Redis Node Count not found")
		} else {
			result.NodeCount = nodeCount
		}

		masterGroupName, exists, err := unstructured.NestedString(configMap, "masterGroupName")
		if err != nil {
			log.Printf("Error reading Redis Master Group Name: %v", err)
			return nil
		}
		if !exists {
			log.Printf("Redis Master Group Name not found")
		} else {
			result.MasterGroupName = masterGroupName
		}
	}
	return &result
}

func extractSiteWhereInstanceConfigurationPersistenceConfiguration(persistenceConfig interface{}) *SiteWhereInstancePersistenceConfiguration {
	var result = SiteWhereInstancePersistenceConfiguration{}

	// if configMap, ok := persistenceConfig.(map[string]interface{}); ok {
	// }
	return &result
}

// Buid a Kubernetes Config from a filepath
func getKubeConfig(pathToCfg string) (*rest.Config, error) {
	if pathToCfg == "" {
		// in cluster access
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", pathToCfg)
}

func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(pathToCfg)

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getClientV1Beta1(pathToCfg string) (*apixv1beta1client.ApiextensionsV1beta1Client, error) {
	config, err := getKubeConfig(pathToCfg)

	if err != nil {
		return nil, err
	}

	return apixv1beta1client.NewForConfig(config)
}
