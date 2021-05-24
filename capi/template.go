package capi

import (
	"capdash/db"
	"fmt"
	"io"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

var gyOpts = &generateYAMLOptions{}

func generateYAML(r io.Reader, w io.Writer) error {
	c, err := client.New(CfgFile)
	if err != nil {
		return err
	}
	options := client.ProcessYAMLOptions{
		ListVariablesOnly: gyOpts.listVariables,
	}
	if gyOpts.url != "" {
		if gyOpts.url == "-" {
			options.ReaderSource = &client.ReaderSourceOptions{
				Reader: r,
			}
		} else {
			options.URLSource = &client.URLSourceOptions{
				URL: gyOpts.url,
			}
		}
	}
	printer, err := c.ProcessYAML(options)
	if err != nil {
		return err
	}
	if gyOpts.listVariables {
		if len(printer.Variables()) > 0 {
			fmt.Fprintln(w, "Variables:")
			for _, v := range printer.Variables() {
				fmt.Fprintf(w, "  - %s\n", v)
			}
		} else {
			fmt.Fprintln(w)
		}
		return nil
	}
	out, err := printer.Yaml()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(out))
	return err
}

func runGetClusterTemplate(name string, cc configClusterOptions) (string,error) {
	c, err := client.New(CfgFile)
	if err != nil {
		return "",err
	}

	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig:        client.Kubeconfig{Path: cc.kubeconfig, Context: cc.kubeconfigContext},
		ClusterName:       name,
		TargetNamespace:   cc.targetNamespace,
		KubernetesVersion: cc.kubernetesVersion,
		ListVariablesOnly: cc.listVariables,
	}

	if cc.url != "" {
		templateOptions.URLSource = &client.URLSourceOptions{
			URL: cc.url,
		}
	}

	if cc.configMapNamespace != "" || cc.configMapName != "" || cc.configMapDataKey != "" {
		templateOptions.ConfigMapSource = &client.ConfigMapSourceOptions{
			Namespace: cc.configMapNamespace,
			Name:      cc.configMapName,
			DataKey:   cc.configMapDataKey,
		}
	}

	if cc.infrastructureProvider != "" || cc.flavor != "" {
		templateOptions.ProviderRepositorySource = &client.ProviderRepositorySourceOptions{
			InfrastructureProvider: cc.infrastructureProvider,
			Flavor:                 cc.flavor,
		}
	}

	template, err := c.GetClusterTemplate(templateOptions)
	if err != nil {
		return "",err
	}

	if cc.listVariables {
		return "",templateListVariablesOutput(template)
	}

	return templateYAMLOutput(template)
}

func RunGetClusterTemplate(database *db.Database, name string, infrastructureProvider string, namespace string, kubernetesVersion string, cpmc int64, vmc int64, templateProvider string) (string, error) {
	err := WriteManagementKubeconfig(database)
	if err != nil {
		return "",err
	}

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
		return "",err
	}

	_, err = database.SetDataWithBase64(name, []byte(out))
	if err != nil {
		return "", err
	}

	fmt.Println(out)

	return out,nil
}

func templateListVariablesOutput(template client.Template) error {
	if len(template.Variables()) > 0 {
		fmt.Println("Variables:")
		for _, v := range template.Variables() {
			fmt.Printf("  - %s\n", v)
		}
	}
	fmt.Println()
	return nil
}

func templateYAMLOutput(template client.Template) (string,error) {
	yaml, err := template.Yaml()
	if err != nil {
		return "",err
	}
	yaml = append(yaml, '\n')

	//if _, err := os.Stdout.Write(yaml); err != nil {
	//	return errors.Wrap(err, "failed to write yaml to Stdout")
	//}
	return string(yaml),nil
}
