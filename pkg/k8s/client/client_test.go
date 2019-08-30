package client

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewClientConfigFromHome(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	v, err := c.ClientSet.ServerVersion()
	require.NoError(t, err)
	require.NotNil(t, v)
}

func TestNewClientConfigFromPath(t *testing.T) {
	c, err := NewClient("./test/ClientConfig")
	require.NoError(t, err)
	v, err := c.ClientSet.ServerVersion()
	require.NoError(t, err)
	require.NotNil(t, v)
}

func TestNewClientConfigFromBadConfig(t *testing.T) {
	c, err := NewClient("./test/config_bad")
	require.Error(t, err)
	require.Nil(t, c)
}

func TestNewClientConfigFromBadPath(t *testing.T) {
	c, err := NewClient("./test/some_bad_config")
	require.Error(t, err)
	require.Nil(t, c)
}

func TestGetConfigClusters(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	clusters, err := c.GetConfigClusters()
	require.NoError(t, err)
	require.NotNil(t, clusters)
	for name, cluster := range clusters {
		fmt.Println(name, cluster.Server)
	}
}

func TestClient_GetConfigContext(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	contexts, err := c.GetConfigContext()
	require.NoError(t, err)
	require.NotNil(t, contexts)
	for name, ctx := range contexts {
		fmt.Println(name, ctx.Cluster, ctx.Namespace)
	}
}

func TestClient_GetCurrentContext(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	ctx, err := c.GetCurrentContext()
	require.NoError(t, err)
	require.NotEmpty(t, ctx)
	fmt.Println(ctx)

}

func TestClient_SwitchContext(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	contexts, err := c.GetConfigContext()
	require.NoError(t, err)
	require.NotNil(t, contexts)
	current, err := c.GetCurrentContext()
	require.NoError(t, err)
	require.NotEmpty(t, current)
	switched := false
	for name, _ := range contexts {
		if name != current {
			err := c.SwitchContext(name)
			require.NoError(t, err)
			newCurrent, err := c.GetCurrentContext()
			require.NoError(t, err)
			require.Equal(t, name, newCurrent)
			switched = true
		}

	}
	require.True(t, switched)
}

func TestClient_GetStatefulSets(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	sts, err := c.GetStatefulSets("")
	require.NoError(t, err)
	require.NotEmpty(t, sts)
	fmt.Println(sts)

}

func TestClient_GetService(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	services, err := c.GetServices("", "kube")
	require.NoError(t, err)
	require.NotEmpty(t, services)
	fmt.Println(services)

}
func TestClient_GetPods(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	pods, err := c.GetPods("", "kubemq-cluster")
	require.NoError(t, err)
	require.NotEmpty(t, pods)
	fmt.Println(pods)

}
func TestClient_ForwardPorts(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	stopCh := make(chan struct{})
	outCh, errOutCh := make(chan string, 1), make(chan string, 1)

	err = c.ForwardPorts("default", "kubemq-cluster-0", []string{"50000:50000", "9090:9090", "8080:8080"}, stopCh, outCh, errOutCh)
	require.NoError(t, err)
	select {
	case out := <-outCh:
		fmt.Println(out)
		close(stopCh)
	case err := <-errOutCh:
		fmt.Println(err)
		require.Empty(t, err)
	case <-time.After(5 * time.Second):
		require.NoError(t, errors.New("timeout"))
	}

}
func TestClient_Scale(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	err = c.Scale("default", "kubemq-cluster", 5)
	require.NoError(t, err)
	err = c.Scale("default", "kubemq-cluster", 3)
	require.NoError(t, err)

}

func TestClient_DescribeStatefulSet(t *testing.T) {
	c, err := NewClient("")
	require.NoError(t, err)
	require.NotNil(t, c)
	y, err := c.DescribeStatefulSet("default", "kubemq-cluster")
	require.NoError(t, err)
	require.NotEmpty(t, y)
	fmt.Println(y)

}
