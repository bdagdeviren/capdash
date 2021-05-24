package capi

var (
	CfgFile = "./hack/cluster.yaml"
)

type generateYAMLOptions struct {
	url           string
	listVariables bool
}

type configClusterOptions struct {
	kubeconfig             string
	kubeconfigContext      string
	flavor                 string
	infrastructureProvider string

	targetNamespace          string
	kubernetesVersion        string
	controlPlaneMachineCount int64
	workerMachineCount       int64

	url                string
	configMapNamespace string
	configMapName      string
	configMapDataKey   string

	listVariables bool
}

type List struct {
	Items []Items `json:"items"`
}

type Items struct {
	Metadata Metadata `json:"metadata"`
	Status  Status  `json:"status"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Status struct {
	Phase  string  `json:"phase"`
}

type ClusterList struct {
	Name string `json:"name"`
	Status string `json:"status"`
}

