package logs

import (
	"context"
	"github.com/kubemq-io/kubetools/k8s/pkg/client"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRun(t *testing.T) {
	c, err := client.NewClient("")
	require.NoError(t, err)
	v, err := c.ClientSet.ServerVersion()
	require.NoError(t, err)
	require.NotNil(t, v)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = Run(ctx, c, nil)
	require.NoError(t, err)
}
