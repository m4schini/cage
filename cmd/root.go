package cmd

import (
	"cage/config"
	"cage/container"
	"cage/state"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:    config.AppName,
	Short:  `isolate "ai" "agents"`,
	PreRun: RequireInitialized,
	Run: func(cmd *cobra.Command, args []string) {
		err := container.Run(cmd.Context())
		cobra.CheckErr(err)
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
