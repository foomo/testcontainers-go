# RustFS Module

[RustFS](https://github.com/rustfs/rustfs) is a high-performance, S3-compatible distributed object storage system written in Rust. Use this module to run RustFS in a testcontainer during integration testing.

## Installation

```bash
go get github.com/foomo/testcontainers-go/modules/rustfs
```

## Usage

Import the module and start a container:

```go
package mypackage

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/foomo/testcontainers-go/modules/rustfs"
)

func TestObjectStorage(t *testing.T) {
	container, err := rustfs.Run(t.Context(), "rustfs/rustfs:latest")
	if err != nil {
		t.Fatal(err)
	}

	endpoint, err := container.Endpoint(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       "us-east-1",
		UsePathStyle: true,
		Credentials: credentials.NewStaticCredentialsProvider(
			container.AccessKey,
			container.SecretKey,
			"",
		),
	})

	resp, err := client.ListBuckets(t.Context(), &s3.ListBucketsInput{})
	if err != nil {
		t.Fatal(err)
	}
	_ = resp
}
```

## Exposed Ports

The RustFS container exposes two ports:

| Port | Protocol | Purpose |
|------|----------|---------|
| `9000` | HTTP | S3-compatible API endpoint |
| `9001` | HTTP | Web console UI |

## Container Configuration

The module starts RustFS with the following defaults:

- **Image:** Any `rustfs/rustfs` image (e.g. `rustfs/rustfs:latest`)
- **Exposed Ports:** 9000/tcp, 9001/tcp
- **Credentials:** `RUSTFS_ACCESS_KEY=rustfsadmin`, `RUSTFS_SECRET_KEY=rustfsadmin`

### Getting the Endpoint

Use `Endpoint()` to get the S3 API URL, and `ConsoleEndpoint()` for the web console:

```go
endpoint, err := container.Endpoint(t.Context())
// endpoint is something like "http://localhost:56789"

console, err := container.ConsoleEndpoint(t.Context())
// console is something like "http://localhost:56790"
```

The container struct also exposes `AccessKey` and `SecretKey` fields with the credentials used to start the container.

## Advanced Configuration

Pass additional [`testcontainers.ContainerCustomizer`](https://pkg.go.dev/github.com/testcontainers/testcontainers-go#ContainerCustomizer) options to customize the container, for example to override credentials:

```go
container, err := rustfs.Run(
    t.Context(),
    "rustfs/rustfs:latest",
    testcontainers.WithEnv(map[string]string{
        "RUSTFS_ACCESS_KEY": "myaccesskey",
        "RUSTFS_SECRET_KEY": "mysecretkey",
    }),
)
```

Note: the `AccessKey` / `SecretKey` fields on the container struct will still hold the module defaults — if you override the env vars, track the new values yourself.

Common customizers:

- `WithEnv(map[string]string{...})` — Set environment variables
- `WithLogger(logger)` — Provide a custom logger
- `WithLogConsumers(consumers...)` — Capture logs from the container

See the [testcontainers-go documentation](https://golang.testcontainers.org/) for a full list of options.

## Cleanup

The container is automatically cleaned up when the test ends. For explicit cleanup, use:

```go
err := container.Terminate(t.Context())
if err != nil {
    t.Fatal(err)
}
```

## Related Resources

- [RustFS GitHub](https://github.com/rustfs/rustfs)
- [RustFS Documentation](https://docs.rustfs.com)
- [AWS SDK for Go v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
- [testcontainers-go Documentation](https://golang.testcontainers.org/)
