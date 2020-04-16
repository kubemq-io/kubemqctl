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
      name: {{.Name}}
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
            - name: RELATED_IMAGE_KUBEMQ_CLUSTER
              value: {{.KubemqImage}}
            - name: RELATED_IMAGE_PROMETHEUS
              value: {{.PrometheusImage}}
            - name: RELATED_IMAGE_GRAFANA
              value: {{.GrafanaImage}}
            - name: KUBEMQ_VIEW_DASHBOARD_SOURCE
              value: {{.KubemqDashboardDashboardSource}}
            - name: KUBEMQ_LICENSE_MODE
              value: "COMMUNITY"
`

type Operator struct {
	Name                           string
	Namespace                      string
	Image                          string
	KubemqImage                    string
	PrometheusImage                string
	GrafanaImage                   string
	KubemqDashboardDashboardSource string
	deployment                     *appsv1.Deployment
}

func CreateOperator(name, namespace string) *Operator {
	return &Operator{
		Name:                           name,
		Namespace:                      namespace,
		Image:                          "docker.io/kubemq/kubemq-operator:latest",
		KubemqImage:                    "docker.io/kubemq/kubemq:latest",
		PrometheusImage:                "prom/prometheus:latest",
		GrafanaImage:                   "grafana/grafana:latest",
		KubemqDashboardDashboardSource: "https://raw.githubusercontent.com/kubemq-io/kubemq-dashboard/master/dashboard.json",
		deployment:                     nil,
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
