package temporal_test

import (
	"testing"

	"github.com/foomo/testcontainers-go/modules/temporal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.temporal.io/api/workflowservice/v1"
	temporalclient "go.temporal.io/sdk/client"
)

func TestRun(t *testing.T) {
	// t.Skip("only for manual testing")
	t.Parallel()

	container, err := temporal.Run(t.Context(), "temporalio/temporal:latest")
	require.NoError(t, err)

	hostPort, err := container.HostPort(t.Context())
	require.NoError(t, err)

	c, err := temporalclient.NewLazyClient(temporalclient.Options{
		HostPort: hostPort,
	})
	require.NoError(t, err)
	t.Cleanup(c.Close)

	resp, err := c.WorkflowService().ListNamespaces(t.Context(), &workflowservice.ListNamespacesRequest{})
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Namespaces)
}
