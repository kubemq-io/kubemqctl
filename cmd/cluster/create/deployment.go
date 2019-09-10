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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatefulSetDeployment struct {
	client      *client.Client
	Namespace   *apiv1.Namespace
	StatefulSet *appsv1.StatefulSet
	Services    []*apiv1.Service
	ConfigMaps  []*apiv1.ConfigMap
	Secrets     []*apiv1.Secret
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
	envVars := o.optionsMenu.ExportEnvVar()
	sd.StatefulSet.Spec.Template.Spec.Containers[0].Env = append(sd.StatefulSet.Spec.Template.Spec.Containers[0].Env, envVars...)

	volMounts := o.optionsMenu.ExportVolumeMounts()
	sd.StatefulSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(sd.StatefulSet.Spec.Template.Spec.Containers[0].VolumeMounts, volMounts...)

	vols := o.optionsMenu.ExportVolumes()
	sd.StatefulSet.Spec.Template.Spec.Volumes = append(sd.StatefulSet.Spec.Template.Spec.Volumes, vols...)

	configMaps := o.optionsMenu.ExportConfigMaps()
	for _, cm := range configMaps {
		sd.ConfigMaps = append(sd.ConfigMaps, &apiv1.ConfigMap{
			TypeMeta: v1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      cm.Name,
				Namespace: o.namespace,
			},
			Data:       map[string]string{cm.FileName: cm.Value},
			BinaryData: nil,
		})
	}

	secrets := o.optionsMenu.ExportSecrets()

	for _, sec := range secrets {
		sd.Secrets = append(sd.Secrets, &apiv1.Secret{
			TypeMeta: v1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      sec.Name,
				Namespace: o.namespace,
			},
			Data:       nil,
			StringData: map[string]string{sec.FileName: sec.Value},
			Type:       apiv1.SecretType(sec.SecretType),
		})
	}
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

	newSts, isUpdated, err := sd.client.CreateOrUpdateStatefulSet(sd.StatefulSet)

	if err != nil {
		utils.Printlnf("StatefulSet %s/%s not created. Error: %s", o.namespace, o.name, utils.Title(err.Error()))
	} else {
		if newSts != nil && isUpdated {
			utils.Printlnf("StatefulSet %s/%s configured", o.namespace, o.name)
		} else if newSts != nil {
			utils.Printlnf("StatefulSet %s/%s created", o.namespace, o.name)
		}
	}

	for _, svc := range sd.Services {
		newSvc, isUpdated, err := sd.client.CreateOrUpdateService(svc)
		if err != nil {
			utils.Printlnf("Service %s/%s not created. Error: %s", svc.Namespace, svc.Name, utils.Title(err.Error()))
		} else {
			if newSvc != nil && isUpdated {
				utils.Printlnf("Service %s/%s configured", svc.Namespace, svc.Name)
			} else if newSvc != nil {
				utils.Printlnf("Service %s/%s created", svc.Namespace, svc.Name)
			}

		}
	}

	for _, cm := range sd.ConfigMaps {
		newCm, isUpdated, err := sd.client.CreateOrUpdateConfigMap(cm)
		if err != nil {
			utils.Printlnf("ConfigMap %s/%s not created. Error: %s", cm.Namespace, cm.Name, utils.Title(err.Error()))
		} else {
			if newCm != nil && isUpdated {
				utils.Printlnf("ConfigMap %s/%s configured", cm.Namespace, cm.Name)
			} else if newCm != nil {
				utils.Printlnf("ConfigMap %s/%s created", cm.Namespace, cm.Name)
			}
		}
	}
	for _, sec := range sd.Secrets {
		newCm, isUpdated, err := sd.client.CreateOrUpdateSecret(sec)
		if err != nil {
			utils.Printlnf("Secret %s/%s not created. Error: %s", sec.Namespace, sec.Name, utils.Title(err.Error()))
		} else {
			if newCm != nil && isUpdated {
				utils.Printlnf("Secret %s/%s configured", sec.Namespace, sec.Name)
			} else if newCm != nil {
				utils.Printlnf("Secret %s/%s created", sec.Namespace, sec.Name)
			}
		}
	}
	return true, nil
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
	for _, cm := range sd.ConfigMaps {
		data, err := yaml.Marshal(cm)
		if err != nil {
			return err
		}
		w.Write(data)
		w.WriteString(fmt.Sprintf("---\n"))
	}
	for _, sec := range sd.Secrets {
		data, err := yaml.Marshal(sec)
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
