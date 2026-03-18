package cmd

import (
	"cage/config"
	"cage/container/runtime"
	"cage/state"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:    config.AppName,
	Short:  `isolate "ai" "agents"`,
	PreRun: RequireInitialized,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
		defer cancel()
		cli, err := runtime.Client(ctx, viper.GetString("runtime"))
		cobra.CheckErr(err)

		r, err := cli.ContainerCreate(ctx, &container.Config{
			Hostname:     "cage",
			Domainname:   "cage",
			User:         "1000",
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			OpenStdin:    true,
			Image:        "docker.io/library/ubuntu:latest",
			Entrypoint:   []string{"bash"},
		}, nil, nil, nil, "cage")
		cobra.CheckErr(err)

		fmt.Println(r.ID)
		for _, warning := range r.Warnings {
			fmt.Println(warning)
		}

		// Start the container
		err = cli.ContainerStart(ctx, r.ID, container.StartOptions{})
		cobra.CheckErr(err)

		// Attach to the container
		attachResp, err := cli.ContainerAttach(ctx, r.ID, container.AttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		})
		cobra.CheckErr(err)
		defer attachResp.Close()

		// Set terminal to raw mode for proper TTY interaction
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		// Connect stdin/stdout - only ONE copy from Reader to Stdout
		go io.Copy(os.Stdout, attachResp.Reader)
		go io.Copy(attachResp.Conn, os.Stdin)

		// Wait for container to finish OR interrupt signal
		statusCh, errCh := cli.ContainerWait(ctx, r.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error waiting for container: %v\n", err)
			}
		case <-statusCh:
			// Container exited normally
		case <-ctx.Done():
			// Interrupted by user (Ctrl+C)
			fmt.Println("\nDetaching...")
		}

		cli.ContainerStop(ctx, r.ID, container.StopOptions{})
		cli.ContainerRemove(ctx, r.ID, container.RemoveOptions{})
	},
}

func RequireInitialized(cmd *cobra.Command, args []string) {
	if !state.IsInitialized() {
		cobra.CheckErr(fmt.Errorf("cage is not initialized: run `%v init` first", config.AppName))
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/.cage.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Init(cfgFile)
}
