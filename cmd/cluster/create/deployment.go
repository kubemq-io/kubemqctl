package create

import (
	"bufio"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

type StatefulSetDeployment struct {
	client      *client.Client
	Namespace   *apiv1.Namespace
	StatefulSet *appsv1.StatefulSet
	Services    []*apiv1.Service
}

func CreateStatefulSetDeployment(o *CreateOptions) (*StatefulSetDeployment, error) {
	sd := &StatefulSetDeployment{
		client:      nil,
		Namespace:   nil,
		StatefulSet: nil,
		Services:    nil,
	}
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	sd.client = c
	ns, existed, err := c.GetNamespace(o.namespace)
	if err != nil {
		return nil, err
	}
	if !existed {
		sd.Namespace = ns
	}
	stSpec, err := NewStatefulSetConfig(o).Spec()
	if err != nil {
		return nil, err
	}
	sts := &appsv1.StatefulSet{}
	err = yaml.Unmarshal(stSpec, sts)
	if err != nil {
		return nil, err
	}
	sd.StatefulSet = sts
	envVars := o.envVars.ExportEnvVar()
	sd.StatefulSet.Spec.Template.Spec.Containers[0].Env = append(sd.StatefulSet.Spec.Template.Spec.Containers[0].Env, envVars...)

	for _, svcfg := range NewServiceConfigs(o) {
		spec, err := svcfg.Spec()
		if err != nil {
			return nil, err
		}
		svc := &apiv1.Service{}
		err = yaml.Unmarshal(spec, svc)
		sd.Services = append(sd.Services, svc)
	}
	return sd, nil
}
func (sd *StatefulSetDeployment) Execute(o *CreateOptions) (bool, error) {
	if sd.Namespace != nil {
		_, err := sd.client.ClientSet.CoreV1().Namespaces().Create(sd.Namespace)
		if err != nil {
			return false, err
		}
		utils.Printlnf("Namespace %s created", o.namespace)
	}
	var err error
	_, err = sd.client.ClientSet.AppsV1().StatefulSets(sd.StatefulSet.Namespace).Create(sd.StatefulSet)
	isCreate := false
	if err != nil {
		utils.Printlnf("StatefulSet %s/%s not created. Error: %s", o.namespace, o.name, utils.Title(err.Error()))
	} else {
		utils.Printlnf("StatefulSet %s/%s created", o.namespace, o.name)
		isCreate = true
	}

	for _, svc := range sd.Services {
		_, err := sd.client.ClientSet.CoreV1().Services(svc.Namespace).Create(svc)
		if err != nil {
			utils.Printlnf("Service %s/%s not created. Error: %s", svc.Namespace, svc.Name, utils.Title(err.Error()))
		} else {
			utils.Printlnf("Service %s/%s created", svc.Namespace, svc.Name)
		}
	}
	return isCreate, nil
}

func (sd *StatefulSetDeployment) Export(out io.Writer, o *CreateOptions) error {
	w := bufio.NewWriter(out)
	if sd.Namespace != nil {
		data, err := yaml.Marshal(sd.Namespace)
		if err != nil {
			return err
		}
		w.Write(data)
		w.WriteString(fmt.Sprintf("---\n"))
	}
	if sd.StatefulSet != nil {
		data, err := yaml.Marshal(sd.StatefulSet)
		if err != nil {
			return err
		}
		w.Write(data)
		w.WriteString(fmt.Sprintf("---\n"))
	}

	for _, svc := range sd.Services {
		data, err := yaml.Marshal(svc)
		if err != nil {
			return err
		}
		w.Write(data)
		w.WriteString(fmt.Sprintf("---\n"))
	}
	var err error
	err = w.Flush()
	if err != nil {
		return err
	}
	utils.Printlnf("create yaml file was exported to %s.yaml", o.name)
	return nil
}
