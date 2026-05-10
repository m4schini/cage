/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the global cage configuration",
}

var configListCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"ls", "list"},
	Short:   "Show configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Using config:", viper.ConfigFileUsed())
		fmt.Println()

		for _, s := range viper.AllKeys() {
			fmt.Printf("%v=%v\n", s, viper.Get(s))
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set config value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set(args[0], args[1])
		cobra.CheckErr(viper.WriteConfig())
	},
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete KEY",
	Aliases: []string{"reset"},
	Short:   "Revert config value to default",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set(args[0], nil)
		cobra.CheckErr(viper.WriteConfig())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
