package capi

import (
	"capdash/db"
	"fmt"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type initOptions struct {
	kubeconfig              string
	kubeconfigContext       string
	coreProvider            string
	bootstrapProviders      []string
	controlPlaneProviders   []string
	infrastructureProviders []string
	targetNamespace         string
	watchingNamespace       string
	listImages              bool
}

func runInit(kubeconfig string,provider []string) (string,error) {
	initOpts := initOptions{
		kubeconfig:              kubeconfig,
		kubeconfigContext:       "",
		coreProvider:            "",
		bootstrapProviders:      nil,
		controlPlaneProviders:   nil,
		infrastructureProviders: provider,
		targetNamespace:         "",
		watchingNamespace:       "",
		listImages:              false,
	}

	c, err := client.New(CfgFile)
	if err != nil {
		return "",err
	}

	options := client.InitOptions{
		Kubeconfig:              client.Kubeconfig{Path: initOpts.kubeconfig, Context: initOpts.kubeconfigContext},
		CoreProvider:            initOpts.coreProvider,
		BootstrapProviders:      initOpts.bootstrapProviders,
		ControlPlaneProviders:   initOpts.controlPlaneProviders,
		InfrastructureProviders: initOpts.infrastructureProviders,
		TargetNamespace:         initOpts.targetNamespace,
		WatchingNamespace:       initOpts.watchingNamespace,
		LogUsageInstructions:    true,
	}

	if initOpts.listImages {
		images, err := c.InitImages(options)
		if err != nil {
			return "",err
		}

		for _, i := range images {
			fmt.Println(i)
		}
		return "",nil
	}

	if _, err := c.Init(options); err != nil {
		return "",err
	}



	return "",nil
}

func RunInitWithParameters(database *db.Database, provider []string) error {
	err := WriteManagementKubeconfig(database)
	if err != nil {
		return err
	}

	_,err = runInit("management-kubeconfig",provider)
	if err != nil {
		return err
	}

	return nil
}
