package cmd

import (
	"cage/cage"
	config2 "cage/cage/config"
	"cage/cage/state"
	"cage/container/runtime"
	"cage/nix"
	"time"

	"fmt"
	"os"
	"os/signal"

	_ "embed"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              config2.AppName,
	Short:            `isolate "ai" "agents"`,
	PersistentPreRun: InitCageApp,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
		defer cancel()

		if dryrun, err := cmd.Flags().GetBool("dry-run"); err == nil && dryrun {
			cfg, lastModified, err := state.Load()
			cobra.CheckErr(err)

			fmt.Println("Image:", cfg.ImageName())
			fmt.Println("Last Modified:", lastModified.Format(time.RFC822))
			fmt.Println()
			fmt.Println("ENV:")
			env, err := cfg.ClaudeCodeEnv(true)
			cobra.CheckErr(err)

			for _, s := range env {
				fmt.Println(s)
			}

			nix, err := nix.NewNixShellString(nix.ShellNixPackages{
				Packages: cfg.Packages,
				Shell:    "bash",
			})
			cobra.CheckErr(err)

			fmt.Println()
			fmt.Println("shell.nix:")
			fmt.Println(nix)

			return
		}

		rt := viper.GetString("runtime")
		fmt.Println("runtime:", rt)
		cli, err := runtime.Client(ctx, rt)
		cobra.CheckErr(err)

		err = cage.Run(ctx, cli)
		cobra.CheckErr(err)
	},
}

func InitCageApp(cmd *cobra.Command, args []string) {

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
	rootCmd.Flags().Bool("dry-run", false, "Dry run, only show build/run config")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config2.Init(cfgFile)
}
