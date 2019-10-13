package deployment

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	conf "github.com/kubemq-io/kubemqctl/pkg/k8s/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type StatefulSetDeployment struct {
	Client      *client.Client
	Namespace   *apiv1.Namespace
	StatefulSet *appsv1.StatefulSet
	Services    []*apiv1.Service
	ConfigMaps  []*apiv1.ConfigMap
	Secrets     []*apiv1.Secret
}

func NewStatefulSetDeployment(cfg *config.Config) (*StatefulSetDeployment, error) {
	sd := &StatefulSetDeployment{
		Client:      nil,
		Namespace:   nil,
		StatefulSet: nil,
		Services:    nil,
		ConfigMaps:  nil,
		Secrets:     nil,
	}
	c, err := client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	sd.Client = c

	return sd, nil
}

func NewStatefulSetDeploymentFromCluster(cfg *config.Config, ns, name string) (*StatefulSetDeployment, error) {
	sd := &StatefulSetDeployment{
		Client:      nil,
		Namespace:   nil,
		StatefulSet: nil,
		Services:    nil,
		ConfigMaps:  nil,
		Secrets:     nil,
	}
	c, err := client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	sd.Client = c

	sd.StatefulSet, err = c.GetStatefulSet(ns, name)
	if err != nil {
		return nil, err
	}
	if sd.StatefulSet == nil {
		return sd, nil
	}
	sd.Services, err = c.GetServices(ns, sd.StatefulSet.Spec.Template.ObjectMeta.Labels)
	if err != nil {
		return nil, err
	}

	sd.ConfigMaps, err = c.GetConfigMaps(ns, sd.StatefulSet.Spec.Template.Spec.Volumes)
	if err != nil {
		return nil, err
	}
	sd.Secrets, err = c.GetSecrets(ns, sd.StatefulSet.Spec.Template.Spec.Volumes)
	if err != nil {
		return nil, err
	}
	return sd, nil
}

func (sd *StatefulSetDeployment) CreateStatefulSetDeployment(o *Options, optionsMenu *conf.Menu) error {

	ns, existed, err := sd.Client.GetNamespace(o.Namespace)
	if err != nil {
		return err
	}
	if !existed {
		sd.Namespace = ns
	}
	stSpec, err := NewStatefulSetConfig(o).Spec()
	if err != nil {
		return err
	}
	sts := &appsv1.StatefulSet{}
	err = yaml.Unmarshal(stSpec, sts)
	if err != nil {
		return err
	}
	sd.StatefulSet = sts
	if len(sd.StatefulSet.Spec.Template.Spec.Containers) > 0 {
		envVars := optionsMenu.ExportEnvVar()
		sd.StatefulSet.Spec.Template.Spec.Containers[0].Env = append(sd.StatefulSet.Spec.Template.Spec.Containers[0].Env, envVars...)

		volMounts := optionsMenu.ExportVolumeMounts()
		sd.StatefulSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(sd.StatefulSet.Spec.Template.Spec.Containers[0].VolumeMounts, volMounts...)
	}
	vols := optionsMenu.ExportVolumes()
	sd.StatefulSet.Spec.Template.Spec.Volumes = append(sd.StatefulSet.Spec.Template.Spec.Volumes, vols...)

	configMaps := optionsMenu.ExportConfigMaps()
	for _, cm := range configMaps {
		sd.ConfigMaps = append(sd.ConfigMaps, &apiv1.ConfigMap{
			TypeMeta: v1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      cm.Name,
				Namespace: o.Namespace,
			},
			Data:       map[string]string{cm.FileName: cm.Value},
			BinaryData: nil,
		})
	}

	secrets := optionsMenu.ExportSecrets()

	for _, sec := range secrets {
		b64Value := []byte(b64.StdEncoding.EncodeToString([]byte(sec.Value)))
		sd.Secrets = append(sd.Secrets, &apiv1.Secret{
			TypeMeta: v1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      sec.Name,
				Namespace: o.Namespace,
			},
			Data: map[string][]byte{sec.FileName: b64Value},
			Type: apiv1.SecretType(sec.SecretType),
		})
	}

	for _, svcfg := range NewServiceConfigs(o) {
		spec, err := svcfg.Spec()
		if err != nil {
			return err
		}
		svc := &apiv1.Service{}
		err = yaml.Unmarshal(spec, svc)
		if err != nil {
			return err
		}
		sd.Services = append(sd.Services, svc)
	}
	return nil
}

