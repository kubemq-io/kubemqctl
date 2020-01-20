package deployment

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var defaultKubeMQStatefulSetTemplate = `
apiVersion: apps/v1
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
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: {{.Name}}
        deployment.id: {{.Id}}
      annotations:
        prometheus.io/scrape: 'true'
        prometheus.io/port: '9102'
        prometheus.io/path: '/metrics'
    spec:
      containers:
        - env:
            - name: CLUSTER_NAME
              value: {{.Name}}
            - name: CLUSTER_ENABLE
              value: 'true'
          envFrom:
            - secretRef:
                name: {{.Name}}
            - configMapRef:
                name: {{.Name}}
          image: 'kubemq/kubemq:{{.ImageTag}}'
          imagePullPolicy: Always
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
            - containerPort: 7000
              name: gateway-port
              protocol: TCP
{{if (gt .Volume 0)  }}
          volumeMounts:
            - name: {{.Name}}-vol
              mountPath: './store'
{{end}}
      restartPolicy: Always
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

type StatefulSetConfig struct {
	Id          string
	Name        string
	Namespace   string
	ImageTag    string
	Replicas    int
	Volume      int
	statefulset *appsv1.StatefulSet
}

func ImportStatefulSetConfig(spec []byte) (*StatefulSetConfig, error) {
	sts := &appsv1.StatefulSet{}
	err := yaml.Unmarshal(spec, sts)
	if err != nil {
		return nil, err
	}
	return &StatefulSetConfig{
		Id:          "",
		Name:        sts.Name,
		Namespace:   sts.Namespace,
		ImageTag:    "",
		Replicas:    0,
		Volume:      0,
		statefulset: sts,
	}, nil
}

func NewStatefulSetConfig(id, name, namespace string) *StatefulSetConfig {
	return &StatefulSetConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		ImageTag:  "",
		Replicas:  0,
		Volume:    0,
	}
}
func DefaultStatefulSetConfig(id, name, namespace string) *StatefulSetConfig {
	return &StatefulSetConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		ImageTag:  "latest",
		Replicas:  3,
		Volume:    0,
	}
}

func (sc *StatefulSetConfig) SetReplicas(value int) *StatefulSetConfig {
	sc.Replicas = value
	return sc
}

func (sc *StatefulSetConfig) SetVolume(value int) *StatefulSetConfig {
	sc.Volume = value
	return sc
}
func (sc *StatefulSetConfig) SetImageTag(value string) *StatefulSetConfig {
	sc.ImageTag = value
	return sc
}

func (sc *StatefulSetConfig) Spec() ([]byte, error) {
	if sc.statefulset == nil {
		t := NewTemplate(defaultKubeMQStatefulSetTemplate, sc)
		return t.Get()
	}
	return yaml.Marshal(sc.statefulset)
}
func (sc *StatefulSetConfig) Set(value *appsv1.StatefulSet) *StatefulSetConfig {
	sc.statefulset = value
	return sc
}

func (sc *StatefulSetConfig) Get() (*appsv1.StatefulSet, error) {
	if sc.statefulset != nil {
		return sc.statefulset, nil
	}
	sts := &appsv1.StatefulSet{}
	data, err := sc.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, sts)
	if err != nil {
		return nil, err
	}
	sc.statefulset = sts
	return sts, nil
}
