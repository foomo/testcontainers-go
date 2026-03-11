package temporal

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

type TemporalContainer struct {
	testcontainers.Container
}

func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*TemporalContainer, error) {
	moduleOpts := []testcontainers.ContainerCustomizer{
		testcontainers.WithEnv(map[string]string{}),
		testcontainers.WithExposedPorts("7233/tcp", "8233/tcp", "8080/tcp"),
		testcontainers.WithEntrypoint("/usr/local/bin/temporal"),
		testcontainers.WithCmd("server", "start-dev", "--ip", "0.0.0.0"),
	}

	moduleOpts = append(moduleOpts, opts...)

	ctr, err := testcontainers.Run(ctx, img, moduleOpts...)

	var c *TemporalContainer
	if ctr != nil {
		c = &TemporalContainer{
			Container: ctr,
		}
	}

	if err != nil {
		return c, fmt.Errorf("run temporal: %w", err)
	}

	return c, nil
}

func (c *TemporalContainer) HostPort(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, "7233/tcp", "")
}
