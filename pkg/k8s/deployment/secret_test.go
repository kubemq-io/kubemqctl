package deployment

import (
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestSecretConfig_Spec(t *testing.T) {
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
				Variables: map[string]string{"LICENSE_KEY_DATA": `-------License---------qweqwweqweqweq`},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SecretConfig{
				Id:        tt.fields.Id,
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
				Variables: tt.fields.Variables,
			}
			got, err := s.Spec()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				sec := &apiv1.Secret{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Data:       nil,
					StringData: nil,
					Type:       "",
				}
				err := yaml.Unmarshal(got, sec)
				require.NoError(t, err)
				deploymentId := sec.Labels["deployment.id"]
				assert.EqualValues(t, tt.fields.Id, deploymentId)
				assert.EqualValues(t, tt.fields.Name, sec.Name)
				assert.EqualValues(t, tt.fields.Namespace, sec.Namespace)
				for key, value := range tt.fields.Variables {
					valueData := sec.StringData[key]
					assert.EqualValues(t, value, valueData)
				}

			}
		})
	}
}
