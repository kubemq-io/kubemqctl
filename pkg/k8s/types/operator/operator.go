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
      app: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
    spec:
      serviceAccountName: kubemq-operator
      containers:
        - name: kubemq-operator
          image: docker.io/kubemq/kubemq-operator:latest
          command:
            - kubemq-operator
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
            - containerPort: 8090
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: SOURCE
              value: "kubemqctl"
            - name: DEBUG_MODE
              value: "false"
            - name: RELATED_IMAGE_KUBEMQ_CLUSTER
              value: "docker.io/kubemq/kubemq:latest"
            - name: RELATED_IMAGE_PROMETHEUS
              value: "prom/prometheus:latest"
            - name: RELATED_IMAGE_GRAFANA
              value: "grafana/grafana:latest"
            - name: KUBEMQ_VIEW_DASHBOARD_SOURCE
              value: "https://raw.githubusercontent.com/kubemq-io/kubemq-dashboard/master/dashboard.json"
            - name: OPERATOR_NAME
              value: "kubemq-operator"
            - name: KUBEMQ_LICENSE_MODE
              value: "COMMUNITY"
            - name: CONNECTOR_TARGETS_IMAGE
              value: "kubemq/kubemq-targets:latest"
            - name: CONNECTOR_SOURCES_IMAGE
              value: "kubemq/kubemq-sources:latest"
            - name: CONNECTOR_BRIDGES_IMAGE
              value: "kubemq/kubemq-bridges:latest"
            - name: CONNECTOR_TASKS_IMAGE
              value: "kubemq/kubemq-tasks:latest"
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
