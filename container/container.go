package container

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func Run(ctx context.Context) error {
	const NAME = "cage3"
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	fmt.Println("started client")
	defer client.Close()

	ctx = namespaces.WithNamespace(ctx, "cage")
	c, err := client.LoadContainer(ctx, NAME)
	if err != nil {
		return err
	}
	c.Delete(ctx)

	// Pull the image if it doesn't exist
	image, err := client.Pull(ctx, "codeberg.org/aur0ra/cage-base:latest", containerd.WithPullUnpack)
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	fmt.Println("pulled image")

	// Create a container
	container, err := client.NewContainer(
		ctx,
		NAME,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(NAME+"-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	fmt.Println("created container")

	// Create a task (running instance of the container)
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	defer task.Delete(ctx)
	fmt.Println("created task")

	// Wait for the task to exit
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for task: %w", err)
	}

	// Start the task
	if err := task.Start(ctx); err != nil {
		return fmt.Errorf("failed to start task: %w", err)
	}
	fmt.Println("started container")

	// Wait for the task to finish
	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return fmt.Errorf("failed to get task result: %w", err)
	}

	fmt.Printf("Container exited with code: %d\n", code)

	return nil
}
