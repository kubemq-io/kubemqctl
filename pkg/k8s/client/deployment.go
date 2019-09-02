package client

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"strings"

	"time"
)

type StatefulSetStatus struct {
	Name      string
	Namespace string
	Desired   int32
	Running   int32
	Ready     int32
	Image     string
	Age       time.Duration
	PVC       bool
}

type ServiceStatus struct {
	Name      string
	Namespace string
	Type      string
	ClusterIP string
	ExternalP string
	Ports     string
	Age       time.Duration
}

type StatefulSetDeployment struct {
	Namespace         string
	Name              string
	Labels            map[string]string
	StatefulSet       *appsv1.StatefulSet
	Services          []apiv1.Service
	VolumesClaims     []apiv1.PersistentVolumeClaim
	StatefulSetStatus *StatefulSetStatus
	ServicesStatus    []*ServiceStatus
}

func (d *StatefulSetDeployment) ServicesStatusString() string {
	list := []string{}
	for _, ss := range d.ServicesStatus {
		if ss.ExternalP == "" {
			list = append(list, fmt.Sprintf(" %s:%s", ss.ClusterIP, ss.Ports))
		} else {
			list = append(list, fmt.Sprintf(" %s/%s:%s", ss.ExternalP, ss.ClusterIP, ss.Ports))
		}

	}
	return strings.Join(list, ",")
}
