package builder

import (
	"fmt"
)

var operatorTemplate = `
kind: Namespace
apiVersion: v1
metadata:
  name: {{.Namespace}}
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubemq-operator-kubemq-{{.Namespace}}-crb
subjects:
  - kind: ServiceAccount
    name: kubemq-operator
    namespace: {{.Namespace}}
roleRef:
  kind: ClusterRole
  name: kubemq-operator
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kubemq-cluster
  namespace: {{.Namespace}}
rules:
  - apiGroups:
      - security.openshift.io
    resources:
      - securitycontextconstraints
    verbs:
      - use
      - delete
      - get
      - list
      - patch
      - update
      - watch
    resourceNames:
      - privileged
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubemq-cluster
  namespace: {{.Namespace}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubemq-operator
  namespace: {{.Namespace}}
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubemq-cluster-{{.Namespace}}-rb
  namespace: {{.Namespace}}
subjects:
  - kind: ServiceAccount
    name: kubemq-cluster
    namespace: {{.Namespace}}
roleRef:
  kind: Role
  name: kubemq-cluster
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubemq-operator
  namespace: {{.Namespace}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubemq-operator
  template:
    metadata:
      labels:
        app: kubemq-operator
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
            - name: RELATED_IMAGE_KUBEMQ_CLUSTER
              value: "docker.io/kubemq/kubemq:latest"
            - name: CONNECTOR_TARGETS_IMAGE
              value: "kubemq/kubemq-targets:latest"
            - name: CONNECTOR_SOURCES_IMAGE
              value: "kubemq/kubemq-sources:latest"
            - name: CONNECTOR_BRIDGES_IMAGE
              value: "kubemq/kubemq-bridges:latest"
`

type OperatorManifest struct {
	Namespace string
}

func NewOperatorManifest() *OperatorManifest {
	return &OperatorManifest{}
}

func (o *OperatorManifest) SetNamespace(value string) *OperatorManifest {
	o.Namespace = value
	return o
}

func (o *OperatorManifest) Manifest() (string, error) {
	t := NewTemplate(operatorTemplate, o)
	data, err := t.Get()
	if err != nil {
		return "", fmt.Errorf("error during deployment manifest generation, error: %s", err.Error())
	}
	return string(data), nil
}
