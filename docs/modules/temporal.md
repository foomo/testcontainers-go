# Temporal Module

[Temporal](https://temporal.io) is an open-source workflow engine designed for building resilient distributed applications. Use this module to run Temporal in a testcontainer during integration testing.

## Installation

```bash
go get github.com/foomo/testcontainers-go/modules/temporal
```

## Usage

Import the module and start a container:

```go
package mypackage

import (
	"testing"

	"github.com/foomo/testcontainers-go/modules/temporal"
	"go.temporal.io/api/workflowservice/v1"
	temporalclient "go.temporal.io/sdk/client"
)

func TestWorkflow(t *testing.T) {
	// Start a Temporal container
	container, err := temporal.Run(t.Context(), "temporalio/temporal:latest")
	if err != nil {
		t.Fatal(err)
	}

	// Get the gRPC endpoint
	hostPort, err := container.HostPort(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// Create a Temporal client
	c, err := temporalclient.NewLazyClient(temporalclient.Options{
		HostPort: hostPort,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Use the client in your tests
	resp, err := c.WorkflowService().ListNamespaces(t.Context(), &workflowservice.ListNamespacesRequest{})
	if err != nil {
		t.Fatal(err)
	}

	// Assert on the response...
	if len(resp.Namespaces) == 0 {
		t.Fatal("Expected namespaces")
	}
}
```

## Exposed Ports

The Temporal container exposes three ports:

| Port | Protocol | Purpose |
|------|----------|---------|
| `7233` | gRPC | Temporal frontend service (default for SDK clients) |
| `8233` | HTTP | Temporal frontend HTTP service |
| `8080` | HTTP | Temporal Web UI (http://localhost:8080 during local development) |

## Container Configuration

The module starts Temporal with development defaults:

- **Image:** Any `temporalio/temporal` image (defaults to `temporalio/temporal:latest`)
- **Command:** `server start-dev --ip 0.0.0.0`
- **Exposed Ports:** 7233/tcp, 8233/tcp, 8080/tcp

### Getting the Endpoint

Use `HostPort()` to get the gRPC endpoint for connecting clients:

```go
hostPort, err := container.HostPort(t.Context())
if err != nil {
    t.Fatal(err)
}
// hostPort is something like "localhost:56789"
```

## Advanced Configuration

Pass additional [`testcontainers.ContainerCustomizer`](https://pkg.go.dev/github.com/testcontainers/testcontainers-go#ContainerCustomizer) options to customize the container:

```go
container, err := temporal.Run(
    t.Context(),
    "temporalio/temporal:latest",
    testcontainers.WithEnv(map[string]string{
        "TEMPORAL_ENVIRONMENT": "development",
    }),
    testcontainers.WithLogger(logger),
)
```

Common customizers:

- `WithEnv(map[string]string{...})` — Set environment variables
- `WithLogger(logger)` — Provide a custom logger
- `WithLogConsumers(consumers...)` — Capture logs from the container

See the [testcontainers-go documentation](https://golang.testcontainers.org/) for a full list of options.

## Example: Testing a Workflow

```go
package myworkflows

import (
	"testing"

	"github.com/foomo/testcontainers-go/modules/temporal"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// MyWorkflow is a simple workflow
func MyWorkflow(ctx workflow.Context) (string, error) {
	// Implementation...
	return "result", nil
}

// MyActivity is a simple activity
func MyActivity(ctx context.Context) (string, error) {
	// Implementation...
	return "activity result", nil
}

func TestMyWorkflow(t *testing.T) {
	// Start Temporal
	container, err := temporal.Run(t.Context(), "temporalio/temporal:latest")
	if err != nil {
		t.Fatal(err)
	}

	hostPort, err := container.HostPort(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// Create client
	c, err := client.NewLazyClient(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Create a worker
	w := worker.New(c, "my-task-queue", worker.Options{})

	// Register workflow and activity
	w.RegisterWorkflow(MyWorkflow)
	w.RegisterActivity(MyActivity)

	// Start the worker
	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	// Execute workflow
	run, err := c.ExecuteWorkflow(t.Context(), client.StartWorkflowOptions{
		TaskQueue: "my-task-queue",
	}, MyWorkflow)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for result
	var result string
	if err := run.Get(t.Context(), &result); err != nil {
		t.Fatal(err)
	}

	if result != "expected-result" {
		t.Fatalf("unexpected result: %v", result)
	}
}
```

## Cleanup

The container is automatically cleaned up when the test ends. For explicit cleanup, use:

```go
err := container.Terminate(t.Context())
if err != nil {
    t.Fatal(err)
}
```

## Related Resources

- [Temporal Documentation](https://docs.temporal.io)
- [Temporal Go SDK](https://pkg.go.dev/go.temporal.io/sdk)
- [testcontainers-go Documentation](https://golang.testcontainers.org/)
