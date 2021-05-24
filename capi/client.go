package capi

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func DynamicResourcesClient(kubeConfig,group,kind string) (dynamic.NamespaceableResourceInterface,error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil,err
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil,err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil,err
	}

	gvk := &schema.GroupVersionKind{Group: group,Version: "",Kind: kind} //

	mapping, err := mapper.RESTMapping(gvk.GroupKind(),gvk.Version)
	if err != nil {
		return nil,err
	}

	dr := dyn.Resource(mapping.Resource)

	return dr,nil
}
