package capi

import (
	"capdash/db"
	"io/ioutil"
	"os"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type getKubeconfigOptions struct {
	kubeconfig        string
	kubeconfigContext string
	namespace         string
}

func init() {
	CfgFile = ""
}

func runGetKubeconfig(workloadClusterName string,gk getKubeconfigOptions) (string,error) {
	c, err := client.New(CfgFile)
	if err != nil {
		return "",err
	}

	options := client.GetKubeconfigOptions{
		Kubeconfig:          client.Kubeconfig{Path: gk.kubeconfig, Context: gk.kubeconfigContext},
		WorkloadClusterName: workloadClusterName,
		Namespace:           gk.namespace,
	}

	out, err := c.GetKubeconfig(options)
	if err != nil {
		return "",err
	}

	return out,nil
}

func RunGetWorkloadClusterKubeconfigWithParameters(database *db.Database,clusterName string) (string,error) {
	err := WriteManagementKubeconfig(database)
	if err != nil {
		return "",err
	}

	gk := getKubeconfigOptions{
		kubeconfig:        "management-kubeconfig",
		kubeconfigContext: "",
		namespace:         "",
	}

	out,err := runGetKubeconfig(clusterName,gk)

	if err != nil {
		return "",err
	}

	return out,nil
}

func WriteManagementKubeconfig(database *db.Database) error {
	_, err := os.Stat("management-kubeconfig")
	if os.IsNotExist(err) {
		kubeConfig, err := database.GetManagementKubeconfig("management-kubeconfig")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile("management-kubeconfig", []byte(kubeConfig), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
