package runtime

import (
	"cage/container/runtime/colima"
	"cage/container/runtime/docker"
	"cage/container/runtime/podman"
	"context"
	"fmt"
)

const (
	Docker = "docker"
	Podman = "podman"
	Colima = "colima"
)

type Socket struct {
	Host      string
	Available bool
}

func Available(ctx context.Context) map[string]*Socket {
	s := make(map[string]*Socket)

	for runtime, socket := range sockets {
		path, err := socket(ctx)
		if err != nil {
			fmt.Println(runtime, err)
			s[runtime] = nil
			continue
		}

		cli, err := docker.Client(ctx, path)
		if err != nil {
			fmt.Println(runtime, err)
			s[runtime] = nil
			continue
		}

		_, err = cli.Info(ctx)
		if err != nil {
			fmt.Println(runtime, err)
			s[runtime] = &Socket{
				Host:      path,
				Available: false,
			}
			continue
		}

		s[runtime] = &Socket{
			Host:      path,
			Available: true,
		}
	}

	return s
}

type socketFunc func(context.Context) (string, error)

var sockets = map[string]socketFunc{
	Docker: docker.Socket,
	Podman: podman.Socket,
	Colima: colima.Socket,
}
