package rustfs

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultAccessKey = "rustfsadmin"
	defaultSecretKey = "rustfsadmin"
)

type RustFSContainer struct {
	testcontainers.Container
	AccessKey string
	SecretKey string
}

func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*RustFSContainer, error) {
	moduleOpts := []testcontainers.ContainerCustomizer{
		testcontainers.WithEnv(map[string]string{
			"RUSTFS_ACCESS_KEY":                         defaultAccessKey,
			"RUSTFS_SECRET_KEY":                         defaultSecretKey,
			"RUSTFS_ALLOW_INSECURE_DEFAULT_CREDENTIALS": "true",
		}),
		testcontainers.WithExposedPorts("9000/tcp", "9001/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("9000/tcp").WithStartupTimeout(time.Minute),
		),
	}

	moduleOpts = append(moduleOpts, opts...)

	ctr, err := testcontainers.Run(ctx, img, moduleOpts...)

	var c *RustFSContainer
	if ctr != nil {
		c = &RustFSContainer{
			Container: ctr,
			AccessKey: defaultAccessKey,
			SecretKey: defaultSecretKey,
		}
	}

	if err != nil {
		return c, fmt.Errorf("run rustfs: %w", err)
	}

	return c, nil
}

func (c *RustFSContainer) Endpoint(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, "9000/tcp", "http")
}

func (c *RustFSContainer) ConsoleEndpoint(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, "9001/tcp", "http")
}
