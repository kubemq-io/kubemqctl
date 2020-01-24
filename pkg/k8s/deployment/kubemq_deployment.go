package deployment

import (
	"bufio"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"io"
	"strings"
)

type KubeMQDeployment struct {
	Client *client.Client
	*KubeMQManifestConfig
}

func NewKubeMQDeployment(cfg *config.Config, manifestConfig *KubeMQManifestConfig) (*KubeMQDeployment, error) {
	sd := &KubeMQDeployment{
		Client:               nil,
		KubeMQManifestConfig: manifestConfig,
	}
	c, err := client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	sd.Client = c

	return sd, nil
}

func NewKubeMQDeploymentFromCluster(client *client.Client, manifestConfig *KubeMQManifestConfig) (*KubeMQDeployment, error) {
	sd := &KubeMQDeployment{
		Client:               client,
		KubeMQManifestConfig: manifestConfig,
	}

	sts, err := sd.Client.GetStatefulSet(sd.StatefulSet.Namespace, sd.Name)
	if err != nil {
		return nil, err
	}
	if sts == nil {
		return sd, nil
	}
	sd.StatefulSet.Set(sts)

	svcList, err := sd.Client.GetServices(sd.Namespace, sts.Spec.Template.ObjectMeta.Labels)
	if err != nil {
		return nil, err
	}

	for _, svc := range svcList {
		svcConfig := NewServiceConfig("", svc.Name, svc.Namespace, sts.Name)
		svcConfig.Set(svc)
		sd.Services[svc.Name] = svcConfig
	}

	cmList, err := sd.Client.GetConfigMaps(sts.Namespace, sts.Name)
	if err != nil {
		return nil, err
	}

	for _, cm := range cmList {
		cmConfig := NewConfigMap("", cm.Name, cm.Namespace)
		cmConfig.Set(cm)
		sd.ConfigMaps[cm.Name] = cmConfig
	}

	secList, err := sd.Client.GetSecrets(sts.Namespace, sts.Name)
	if err != nil {
		return nil, err
	}
	for _, sec := range secList {
		secConfig := NewSecretConfig("", sec.Name, sec.Namespace)
		secConfig.Set(sec)
		sd.Secrets[sec.Name] = secConfig
	}

	ingressList, err := sd.Client.GetIngress(sts.Namespace, sts.Name)
	if err != nil {
		return nil, err
	}
	for _, ingress := range ingressList {
		ingressConfig := NewIngressConfig("", ingress.Name, ingress.Namespace, sts.Name)
		ingressConfig.Set(ingress)
		sd.Ingress[ingress.Name] = ingressConfig
	}
	return sd, nil
}

