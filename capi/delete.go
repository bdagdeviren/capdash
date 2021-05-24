package capi

import (
	"capdash/db"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type deleteOptions struct {
	kubeconfig              string
	kubeconfigContext       string
	targetNamespace         string
	coreProvider            string
	bootstrapProviders      []string
	controlPlaneProviders   []string
	infrastructureProviders []string
	includeNamespace        bool
	includeCRDs             bool
	deleteAll               bool
}

func init() {
	CfgFile = "./cluster.yaml"
}

func runDelete(kubeconfig string,provider []string) (string,error) {
	
	dd := deleteOptions{
		kubeconfig:              kubeconfig,
		kubeconfigContext:       "",
		targetNamespace:         "",
		coreProvider:            "",
		bootstrapProviders:      nil,
		controlPlaneProviders:   nil,
		infrastructureProviders: provider,
		includeNamespace:        false,
		includeCRDs:             false,
		deleteAll:               false,
	}
	c, err := client.New(CfgFile)
	if err != nil {
		return "",err
	}

	hasProviderNames := (dd.coreProvider != "") ||
		(len(dd.bootstrapProviders) > 0) ||
		(len(dd.controlPlaneProviders) > 0) ||
		(len(dd.infrastructureProviders) > 0)

	if dd.deleteAll && hasProviderNames {
		return "",errors.New("The --all flag can't be used in combination with --core, --bootstrap, --control-plane, --infrastructure")
	}

	if !dd.deleteAll && !hasProviderNames {
		return "",errors.New("At least one of --core, --bootstrap, --control-plane, --infrastructure should be specified or the --all flag should be set")
	}

	if err := c.Delete(client.DeleteOptions{
		Kubeconfig:              client.Kubeconfig{Path: dd.kubeconfig, Context: dd.kubeconfigContext},
		IncludeNamespace:        dd.includeNamespace,
		IncludeCRDs:             dd.includeCRDs,
		Namespace:               dd.targetNamespace,
		CoreProvider:            dd.coreProvider,
		BootstrapProviders:      dd.bootstrapProviders,
		InfrastructureProviders: dd.infrastructureProviders,
		ControlPlaneProviders:   dd.controlPlaneProviders,
		DeleteAll:               dd.deleteAll,
	}); err != nil {
		return "",err
	}

	return "",nil
}

func RunDeleteWithParameters(database *db.Database, provider []string) error {
	err := WriteManagementKubeconfig(database)
	if err != nil {
		return err
	}
	_, err = runDelete("management-kubeconfig", provider)
	if err != nil {
		return err
	}

	return nil
}
