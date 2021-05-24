package main

import (
	"capdash/capi"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main()  {
	dr,err := capi.DynamicResourcesClient("management-kubeconfig","cluster.x-k8s.io","MachineDeployment")
	if err != nil {
		fmt.Println(err.Error())
	}

	list, err := dr.List(metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(list)
}
