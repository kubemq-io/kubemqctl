package operator

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var operator = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  replicas: 1
  selector:
    matchLabels:
      name:{{.Name}}
  template:
    metadata:
      labels:
        name: {{.Name}}
    spec:
      serviceAccountName: {{.Name}}
      containers:
        - name: {{.Name}}
          image: {{.Image}}
          command:
          - kubemq-operator
          imagePullPolicy: Always
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "{{.Name}}"
            - name: KUBEMQ_REGISTRY
              value: "{{.KubeMQRegistry}}"
            - name: KUBEMQ_REPOSITORY
              value: "{{.KubeMQRepository}}"
            - name: KUBEMQ_IMAGE_TAG
              value: "{{.KubeMQTag}}"
`

type Operator struct {
	Name             string
	Namespace        string
	Image            string
	KubeMQRegistry   string
	KubeMQRepository string
	KubeMQTag        string
	deployment       *appsv1.Deployment
}

func CreateOperator(name, namespace string) *Operator {
	return &Operator{
		Name:             name,
		Namespace:        namespace,
		Image:            "docker.io/kubemq/kubemq-operator:latest",
		KubeMQRegistry:   "docker.io",
		KubeMQRepository: "kubemq/kubemq",
		KubeMQTag:        "latest",
		deployment:       nil,
	}
}
func (op *Operator) Spec() ([]byte, error) {
	t := NewTemplate(operator, op)
	return t.Get()
}
func (op *Operator) Get() (*appsv1.Deployment, error) {
	if op.deployment != nil {
		return op.deployment, nil
	}
	deployment := &appsv1.Deployment{}
	data, err := op.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, deployment)
	if err != nil {
		return nil, err
	}
	op.deployment = deployment
	return deployment, nil
}
