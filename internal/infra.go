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
		fmt.Printf("Deploymene sitewhere-infrastructure-mosquitto: Available\n")
	}

	err = waitForDeploymentAvailable(config.GetClientset(), "sitewhere-kafka-entity-operator", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Deploymene sitewhere-kafka-entity-operator: Available\n")
	}

	err = waitForDeploymentAvailable(config.GetClientset(), "sitewhere-syncope", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Deploymene sitewhere-syncope: Available\n")
	}

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-infrastructure-zookeeper-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Pod sitewhere-infrastructure-zookeeper-0: Ready\n")
	}
	// TODO if not minimal, wait for -1 and -2

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-kafka-kafka-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Pod sitewhere-kafka-kafka-0: Ready\n")
	}

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-postgresql-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Pod sitewhere-postgresql-0: Ready\n")
	}

	err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-infrastructure-redis-ha-server-0", sitewhereSystemNamespace)
	if err != nil {
		return err
	}
	if config.IsVerbose() {
		fmt.Printf("Pod sitewhere-infrastructure-redis-ha-server-0: Ready\n")
	}

	// err = waitForPodContainersRunning(config.GetClientset(), "sitewhere-infrastructure-warp10-0", sitewhereSystemNamespace)
	// if err != nil {
	// 	return err
	// }
	// if config.IsVerbose() {
	// 	fmt.Printf("Pod sitewhere-infrastructure-warp10-0: Ready\n")
	// }

	if config.IsVerbose() {
		fmt.Printf("SiteWhere Infrastructure: Installed\n")
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
		fmt.Printf("SiteWhere Infrastructure: Uninstalled\n")
	}
	return nil
}
