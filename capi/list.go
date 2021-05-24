package capi

import (
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type Providers struct {
	Provider []Provider `json:"providers"`
}

type Provider struct {
	Index int `json:"index"`
	Name string `json:"name"`
	Provider string `json:"provider"`
}

func RunListProvider() (*Providers,error) {
	c, err := client.New(CfgFile)
	if err != nil {
		return nil,err
	}

	configGet,err := c.GetProvidersConfig()
	if err != nil {
		return nil,err
	}

	providers := Providers{}

	for i,element := range configGet {
		providers.Provider = append(providers.Provider, Provider{
			Index:    i,
			Name:     element.Name(),
			Provider: string(element.Type()),
		})
	}

	return &providers,nil
}
