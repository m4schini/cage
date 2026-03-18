package colima

import (
	"cage/errors"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type StatusJSON struct {
	DisplayName  string `json:"display_name"`
	Driver       string `json:"driver"`
	Arch         string `json:"arch"`
	Runtime      string `json:"runtime"`
	DockerSocket string `json:"docker_socket"`
	CPU          int    `json:"cpu"`
	Memory       int    `json:"memory"`
	Disk         int    `json:"disk"`
}

func Socket(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "colima", "status", "--json").Output()
	if err != nil {
		return "", err
	}

	var status StatusJSON
	err = json.Unmarshal(out, &status)
	if err != nil {
		return "", err
	}

	if status.Runtime != "docker" {
		return "", UnsupportedRuntimeErr{Runtime: status.Runtime}
	}

	if status.DockerSocket == "" {
		return "", errors.MissingDockerSocketErr{}
	}

	return status.DockerSocket, nil
}

type UnsupportedRuntimeErr struct {
	Runtime string
}

func (u UnsupportedRuntimeErr) Error() string {
	return fmt.Sprintf("colima: %v unsupported (start colima with docker runtime)", u.Runtime)
}
