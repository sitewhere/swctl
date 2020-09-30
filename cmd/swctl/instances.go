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

package main

import (
	"io"

	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/sitewhere/swctl/cmd/swctl/require"
	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/apis/v1/alpha3"
	"github.com/sitewhere/swctl/pkg/cli/output"
	"github.com/sitewhere/swctl/pkg/instance"
)

var instancesHelp = `
Use this command to list SiteWhere Intances.
`

func newInstancesCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewInstances(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:               "instances",
		Short:             "show SiteWhere instances",
		Long:              instancesHelp,
		Args:              require.NoArgs,
		ValidArgsFunction: noCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outFmt.Write(out, newInstancesWriter(results))
		},
	}
	bindOutputFlag(cmd, &outFmt)
	return cmd
}

type instancesWriter struct {
	// Instances found
	Instances []*alpha3.SiteWhereInstance
}

func newInstancesWriter(result *instance.ListSiteWhereInstance) *instancesWriter {
	return &instancesWriter{
		Instances: result.Instances,
	}
}

func (i *instancesWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE", "CONFIG TMPL", "DATESET TMPL", "TM STATUS", "UM STATUS")
	for _, item := range i.Instances {
		tmState := renderState(item.Status.TenantManagementStatus)
		umStatus := renderState(item.Status.UserManagementStatus)
		table.AddRow(item.Name, item.Namespace, item.ConfigurationTemplate, item.DatasetTemplate, tmState, umStatus)
	}
	return output.EncodeTable(out, table)
}

func (i *instancesWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, i)
}

func (i *instancesWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, i)
}

func renderState(state string) string {
	switch state {
	case "Unknown":
		return color.Warn.Render("Unknown")
	case "Bootstrapped":
		return color.Info.Render("Bootstrapped")
	case "NotBootstrapped":
		return color.Error.Render("Not Bootstrapped")
	default:
		return state
	}
}

// const (
// 	frmtAttr                  = "%-35s: %-32s\n"
// 	firstLevelTemplateString  = "    %-31s: %-32s\n"
// 	secondLevelTemplateFloat  = "      %-29s: %-6.2f\n"
// 	secondLevelTemplateInt    = "      %-29s: %-d\n"
// 	secondLevelTemplateBool   = "      %-29s: %-t\n"
// 	secondLevelTemplateString = "      %-29s: %-32s\n"
// 	thirdLevelTemplateString  = "        %-27s: %-32s\n"
// 	thirdLevelTemplateInt     = "        %-27s: %-d\n"
// )

