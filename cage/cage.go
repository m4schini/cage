package cage

import (
	"cage/cage/state"
	"cage/nix"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

type ContainerRunner interface {
	Run(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) error
}

func Run(ctx context.Context, name string, runner ContainerRunner) error {
	err := validateCageName(name)
	if err != nil {
		return err
	}
	cagePath := filepath.Join(state.DataDirPath, name)
	cageDefinition, err := Load(name)
	if err != nil {
		return err
	}

	shellNixPath := filepath.Join(cagePath, "shell.nix")
	if _, err := os.Stat(shellNixPath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(shellNixPath)
		if err != nil {
			return err
		}
		err = nix.NewNixShell(nix.ShellNixPackages{
			Packages: cageDefinition.Packages,
			Shell:    cageDefinition.Shell,
		}, f)
		f.Close()
		if err != nil {
			return err
		}
	}

	var env []string
	for _, envVar := range cageDefinition.Env {
		env = append(env, fmt.Sprintf("%v=%v", envVar.Key, envVar.Value))
	}

	ctrConfig := &container.Config{
		Hostname:     name,
		Domainname:   name,
		User:         "developer",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    false,
		StdinOnce:    false,
		Env:          env,
		//Cmd:             []string{"claude"},
		Healthcheck:     nil,
		ArgsEscaped:     false,
		Image:           "localhost/my-devcontainer",
		WorkingDir:      "/home/developer/workspace",
		Entrypoint:      []string{"nix-shell", "../shell.nix"},
		NetworkDisabled: false,
		MacAddress:      "",
		OnBuild:         nil,
		Labels:          nil,
		StopSignal:      "",
		StopTimeout:     nil,
		Shell:           nil,
	}

	return runner.Run(ctx, ctrConfig, &container.HostConfig{
		//Binds: []string{
		//	filepath.Join(cagePath, "shell.nix") + ":/home/developer/shell.nix:ro",
		//},
		ContainerIDFile: "",
		LogConfig:       container.LogConfig{},
		NetworkMode:     "",
		PortBindings:    nil,
		RestartPolicy:   container.RestartPolicy{},
		AutoRemove:      false,
		VolumeDriver:    "",
		VolumesFrom:     nil,
		ConsoleSize:     [2]uint{},
		Annotations:     nil,
		CapAdd:          nil,
		CapDrop:         nil,
		CgroupnsMode:    "",
		DNS:             nil,
		DNSOptions:      nil,
		DNSSearch:       nil,
		ExtraHosts:      nil,
		GroupAdd:        nil,
		IpcMode:         "",
		Cgroup:          "",
		Links:           nil,
		OomScoreAdj:     0,
		PidMode:         "",
		Privileged:      false,
		PublishAllPorts: false,
		ReadonlyRootfs:  false,
		SecurityOpt:     nil,
		StorageOpt:      nil,
		Tmpfs:           nil,
		UTSMode:         "",
		UsernsMode:      "",
		ShmSize:         0,
		Sysctls:         nil,
		Runtime:         "",
		Isolation:       "",
		Resources:       container.Resources{},
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   filepath.Join(cagePath, "shell.nix"),
				Target:   "/home/developer/shell.nix",
				ReadOnly: true,
				BindOptions: &mount.BindOptions{
					Propagation: mount.PropagationRPrivate,
				},
			},
			mount.Mount{
				Type:   mount.TypeVolume,
				Source: "nixstore",
				Target: "/nix",
				VolumeOptions: &mount.VolumeOptions{
					DriverConfig: &mount.Driver{},
				},
			},
		},
		MaskedPaths:   nil,
		ReadonlyPaths: nil,
		Init:          nil,
	}, nil, name)
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
