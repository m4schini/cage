package cmd

import (
	"cage/cage"
	config2 "cage/cage/config"
	"cage/cage/state"
	ctr "cage/container"
	"cage/container/runtime"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:    config2.AppName,
	Short:  `isolate "ai" "agents"`,
	PreRun: RequireInitialized,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
		defer cancel()
		cli, err := runtime.Client(ctx, viper.GetString("runtime"))
		cobra.CheckErr(err)

		d := ctr.Docker{Client: cli}
		err = cage.Run(ctx, "ais", &d)
		cobra.CheckErr(err)
	},
}

func RequireInitialized(cmd *cobra.Command, args []string) {
	if !state.IsInitialized() {
		cobra.CheckErr(fmt.Errorf("cage is not initialized: run `%v init` first", config2.AppName))
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
	config2.Init(cfgFile)
}
