/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cage/cage"
	"cage/cage/state"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:    "new NAME",
	Short:  "Create new cage",
	Args:   cobra.ExactArgs(1),
	PreRun: RequireInitialized,
	Run: func(cmd *cobra.Command, args []string) {
		cageName := args[0]

		err := cage.New(cageName, state.CageDefinition{
			Shell: "zsh",
			Packages: []string{
				"go",
				"claude-code",
			},
			Env: nil,
		})
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
