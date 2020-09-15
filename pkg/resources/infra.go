/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package resources

import (
	"fmt"
	"net/http"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	errors "k8s.io/apimachinery/pkg/api/errors"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// InfraTemplateResource template for resources files
type InfraTemplateResource struct {
	FileTemplate string
	FileCount    int
	Enabled      bool
}

var infraTemplate = []InfraTemplateResource{
	InfraTemplateResource{
		FileTemplate: "/infra/01-mosquitto/infra-mosquitto-%02d.yaml",
		FileCount:    4,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/02-redis/infra-redis-%02d.yaml",
		FileCount:    9,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/03-zookeeper/infra-zookeeper-%02d.yaml",
		FileCount:    5,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/04-kafka/infra-kafka-%02d.yaml",
		FileCount:    1,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/05-postgresql/infra-postgresql-%02d.yaml",
		FileCount:    3,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/06-syncope/infra-syncope-%02d.yaml",
		FileCount:    8,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/07-warp10/infra-warp10-%02d.yaml",
		FileCount:    4,
		Enabled:      false,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/08-influxdb/infra-influxdb-%02d.yaml",
		FileCount:    4,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/09-nifi/infra-nifi-%02d.yaml",
		FileCount:    4,
		Enabled:      true,
	},
	InfraTemplateResource{
		FileTemplate: "/infra/10-zoo-entrance/infra-zoo-entrance-%02d.yaml",
		FileCount:    3,
		Enabled:      true,
	},
}

var infraTemplateMin = []InfraTemplateResource{
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
func InstallSiteWhereInfrastructure(minimal bool,
	waitReady bool,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	if minimal {
		err = installSiteWhereInfrastructureMinimal(waitReady, statikFS, clientset, apiextensionsClientset, config)
	} else {
		err = installSiteWhereInfrastructureDefault(waitReady, statikFS, clientset, apiextensionsClientset, config)
	}

	// if config.IsVerbose() {
	// 	color.Info.Println("SiteWhere Infrastructure: Installed")
	// }
	return err
}

func installSiteWhereInfrastructureDefault(waitReady bool,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	for _, tpl := range infraTemplate {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = InstallResourceFromFile(infraResource, statikFS, clientset, apiextensionsClientset, config)
				if err != nil && !errors.IsAlreadyExists(err) {
					return err
				}
			}
		}
	}

	if waitReady {
		err = waitForDeploymentAvailable(clientset, "sitewhere-infrastructure-mosquitto", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Deployment sitewhere-infrastructure-mosquitto: ")
		// 	color.Info.Println("Available")
		// }

		err = waitForDeploymentAvailable(clientset, "sitewhere-syncope", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Deployment sitewhere-syncope: ")
		// 	color.Info.Println("Available")
		// }

		for i := 0; i < 3; i++ {
			podName := fmt.Sprintf("sitewhere-kafka-zookeeper-%d", i)
			err = waitForPodContainersRunning(clientset, podName, sitewhereSystemNamespace)
			if err != nil {
				return err
			}
			// if config.IsVerbose() {
			// 	fmt.Print(fmt.Sprintf("Pod %s: ", podName))
			// 	color.Info.Println("Ready")
			// }
		}

		for i := 0; i < 3; i++ {
			podName := fmt.Sprintf("sitewhere-kafka-kafka-%d", i)
			err = waitForPodContainersRunning(clientset, podName, sitewhereSystemNamespace)
			if err != nil {
				return err
			}
			// if config.IsVerbose() {
			// 	fmt.Print(fmt.Sprintf("Pod %s: ", podName))
			// 	color.Info.Println("Ready")
			// }
		}

		err = waitForPodContainersRunning(clientset, "sitewhere-postgresql-0", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Pod sitewhere-postgresql-0: ")
		// 	color.Info.Println("Ready")
		// }

		for i := 0; i < 3; i++ {
			podName := fmt.Sprintf("sitewhere-infrastructure-redis-ha-server-%d", i)
			err = waitForPodContainersRunning(clientset, podName, sitewhereSystemNamespace)
			if err != nil {
				return err
			}
			// if config.IsVerbose() {
			// 	fmt.Print(fmt.Sprintf("Pod %s: ", podName))
			// 	color.Info.Println("Ready")
			// }
		}

		// err = waitForPodContainersRunning(clientset, "sitewhere-infrastructure-warp10-0", sitewhereSystemNamespace)
		// if err != nil {
		// 	return err
		// }
		// if config.IsVerbose() {
		// 	fmt.Print("Pod sitewhere-infrastructure-warp10-0: ")
		// 	color.Info.Println("Ready")
		// }
	}
	return err
}

func installSiteWhereInfrastructureMinimal(waitReady bool,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	for _, tpl := range infraTemplateMin {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = InstallResourceFromFile(infraResource, statikFS, clientset, apiextensionsClientset, config)
				if err != nil && !errors.IsAlreadyExists(err) {
					return err
				}
			}
		}
	}

	if waitReady {
		err = waitForDeploymentAvailable(clientset, "sitewhere-infrastructure-mosquitto", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Deployment sitewhere-infrastructure-mosquitto: ")
		// 	color.Info.Println("Available")
		// }

		err = waitForDeploymentAvailable(clientset, "sitewhere-kafka-entity-operator", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Deployment sitewhere-kafka-entity-operator: ")
		// 	color.Info.Println("Available")
		// }

		err = waitForDeploymentAvailable(clientset, "sitewhere-syncope", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Deployment sitewhere-syncope: ")
		// 	color.Info.Println("Available")
		// }

		err = waitForPodContainersRunning(clientset, "sitewhere-infrastructure-zookeeper-0", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Pod sitewhere-infrastructure-zookeeper-0: ")
		// 	color.Info.Println("Ready")
		// }
		// TODO if not minimal, wait for -1 and -2

		err = waitForPodContainersRunning(clientset, "sitewhere-kafka-kafka-0", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Pod sitewhere-kafka-kafka-0: ")
		// 	color.Info.Println("Ready")
		// }

		err = waitForPodContainersRunning(clientset, "sitewhere-postgresql-0", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Pod sitewhere-postgresql-0: ")
		// 	color.Info.Println("Ready")
		// }

		err = waitForPodContainersRunning(clientset, "sitewhere-infrastructure-redis-ha-server-0", sitewhereSystemNamespace)
		if err != nil {
			return err
		}
		// if config.IsVerbose() {
		// 	fmt.Print("Pod sitewhere-infrastructure-redis-ha-server-0: ")
		// 	color.Info.Println("Ready")
		// }
	}
	return nil
}

// UninstallSiteWhereInfrastructure Uninstall SiteWhere Infrastructure components in the cluster
func UninstallSiteWhereInfrastructure(minimal bool,
	statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	if minimal {
		err = uninstallSiteWhereInfrastructureMinimal(statikFS, clientset, apiextensionsClientset, config)
	} else {
		err = uninstallSiteWhereInfrastructureDefault(statikFS, clientset, apiextensionsClientset, config)
	}

	// if config.Verbose {
	// 	fmt.Print("SiteWhere Infrastructure: ")
	// 	color.Info.Println("Uninstalled")
	// }
	return err
}

func uninstallSiteWhereInfrastructureMinimal(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	for _, tpl := range infraTemplateMin {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = UninstallResourceFromFile(infraResource, statikFS, clientset, apiextensionsClientset, config)
				if err != nil && !errors.IsNotFound(err) {
					return err
				}
			}
		}
	}
	return nil
}

func uninstallSiteWhereInfrastructureDefault(statikFS http.FileSystem,
	clientset kubernetes.Interface,
	apiextensionsClientset apiextensionsclientset.Interface,
	config *rest.Config) error {
	var err error

	for _, tpl := range infraTemplate {
		if tpl.Enabled {
			for i := 1; i <= tpl.FileCount; i++ {
				var infraResource = fmt.Sprintf(tpl.FileTemplate, i)
				err = UninstallResourceFromFile(infraResource, statikFS, clientset, apiextensionsClientset, config)
				if err != nil && !errors.IsNotFound(err) {
					return err
				}
			}
		}
	}
	return nil
}
