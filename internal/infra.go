/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

// Package internal Implements swctl internal use only functions
package internal

import (
	"fmt"

	"github.com/gookit/color"

	"k8s.io/apimachinery/pkg/api/errors"
)

var infraTemplate = []InfraTemplateResource{
	InfraTemplateResource{
		FileTemplate: "/infra-min/01-mosquitto/infra-mosquitto-%02d.yaml",
		FileCount:    4,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/02-redis/infra-redis-%02d.yaml",
		FileCount:    7,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/03-zookeeper/infra-zookeeper-%02d.yaml",
		FileCount:    5,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/04-kafka/infra-kafka-%02d.yaml",
		FileCount:    1,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/05-postgresql/infra-postgresql-%02d.yaml",
		FileCount:    3,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/06-syncope/infra-syncope-%02d.yaml",
		FileCount:    4,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/07-warp10/infra-warp10-%02d.yaml",
		FileCount:    4,
		Enabled:      false,
	},
	InfraTemplateResource{
		FileTemplate: "/infra-min/08-influxdb/infra-influxdb-%02d.yaml",
		FileCount:    4,
		Enabled:      true,
	},
}

// InstallSiteWhereInfrastructure Install SiteWhere Infrastructure components in the cluster
func InstallSiteWhereInfrastructure(config SiteWhereConfiguration) error {
	var err error

	for _, tpl := range infraTemplate {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = InstallResourceFromFile(infraResource, config)
				if err != nil && !errors.IsAlreadyExists(err) {
					return err
				}
			}
		}
	}

	err = waitForDeploymentAvailable(config.GetClientset(), "sitewhere-infrastructure-mosquitto", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Deployment sitewhere-infrastructure-mosquitto: ")
		color.Info.Println("Available")
	}

	err = waitForDeploymentAvailable(config.GetClientset(), "sitewhere-kafka-entity-operator", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Deployment sitewhere-kafka-entity-operator: ")
		color.Info.Println("Available")
	}

	err = waitForDeploymentAvailable(config.GetClientset(), "sitewhere-syncope", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Deployment sitewhere-syncope: ")
		color.Info.Println("Available")
	}

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-infrastructure-zookeeper-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Pod sitewhere-infrastructure-zookeeper-0: ")
		color.Info.Println("Ready")
	}
	// TODO if not minimal, wait for -1 and -2

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-kafka-kafka-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Pod sitewhere-kafka-kafka-0: ")
		color.Info.Println("Ready")
	}

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-postgresql-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Pod sitewhere-postgresql-0: ")
		color.Info.Println("Ready")
	}

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-infrastructure-redis-ha-server-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Print("Pod sitewhere-infrastructure-redis-ha-server-0: ")
		color.Info.Println("Ready")
	}

	// err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-infrastructure-warp10-0", sitewhereSystemNamespace)
	// if err != nil {
	// 	return err
	// }
	// if config.IsVerbose() {
	// 	fmt.Printf("Pod sitewhere-infrastructure-warp10-0: Ready\n")
	// }

	if config.IsVerbose() {
		color.Info.Println("SiteWhere Infrastructure: Installed")
	}
	return nil
}

// UninstallSiteWhereInfrastructure Uninstall SiteWhere Infrastructure components in the cluster
func UninstallSiteWhereInfrastructure(config *SiteWhereInstallConfiguration) error {
	var err error

	for _, tpl := range infraTemplate {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = UninstallResourceFromFile(infraResource, config.KubernetesConfig, config.StatikFS)
				if err != nil && !errors.IsNotFound(err) {
					return err
				}
			}
		}
	}
	if config.Verbose {
		color.Info.Println("\nSiteWhere Infrastructure: Uninstalled")
	}
	return nil
}