func (sd *StatefulSetDeployment) Execute(name, namespace string) (bool, error) {
	if sd.Namespace != nil {
		_, err := sd.Client.ClientSet.CoreV1().Namespaces().Create(sd.Namespace)
		if err != nil {
			return false, err
		}
		utils.Printlnf("Namespace %s created", sd.Namespace)
	}
	var err error

	newSts, isUpdated, err := sd.Client.CreateOrUpdateStatefulSet(sd.StatefulSet)

	if err != nil {
		utils.Printlnf("StatefulSet %s/%s not created. Error: %s", namespace, name, utils.Title(err.Error()))
	} else {
		if newSts != nil && isUpdated {
			utils.Printlnf("StatefulSet %s/%s configured", namespace, name)
		} else if newSts != nil {
			utils.Printlnf("StatefulSet %s/%s created", namespace, name)
		}
	}

	for _, svc := range sd.Services {
		newSvc, isUpdated, err := sd.Client.CreateOrUpdateService(svc)
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
		newCm, isUpdated, err := sd.Client.CreateOrUpdateConfigMap(cm)
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
		newCm, isUpdated, err := sd.Client.CreateOrUpdateSecret(sec)
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

func writeSection(w *bufio.Writer, data []byte) error {
	var err error
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.WriteString(fmt.Sprintf("---\n"))
	if err != nil {
		return err
	}
	return nil
}
func (sd *StatefulSetDeployment) Export(out io.Writer) error {

	w := bufio.NewWriter(out)
	if sd.Namespace != nil {
		data, err := yaml.Marshal(sd.Namespace)
		if err != nil {
			return err
		}
		if err := writeSection(w, data); err != nil {
			return err
		}
	}
	if sd.StatefulSet != nil {
		data, err := yaml.Marshal(sd.StatefulSet)
		if err != nil {
			return err
		}
		if err := writeSection(w, data); err != nil {
			return err
		}
	}

	for _, svc := range sd.Services {
		data, err := yaml.Marshal(svc)
		if err != nil {
			return err
		}
		if err := writeSection(w, data); err != nil {
			return err
		}
	}
	for _, cm := range sd.ConfigMaps {
		data, err := yaml.Marshal(cm)
		if err != nil {
			return err
		}
		if err := writeSection(w, data); err != nil {
			return err
		}
	}
	for _, sec := range sd.Secrets {
		data, err := yaml.Marshal(sec)
		if err != nil {
			return err
		}
		if err := writeSection(w, data); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func (sd *StatefulSetDeployment) Import(input string) error {
	segments := strings.Split(input, "---")
	if len(segments) == 0 {
		segments = append(segments, input)
	}
	for index, seg := range segments {
		if strings.Contains(seg, "kind: StatefulSet") {
			sts := &appsv1.StatefulSet{}
			err := yaml.Unmarshal([]byte(seg), sts)
			if err != nil {
				return fmt.Errorf("error parsing StatefulSet yaml (segment %d), %s", index, err)
			}
			sd.StatefulSet = sts
			continue
		}
		if strings.Contains(seg, "kind: Service") {
			svc := &apiv1.Service{}
			err := yaml.Unmarshal([]byte(seg), svc)
			if err != nil {
				return fmt.Errorf("error parsing Service yaml (segment %d), %s", index, err)
			}
			sd.Services = append(sd.Services, svc)
			continue
		}
		if strings.Contains(seg, "kind: ConfigMap") {
			cm := &apiv1.ConfigMap{}
			err := yaml.Unmarshal([]byte(seg), cm)
			if err != nil {
				return fmt.Errorf("error parsing ConfigMap yaml (segment %d), %s", index, err)
			}
			sd.ConfigMaps = append(sd.ConfigMaps, cm)
			continue
		}
		if strings.Contains(seg, "kind: Secret") {
			sec := &apiv1.Secret{}
			err := yaml.Unmarshal([]byte(seg), sec)
			if err != nil {
				return fmt.Errorf("error parsing Secret yaml (segment %d), %s", index, err)
			}
			sd.Secrets = append(sd.Secrets, sec)
			continue
		}
	}

	return sd.Validate()
}

func (sd *StatefulSetDeployment) Validate() error {
	if sd.StatefulSet == nil {
		return fmt.Errorf("invalid import file, no StatefulSet was defined")
	}
	if len(sd.StatefulSet.Spec.Template.Spec.Containers) > 0 {
		if !strings.Contains(sd.StatefulSet.Spec.Template.Spec.Containers[0].Image, "kubemq") {
			return fmt.Errorf("invalid KubeMQ StatefulSet definitions, docker image is invalid")
		}
	}
	if len(sd.Services) == 0 {
		return fmt.Errorf("invalid KubeMQ Services definitions, at least one serivce must be defined")

	}
	return nil
}
