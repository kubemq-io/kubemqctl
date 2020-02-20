package deployment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNamespaceConfig(t *testing.T) {
	type fields struct {
		Id   string
		Name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "full",
			fields: fields{
				Id:   "some-id",
				Name: "namespace",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NamespaceConfig{
				Id:   tt.fields.Id,
				Name: tt.fields.Name,
			}
			ns, err := n.Get()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				deploymentId := ns.Labels["deployment.id"]
				assert.EqualValues(t, tt.fields.Id, deploymentId)
				assert.EqualValues(t, tt.fields.Name, ns.Name)
			}
		})
	}
}
