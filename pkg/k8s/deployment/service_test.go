package deployment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceConfig_Spec(t *testing.T) {
	type fields struct {
		Id            string
		Name          string
		Namespace     string
		AppName       string
		Type          string
		ContainerPort int
		TargetPort    int
		PortName      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "full",
			fields: fields{
				Id:            "some-id",
				Name:          "kubemq",
				Namespace:     "kubemq-namespace",
				AppName:       "svc",
				Type:          "ClusterIP",
				ContainerPort: 5000,
				TargetPort:    6000,
				PortName:      "kube-rest",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceConfig{
				Id:            tt.fields.Id,
				Name:          tt.fields.Name,
				Namespace:     tt.fields.Namespace,
				AppName:       tt.fields.AppName,
				Type:          tt.fields.Type,
				ContainerPort: tt.fields.ContainerPort,
				TargetPort:    tt.fields.TargetPort,
				PortName:      tt.fields.PortName,
			}
			svc, err := s.Get()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				deploymentId := svc.Labels["deployment.id"]
				assert.EqualValues(t, tt.fields.Id, deploymentId)
				assert.EqualValues(t, tt.fields.Name, svc.Name)
				assert.EqualValues(t, tt.fields.Namespace, svc.Namespace)
				assert.EqualValues(t, tt.fields.Type, svc.Spec.Type)
				assert.EqualValues(t, tt.fields.ContainerPort, svc.Spec.Ports[0].Port)
				assert.EqualValues(t, tt.fields.TargetPort, svc.Spec.Ports[0].TargetPort.IntVal)
				assert.EqualValues(t, tt.fields.PortName, svc.Spec.Ports[0].Name)
			}
		})
	}
}
