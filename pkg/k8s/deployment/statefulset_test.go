package deployment

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatefulSetConfig_Spec(t *testing.T) {
	type fields struct {
		Id        string
		Name      string
		Namespace string
		Replicas  int
		Volume    int
		ImageTag  string
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
				Name:      "kubemq-test",
				Namespace: "kubemq-namesapce",
				Replicas:  5,
				Volume:    30,
				ImageTag:  "latest",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := StatefulSetConfig{
				Id:        tt.fields.Id,
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
				Replicas:  tt.fields.Replicas,
				Volume:    tt.fields.Volume,
				ImageTag:  tt.fields.ImageTag,
			}
			sts, err := sc.Get()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, tt.fields.Name, sts.Name)
				deploymentId := sts.Spec.Template.Labels["deployment.id"]
				assert.EqualValues(t, tt.fields.Id, deploymentId)
				assert.EqualValues(t, tt.fields.Namespace, sts.Namespace)
				assert.EqualValues(t, tt.fields.Replicas, *sts.Spec.Replicas)
				assert.EqualValues(t, "kubemq/kubemq:"+tt.fields.ImageTag, sts.Spec.Template.Spec.Containers[0].Image)

				if tt.fields.Volume > 0 {
					storage := sts.Spec.VolumeClaimTemplates[0].Spec.Resources.Requests["storage"]
					assert.EqualValues(t, fmt.Sprintf("%dGi", tt.fields.Volume), storage.String())
				}
			}
		})
	}
}
