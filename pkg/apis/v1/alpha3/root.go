/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package alpha3 defines SiteWhere Structures
package alpha3

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

// SiteWhereInstancePersistenceCassandraConfiguration SiteWhere Instance Persistence Cassandra configurations
type SiteWhereInstancePersistenceCassandraConfiguration struct {
	ContactPoints string `json:"contactPoints"`
	Keyspace      string `json:"keyspace"`
}

// SiteWhereInstancePersistenceInfluxDBConfiguration SiteWhere Instance Persistence InfuxDB configurations
type SiteWhereInstancePersistenceInfluxDBConfiguration struct {
	Hostname     string `json:"hostname"`
	Port         int64  `json:"port"`
	DatabaseName string `json:"databaseName"`
}

// SiteWhereInstancePersistenceRDBConfiguration SiteWhere Instance Persistence Relational Database configurations
type SiteWhereInstancePersistenceRDBConfiguration struct {
}

// SiteWhereInstancePersistenceConfiguration SiteWhere Instance Persistence configurations
type SiteWhereInstancePersistenceConfiguration struct {
	CassandraConfigurations map[string]SiteWhereInstancePersistenceCassandraConfiguration `json:"cassandraConfigurations"`
	InfluxDBConfigurations  map[string]SiteWhereInstancePersistenceInfluxDBConfiguration  `json:"influxDbConfigurations"`
	RDBConfigurations       map[string]SiteWhereInstancePersistenceRDBConfiguration       `json:"rdbConfigurations"`
}

// SiteWhereInstanceConfiguration SiteWhere Instance configurations
type SiteWhereInstanceConfiguration struct {
	Infrastructure *SiteWhereInstanceInfrastructureConfiguration `json:"infrastructure"`
	Persistence    *SiteWhereInstancePersistenceConfiguration    `json:"persistenceConfigurations"`
}

// SiteWhereInstanceStatus SiteWhere Instance Tenant Management and User Management Status
type SiteWhereInstanceStatus struct {
	TenantManagementStatus string `json:"tenantManagementStatus"`
	UserManagementStatus   string `json:"userManagementStatus"`
}

// SiteWhereMicroserviceStatus SiteWhere Instance Microservice Status
type SiteWhereMicroserviceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// SiteWhereInstance represents an Instacen in SiteWhere
type SiteWhereInstance struct {
	Name                  string                          `json:"name"`
	Namespace             string                          `json:"namespace"`
	ConfigurationTemplate string                          `json:"configurationTemaplate"`
	DatasetTemplate       string                          `json:"datasetTemplate"`
	Configuration         *SiteWhereInstanceConfiguration `json:"configuration"`
	Status                *SiteWhereInstanceStatus        `json:"status"`
	Microservices         []SiteWhereMicroserviceStatus   `json:"microservices"`
}

var microservices = []string{
	"asset-management",
	"batch-operations",
	"command-delivery",
	"device-management",
	"device-registration",
	"device-state",
	"event-management",
	"event-sources",
	"inbound-processing",
	"instance-management",
	"label-generation",
	"outbound-connectors",
	"schedule-management",
}

// GetSiteWhereMicroservicesList Returns the list of SiteWhere Microservices Names
func GetSiteWhereMicroservicesList() []string {
	return microservices
}
