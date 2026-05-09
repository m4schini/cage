package cage

import (
	"cage/cage/containerfile"
	"cage/cage/state"
	ctr "cage/container"
	"cage/nix"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type ContainerRunner interface {
	Run(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) error
}

func Run(ctx context.Context, cli *client.Client, name string, runner ContainerRunner) error {
	cfg, modTime, err := LoadConfig()
	if err != nil {
		return err
	}
	imageName := "localhost/cage:" + cfg.Name
	packages := cfg.Packages
	packages = append(packages, "claude-code")

	imageCreatedTime, err := ctr.ImageCreatedAt(ctx, cli, imageName)
	if err != nil || modTime.After(imageCreatedTime) {
		err = ctr.BuildImage(ctx, cli, containerfile.Containerfile, []string{imageName}, nix.ShellNixPackages{
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

func Load(name string) (state.CageDefinition, error) {
	root, err := state.DataDir.OpenRoot(name)
	if err != nil {
		return state.CageDefinition{}, err
	}

	f, err := root.Open("cage.yaml")
	if err != nil {
		return state.CageDefinition{}, err
	}

	return state.Read(f)
}

func New(name string, definition state.CageDefinition) error {
	err := checkDirAvailability(name)
	if err != nil {
		return err
	}

	err = state.DataDir.Mkdir(name, 0750)
	if err != nil {
		return err
	}

	root, err := state.DataDir.OpenRoot(name)
	if err != nil {
		return err
	}

	f, err := root.Create(state.DefinitionFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return state.Write(definition, f)
}

func validateCageName(name string) error {
	//TODO
	return nil
}

func checkDirAvailability(name string) error {
	fi, err := state.DataDir.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return ErrNameConflict{
			Name:   name,
			Reason: "file with name exists",
		}
	}

	return ErrAlreadyExists{Name: name}
}

type ErrAlreadyExists struct {
	Name string
}

func (a ErrAlreadyExists) Error() string {
	return fmt.Sprintf(`cage "%v" already exists`, a.Name)
}

type ErrNameConflict struct {
	Name   string
	Reason string
}

func (n ErrNameConflict) Error() string {
	return fmt.Sprintf(`cage name (%v) is unavailable: %v`, n.Name, n.Reason)
}
