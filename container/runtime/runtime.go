package runtime

import (
	"cage/container/runtime/docker"
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

func Client(ctx context.Context, runtime string) (*client.Client, error) {
	getSocket, ok := sockets[runtime]
	if !ok {
		return nil, UnsupportedRuntimeErr{Runtime: runtime}
	}

	host, err := getSocket(ctx)
	if err != nil {
		return nil, err
	}

	return docker.Client(ctx, host)
}

type UnsupportedRuntimeErr struct {
	Runtime string
}

func (u UnsupportedRuntimeErr) Error() string {
	return fmt.Sprintf("unsupported runtime: %v", u.Runtime)
}
