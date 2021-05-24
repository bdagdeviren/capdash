package capi

import (
	"bytes"
	"capdash/db"
	"context"
	"encoding/json"
	"github.com/johandry/klient"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func GetClusters(kubeConfig string) ([]ClusterList, error) {
	dr,err := DynamicResourcesClient(kubeConfig,"cluster.x-k8s.io","cluster")
	if err != nil {
		return nil,err
	}

	list, err := dr.List(metav1.ListOptions{})
	if err != nil {
		return nil,err
	}

	marshList, err := json.Marshal(list)
	if err != nil {
		return nil,err
	}

	listLast := List{}
	err = json.Unmarshal(marshList,&listLast)
	if err != nil {
		return nil,err
	}

	var clusterList []ClusterList

	for _, s := range listLast.Items {
		clusterList = append(clusterList,ClusterList{
			Name:   s.Metadata.Name,
			Status: s.Status.Phase,
		})
	}

	//json,err := json.Marshal(clusterList)
	//if err != nil {
	//	return "",err
	//}

	return clusterList,nil
}

func CreateCluster(database *db.Database, name string, infrastructureProvider string, namespace string, kubernetesVersion string, cpmc int64, vmc int64, templateProvider string) error {
	err := WriteManagementKubeconfig(database)
	if err != nil {
		return err
	}

	out,err := database.Client.Get(context.Background(),name).Result()
	if err!=nil {
		if err.Error() == "redis: nil" {
			cc := configClusterOptions{
				kubeconfig:               "management-kubeconfig",
				kubeconfigContext:        "",
				flavor:                   "",
				infrastructureProvider:   infrastructureProvider,
				targetNamespace:          namespace,
				kubernetesVersion:        kubernetesVersion,
				controlPlaneMachineCount: cpmc,
				workerMachineCount:       vmc,
				url:                      templateProvider,
				configMapNamespace:       "",
				configMapName:            "",
				configMapDataKey:         "",
				listVariables:            false,
			}

			out, err := runGetClusterTemplate(name, cc)
			if err != nil {
				return err
			}

			_, err = database.SetData(name, []byte(out))
			if err != nil {
				return err
			}
		}
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(out)),100)
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil{
			break
		}

		if rawObj.Raw != nil {
			c := klient.New("", "management-kubeconfig")
			if err := c.Apply(rawObj.Raw); err != nil {
				return err
			}
		}
	}

	return	nil
}

func DeleteCluster(database *db.Database, name string) error {
	err := WriteManagementKubeconfig(database)
	if err != nil {
		return err
	}

	out,err := database.Client.Get(context.Background(),name).Result()
	if err!=nil {
		return errors.New("Cannot find cluster!")
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(out)),100)
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil{
			break
		}

		if rawObj.Raw != nil {
			c := klient.New("", "management-kubeconfig")
			if err := c.Delete(rawObj.Raw); err != nil {
				return err
			}
		}
	}

	return	nil
}