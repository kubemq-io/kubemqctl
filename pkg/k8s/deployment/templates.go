package deployment

import (
	"text/template"
)

var defaultStsTemplate = `
apiVersion: {{.ApiVersion}}
kind: StatefulSet
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  selector:
    matchLabels:
      app: {{.Name}}
  replicas: {{.Replicas}}
  serviceName: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
      annotations:
        prometheus.io/scrape: 'true'
        prometheus.io/port: '9102'
        prometheus.io/path: '/metrics'
    spec:
      containers:
        - env:
            - name: KUBEMQ_TOKEN
              value: {{.Token}}
            - name: CLUSTER_ROUTES
              value: '{{.Name}}:5228'
            - name: CLUSTER_PORT
              value: '5228'
            - name: CLUSTER_ENABLE
              value: 'true'
            - name: GRPC_PORT
              value: '50000'
            - name: REST_PORT
              value: '9090'
            - name: KUBEMQ_PORT
              value: '8080'
            - name: STORE_DIR
              value: '/store'
          image: 'kubemq/kubemq:{{.Version}}'
          name: {{.Name}}
          ports:
            - containerPort: 50000
              name: grpc-port
              protocol: TCP
            - containerPort: 8080
              name: api-port
              protocol: TCP
            - containerPort: 9090
              name: rest-port
              protocol: TCP
            - containerPort: 5228
              name: cluster-port
              protocol: TCP
{{if (gt .Volume 0)  }}
          volumeMounts:
            - name: {{.Name}}-vol
              mountPath: '/store'
{{end}}
{{if (gt .Volume 0)  }}  
  volumeClaimTemplates:
    - metadata:
        name: {{.Name}}-vol
      spec:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName:
        resources:
          requests:
            storage: {{.Volume}}Gi
{{end}}
`

var defaultServiceTemplate = `
apiVersion: {{.ApiVersion}}
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  ports:
    - name: {{.PortName}}
      port: {{.ContainerPort}}
      protocol: {{.Protocol}}
      targetPort: {{.TargetPort}}
  sessionAffinity: None
  type: {{.Type}}
  selector:
    app: {{.AppName}}
`

var defaultNameSpaceTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Namespace}}
`

type Template struct {
	Structure string
	Data      interface{}
	output    []byte
}

func NewTemplate(str string, data interface{}) *Template {
	return &Template{
		Structure: str,
		Data:      data,
	}
}

func (t *Template) Write(p []byte) (n int, err error) {
	t.output = append(t.output, p...)
	return len(t.output), nil
}
func (t *Template) Get() ([]byte, error) {
	tmpl, err := template.New("tmpl").Parse(t.Structure)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(t, t.Data)
	if err != nil {
		return nil, err
	}
	return t.output, nil
}
