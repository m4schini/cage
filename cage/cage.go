package cage

import (
	"cage/cage/config"
	"cage/cage/containerfile"
	"cage/cage/state"
	ctr "cage/container"
	"cage/nix"
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

func Run(ctx context.Context, cli *client.Client) error {
	cfg, modTime, err := state.Load()
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("localhost/%v:%v", config.AppName, cfg.Name)
	packages := cfg.Packages
	packages = append(packages, "claude-code")

	imageCreatedTime, err := ctr.ImageCreatedAt(ctx, cli, imageName)
	if err != nil || modTime.After(imageCreatedTime) {
		err = ctr.BuildImage(ctx, cli, true, containerfile.Containerfile, []string{imageName}, nix.ShellNixPackages{
			Packages: packages,
			Shell:    "bash",
		})
		if err != nil {
			return err
		}
	}

	env, err := cfg.ClaudeCodeEnv()
	if err != nil {
		return err
	}

	err = ctr.RunImage(ctx, cli, imageName, env)
	if err != nil {
		return err
	}

	return nil
}
