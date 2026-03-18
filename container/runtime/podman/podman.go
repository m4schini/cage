package podman

import (
	"cage/errors"
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func Socket(ctx context.Context) (string, error) {
	dockerHost := os.Getenv("CONTAINER_HOST")
	if dockerHost != "" {
		return "", nil
	}

	defaultSocketPath := filepath.Join("/run/user", fmt.Sprintf("%d", os.Getuid()), "podman/podman.sock")
	f, err := os.Stat(defaultSocketPath)
	if err != nil {
		return "", err
	}
	if f.IsDir() {
		return "", errors.MissingDockerSocketErr{}
	}

	return "unix://" + defaultSocketPath, nil
}
