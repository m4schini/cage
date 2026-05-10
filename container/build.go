package container

import (
	"archive/tar"
	"bytes"
	"cage/nix"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BuildImage(ctx context.Context, cli *client.Client, dockerfile string, tags []string, packages nix.ShellNixPackages) error {
	shellNix, err := nix.NewNixShellString(packages)
	if err != nil {
		return err
	}

	fmt.Println(shellNix)

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	if err := tw.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Size: int64(len(dockerfile)),
		Mode: 0644,
	}); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(dockerfile)); err != nil {
		return err
	}

	if err := tw.WriteHeader(&tar.Header{
		Name: "shell.nix",
		Size: int64(len(shellNix)),
		Mode: 0644,
	}); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(shellNix)); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	resp, err := cli.ImageBuild(ctx, buf, types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func ImageCreatedAt(ctx context.Context, cli *client.Client, image string) (time.Time, error) {
	info, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, info.Created)
}