// // instancesCmd represents the instances command
// var (
// 	instancesOutput = ""
// 	instancesCmd    = &cobra.Command{
// 		Use:   "instances [OPTIONS] [instance]",
// 		Short: "Manage SiteWhere Instance",
// 		Long:  `Manage SiteWhere Instance.`,
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) > 2 {
// 				return errors.New("requires one or zero arguments")
// 			}
// 			return nil
// 		},
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) < 1 {
// 				handleListInstances()
// 			} else {
// 				name := args[0]
// 				handleInstance(name)
// 			}
// 		},
// 	}
// )

// func printSiteWhereInstance(crSiteWhereInstace *unstructured.Unstructured) {
// 	sitewhereInstace := extractFromResource(crSiteWhereInstace)

// 	if strings.TrimSpace(instancesOutput) == "json" {
// 		printJSONSiteWhereInstance(sitewhereInstace)
// 	} else if strings.TrimSpace(instancesOutput) == "yaml" {
// 		printYAMLSiteWhereInstance(sitewhereInstace)
// 	} else {
// 		printStandardSiteWhereInstance(sitewhereInstace)
// 	}
// }

// func printJSONSiteWhereInstance(sitewhereInstace *alpha3.SiteWhereInstance) {
// 	e, err := json.Marshal(sitewhereInstace)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(e))
// }

// func printYAMLSiteWhereInstance(sitewhereInstace *alpha3.SiteWhereInstance) {
// 	e, err := yaml.Marshal(sitewhereInstace)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(e))
// }

// func printStandardSiteWhereInstance(sitewhereInstace *alpha3.SiteWhereInstance) {
// 	fmt.Printf(frmtAttr, "Instance Name", sitewhereInstace.Name)
// 	fmt.Printf(frmtAttr, "Instance Namespace", sitewhereInstace.Namespace)
// 	fmt.Printf(frmtAttr, "Configuration Template", sitewhereInstace.ConfigurationTemplate)
// 	fmt.Printf(frmtAttr, "Dataset Template", sitewhereInstace.DatasetTemplate)
// 	fmt.Printf(frmtAttr, "Tenant Management Status", sitewhereInstace.Status.TenantManagementStatus)
// 	fmt.Printf(frmtAttr, "User Management Status", sitewhereInstace.Status.UserManagementStatus)
// 	printSiteWhereInstanceConfiguration(sitewhereInstace.Configuration)
// 	printSiteWhereMicroservicesStatuses(sitewhereInstace.Microservices)
// }

// func printSiteWhereInstanceConfiguration(config *alpha3.SiteWhereInstanceConfiguration) {
// 	fmt.Printf("Configuration:\n")
// 	printSiteWhereInstanceConfigurationInfrastructure(config.Infrastructure)
// 	printSiteWhereInstanceConfigurationPersistence(config.Persistence)
// }

// func printSiteWhereInstanceConfigurationInfrastructure(config *alpha3.SiteWhereInstanceInfrastructureConfiguration) {
// 	fmt.Printf("  Infrastructure:\n")

// 	if config != nil {
// 		fmt.Printf(firstLevelTemplateString, "Namespace", config.Namespace)
// 		printSiteWhereInstanceConfigurationInfrastructureGRPC(config.GRPC)
// 		printSiteWhereInstanceConfigurationInfrastructureKafka(config.Kafka)
// 		printSiteWhereInstanceConfigurationInfrastructureMetrics(config.Metrics)
// 		printSiteWhereInstanceConfigurationInfrastructureRedis(config.Redis)
// 	}
// }

// func printSiteWhereInstanceConfigurationInfrastructureGRPC(config *alpha3.SiteWhereInstanceInfrastructureGRPCConfiguration) {
// 	fmt.Printf("    gRPC:\n")
// 	if config != nil {
// 		fmt.Printf(secondLevelTemplateFloat, "Backoff Multiplier", config.BackoffMultiplier)
// 		fmt.Printf(secondLevelTemplateInt, "Initial Backoff (sec)", config.InitialBackoffSeconds)
// 		fmt.Printf(secondLevelTemplateInt, "Max Backoff (sec)", config.MaxBackoffSeconds)
// 		fmt.Printf(secondLevelTemplateInt, "Max Retry", config.MaxRetryCount)
// 		fmt.Printf(secondLevelTemplateBool, "Resolve FQDN", config.ResolveFQDN)
// 	}
// }

// func printSiteWhereInstanceConfigurationInfrastructureKafka(config *alpha3.SiteWhereInstanceInfrastructureKafkaConfiguration) {
// 	fmt.Printf("    Kafka:\n")
// 	if config != nil {
// 		fmt.Printf(secondLevelTemplateString, "Hostname", config.Hostname)
// 		fmt.Printf(secondLevelTemplateInt, "Port", config.Port)
// 		fmt.Printf(secondLevelTemplateInt, "Def Topic Partitions", config.DefaultTopicPartitions)
// 		fmt.Printf(secondLevelTemplateInt, "Def Topic Replication Factor", config.DefaultTopicReplicationFactor)
// 	}
// }

// func printSiteWhereInstanceConfigurationInfrastructureMetrics(config *alpha3.SiteWhereInstanceInfrastructureMetricsConfiguration) {
// 	fmt.Printf("    Metrics:\n")
// 	if config != nil {
// 		fmt.Printf(secondLevelTemplateBool, "Enabled", config.Enabled)
// 		fmt.Printf(secondLevelTemplateInt, "HTTP Port", config.HTTPPort)
// 	}
// }

// func printSiteWhereInstanceConfigurationInfrastructureRedis(config *alpha3.SiteWhereInstanceInfrastructureRedisConfiguration) {
// 	fmt.Printf("    Redis:\n")
// 	if config != nil {
// 		fmt.Printf(secondLevelTemplateString, "Hostname", config.Hostname)
// 		fmt.Printf(secondLevelTemplateInt, "Port", config.Port)
// 		fmt.Printf(secondLevelTemplateInt, "Node Count", config.NodeCount)
// 		fmt.Printf(secondLevelTemplateString, "Master Group Name", config.MasterGroupName)
// 	}
// }

// func printSiteWhereInstanceConfigurationPersistence(config *alpha3.SiteWhereInstancePersistenceConfiguration) {
// 	fmt.Printf("  Persistence:\n")
// 	if config != nil {
// 		printSiteWhereInstanceConfigurationCassandraPersistence(config.CassandraConfigurations)
// 		printSiteWhereInstanceConfigurationInfluxDBPersistence(config.InfluxDBConfigurations)
// 		printSiteWhereInstanceConfigurationRDBPersistence(config.RDBConfigurations)
// 	}
// }

// func printSiteWhereInstanceConfigurationCassandraPersistence(config map[string]alpha3.SiteWhereInstancePersistenceCassandraConfiguration) {
// 	fmt.Printf("    Cassandra:\n")
// 	if config != nil {
// 		for key, value := range config {
// 			fmt.Printf(secondLevelTemplateString, "Entry", key)
// 			fmt.Printf(thirdLevelTemplateString, "Contact Points", value.ContactPoints)
// 			fmt.Printf(thirdLevelTemplateString, "Keyspace", value.Keyspace)
// 		}
// 	}
// }

// func printSiteWhereInstanceConfigurationInfluxDBPersistence(config map[string]alpha3.SiteWhereInstancePersistenceInfluxDBConfiguration) {
// 	fmt.Printf("    InfluxDB:\n")
// 	if config != nil {
// 		for key, value := range config {
// 			fmt.Printf(secondLevelTemplateString, "Entry", key)
// 			fmt.Printf(thirdLevelTemplateString, "Hostname", value.Hostname)
// 			fmt.Printf(thirdLevelTemplateInt, "Port", value.Port)
// 			fmt.Printf(thirdLevelTemplateString, "Database Name", value.DatabaseName)
// 		}
// 	}
// }

// func printSiteWhereInstanceConfigurationRDBPersistence(config map[string]alpha3.SiteWhereInstancePersistenceRDBConfiguration) {
// 	fmt.Printf("    RDB:\n")
// }

// func printSiteWhereMicroservicesStatuses(microservices []alpha3.SiteWhereMicroserviceStatus) {
// 	fmt.Printf("  Microservices:\n")
// 	for _, micrservice := range microservices {
// 		fmt.Printf(secondLevelTemplateString, micrservice.Name, micrservice.Status)
// 	}
// }
