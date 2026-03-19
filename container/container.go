package container

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/term"
)

type Docker struct {
	Client *client.Client
}

func (d *Docker) Run(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) error {
	cli := d.Client
	containerName = strings.Join([]string{"cage", containerName}, "-")
	r, err := cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return err
	}

	// Start the container
	err = cli.ContainerStart(ctx, r.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	// Attach to the container
	attachResp, err := cli.ContainerAttach(ctx, r.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer attachResp.Close()

	// Set terminal to raw mode for proper TTY interaction
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Connect stdin/stdout - only ONE copy from Reader to Stdout
	go io.Copy(os.Stdout, attachResp.Reader)
	go io.Copy(attachResp.Conn, os.Stdin)

	// Wait for container to finish OR interrupt signal
	statusCh, errCh := cli.ContainerWait(ctx, r.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
		// Container exited normally
	case <-ctx.Done():
		// Interrupted by user (Ctrl+C)
		fmt.Println("\nDetaching...")
	}

	return errors.Join(
		cli.ContainerStop(ctx, r.ID, container.StopOptions{}),
		cli.ContainerRemove(ctx, r.ID, container.RemoveOptions{}),
	)
}
