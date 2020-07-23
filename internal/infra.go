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
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/rest"
)

// InfraTemplateResource template for resources files
type InfraTemplateResource struct {
	FileTemplate string
	FileCount    int
	Enabled      bool
}

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
func InstallSiteWhereInfrastructure(config *rest.Config, statikFS http.FileSystem) error {
	var err error

	for _, tpl := range infraTemplate {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = InstallResourceFromFile(infraResource, config, statikFS)
				if err != nil && !errors.IsAlreadyExists(err) {
					return err
				}
			}
		}
	}

	return nil
}

// UninstallSiteWhereInfrastructure Uninstall SiteWhere Infrastructure components in the cluster
func UninstallSiteWhereInfrastructure(config *rest.Config, statikFS http.FileSystem) error {
	var err error

	for _, tpl := range infraTemplate {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = UninstallResourceFromFile(infraResource, config, statikFS)
				if err != nil && !errors.IsNotFound(err) {
					return err
				}
			}
		}
	}

	return nil
}
