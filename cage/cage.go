package cage

import (
	"cage/cage/containerfile"
	ctr "cage/container"
	"cage/nix"
	"context"

	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

func Run(ctx context.Context, cli *client.Client) error {
	cfg, modTime, err := LoadConfig()
	if err != nil {
		return err
	}
	imageName := "localhost/cage:" + cfg.Name
	packages := cfg.Packages
	packages = append(packages, "claude-code")

	imageCreatedTime, err := ctr.ImageCreatedAt(ctx, cli, imageName)
	if err != nil || modTime.After(imageCreatedTime) {
		err = ctr.BuildImage(ctx, cli, viper.GetBool("verbose"), containerfile.Containerfile, []string{imageName}, nix.ShellNixPackages{
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
