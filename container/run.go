package container

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	xterm "golang.org/x/term"
)

func RunImage(ctx context.Context, cli *client.Client, image string, env []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	env = append(env, "TERM="+term)

	r, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Env:          env,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    false,
		Tty:          true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cwd,
				Target: "/workspace",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, r.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	if cols, rows, err := xterm.GetSize(int(os.Stdin.Fd())); err == nil {
		_ = cli.ContainerResize(ctx, r.ID, container.ResizeOptions{
			Width:  uint(cols),
			Height: uint(rows),
		})
	}

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

	cleanup, err := PrepareTTY()
	if err != nil {
		return err
	}
	defer cleanup()

	go io.Copy(os.Stdout, attachResp.Reader)
	go io.Copy(attachResp.Conn, os.Stdin)

	statusCh, errCh := cli.ContainerWait(ctx, r.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	case <-ctx.Done():
		fmt.Println("\nDetaching...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return cli.ContainerRemove(ctx, r.ID, container.RemoveOptions{Force: true})
}
