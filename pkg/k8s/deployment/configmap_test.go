package deployment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfigMapConfig_Spec(t *testing.T) {
	type fields struct {
		Id        string
		Name      string
		Namespace string
		Variables map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "full",
			fields: fields{
				Id:        "some-id",
				Name:      "kubemq-sec",
				Namespace: "kubemq-namespace",
				Variables: map[string]string{"CONFIG": `-------License---------qweqwweqweqweq`},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConfigMapConfig{
				Id:        tt.fields.Id,
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
				Variables: tt.fields.Variables,
			}
			cm, err := c.Get()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NoError(t, err)
				deploymentId := cm.Labels["deployment.id"]
				assert.EqualValues(t, tt.fields.Id, deploymentId)
				assert.EqualValues(t, tt.fields.Name, cm.Name)
				assert.EqualValues(t, tt.fields.Namespace, cm.Namespace)
				for key, value := range tt.fields.Variables {
					valueData := cm.Data[key]
					assert.EqualValues(t, value, valueData)
				}
			}
		})
	}
}
