package openshift

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/api/errors"
	kclientcmd "k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"

	projectv1 "github.com/openshift/api/project/v1"
	projectclient "github.com/openshift/client-go/project/clientset/versioned"
	"github.com/openshift/origin/pkg/oc/lib/kubeconfig"
)

const requestProjectSwitchProjectOutput = `Project %[2]q created on server %[3]q.

To switch to this project and start adding applications, use:

    %[1]s project %[2]s
`

// CreateProject creates a new project for the user
func CreateProject(name, display, desc string, out io.Writer) error {
	userFactory, err := LoggedInUserFactory()
	if err != nil {
		return err
	}
	clientConfig, err := userFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	projectClient, err := projectclient.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	projectRequest := &projectv1.ProjectRequest{}
	projectRequest.Name = name
	projectRequest.DisplayName = display
	projectRequest.Description = desc

	newProject, err := projectClient.ProjectV1().ProjectRequests().Create(projectRequest)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return setCurrentProject(userFactory, name, out)
		}
		return err
	}

	fmt.Fprintf(out, requestProjectSwitchProjectOutput, "oc", newProject.Name, clientConfig.Host)
	return nil
}

func setCurrentProject(f genericclioptions.RESTClientGetter, name string, out io.Writer) error {
	config, err := f.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}
	config.CurrentContext = name

	// TODO: Make the kubeconfig shared helper
	pathOptions := kubeconfig.NewPathOptionsWithConfig("")
	return kclientcmd.ModifyConfig(pathOptions, config, true)
}

func LoggedInUserFactory() (genericclioptions.RESTClientGetter, error) {
	cfg, err := kclientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, err
	}
	defaultCfg := kclientcmd.NewDefaultClientConfig(*cfg, &kclientcmd.ConfigOverrides{})

	return genericclioptions.NewTestConfigFlags().WithClientConfig(defaultCfg), nil
}
