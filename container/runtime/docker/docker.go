package docker

import (
	"cage/errors"
	"context"
	"os"

	"github.com/docker/docker/client"
)

func Socket(ctx context.Context) (string, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost != "" {
		return "", nil
	}

	return "", errors.MissingDockerSocketErr{}
}

func Client(ctx context.Context, host string) (*client.Client, error) {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	_, err = cli.Info(ctx)
	if err != nil {
		cli.Close()
		return nil, err
	}

	return cli, err
}