func (sd *KubeMQDeployment) Execute(name, namespace string) (bool, error) {
	if sd.NamespaceConfig != nil {
		ns, err := sd.NamespaceConfig.Get()
		if err != nil {

			return false, err
		}
		createdNamespace, created, err := sd.Client.CheckAndCreateNamespace(ns)
		if err != nil {
			return false, err
		}
		if created {
			utils.Printlnf("Namespace %s created", sd.Namespace)
		}
		sd.NamespaceConfig.Set(createdNamespace)

	}

	var err error
	requestedSts, err := sd.StatefulSet.Get()
	if err != nil {
		return false, err
	}
	newSts, isUpdated, err := sd.Client.CreateOrUpdateStatefulSet(requestedSts)
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
		requestedSvc, err := svc.Get()
		if err != nil {
			return false, err
		}
		newSvc, isUpdated, err := sd.Client.CreateOrUpdateService(requestedSvc)
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
		requestedConfigMap, err := cm.Get()
		if err != nil {
			return false, err
		}
		newCm, isUpdated, err := sd.Client.CreateOrUpdateConfigMap(requestedConfigMap)
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
		requestedSec, err := sec.Get()
		if err != nil {
			return false, err
		}
		newCm, isUpdated, err := sd.Client.CreateOrUpdateSecret(requestedSec)
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
	for _, ing := range sd.Ingress {
		requestedIng, err := ing.Get()
		if err != nil {
			return false, err
		}
		newIngress, isUpdated, err := sd.Client.CreateOrUpdateIngress(requestedIng)
		if err != nil {
			utils.Printlnf("Ingress %s/%s not created. Error: %s", ing.Namespace, ing.Name, utils.Title(err.Error()))
		} else {
			if newIngress != nil && isUpdated {
				utils.Printlnf("Ingress %s/%s configured", ing.Namespace, ing.Name)
			} else if newIngress != nil {
				utils.Printlnf("Ingress %s/%s created", ing.Namespace, ing.Name)
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
func (sd *KubeMQDeployment) Export(out io.Writer) error {
	w := bufio.NewWriter(out)
	data, err := sd.Spec()
	if err != nil {
		return err
	}
	if err := writeSection(w, data); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func (sd *KubeMQDeployment) Import(input string) error {
	segments := strings.Split(input, "---")
	if len(segments) == 0 {
		segments = append(segments, input)
	}
	for index, seg := range segments {
		if strings.Contains(seg, "#") {
			continue
		}
		if strings.Contains(seg, "kind: Namespace") {
			ns, err := ImportNamespaceConfig([]byte(seg))
			if err != nil {
				return fmt.Errorf("error parsing namespace yaml (segment %d), %s", index, err)
			}
			sd.NamespaceConfig = ns
			continue
		}
		if strings.Contains(seg, "kind: StatefulSet") {
			sts, err := ImportStatefulSetConfig([]byte(seg))
			if err != nil {
				return fmt.Errorf("error parsing StatefulSet yaml (segment %d), %s", index, err)
			}
			sd.StatefulSet = sts
			sd.KubeMQManifestConfig.Name = sts.Name
			sd.KubeMQManifestConfig.Namespace = sts.Namespace
			continue
		}
		if strings.Contains(seg, "kind: Service") {
			svc, err := ImportServiceConfig([]byte(seg))
			if err != nil {
				return fmt.Errorf("error parsing Service yaml (segment %d), %s", index, err)
			}
			sd.Services[svc.Name] = svc
			continue
		}
		if strings.Contains(seg, "kind: ConfigMap") {
			cm, err := ImportConfigMap([]byte(seg))

			if err != nil {
				return fmt.Errorf("error parsing ConfigMap yaml (segment %d), %s", index, err)
			}
			sd.ConfigMaps[cm.Name] = cm
			continue
		}
		if strings.Contains(seg, "kind: Secret") {

			sec, err := ImportSecret([]byte(seg))
			if err != nil {
				return fmt.Errorf("error parsing Secret yaml (segment %d), %s", index, err)
			}
			sd.Secrets[sec.Name] = sec
			continue
		}
		if strings.Contains(seg, "kind: Ingress") {

			ingress, err := ImportIngress([]byte(seg))
			if err != nil {
				return fmt.Errorf("error parsing Ingress yaml (segment %d), %s", index, err)
			}
			sd.Ingress[ingress.Name] = ingress
			continue
		}
	}

	return sd.Validate()
}

func (sd *KubeMQDeployment) Validate() error {
	if sd.StatefulSet == nil {
		return fmt.Errorf("invalid import file, no StatefulSet was defined")
	}
	sts, err := sd.StatefulSet.Get()
	if err != nil {
		return fmt.Errorf("invalid import file, no StatefulSet was defined")
	}
	if len(sts.Spec.Template.Spec.Containers) > 0 {
		if !strings.Contains(sts.Spec.Template.Spec.Containers[0].Image, "kubemq") {
			return fmt.Errorf("invalid KubeMQ StatefulSet definitions, docker image is invalid")
		}
	}
	if len(sd.Services) == 0 {
		return fmt.Errorf("invalid KubeMQ Services definitions, at least one serivce must be defined")

	}
	return nil
}
