/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cage/agent/secrets"
	"fmt"

	"github.com/spf13/cobra"
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var secretsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list secrets",
	Run: func(cmd *cobra.Command, args []string) {
		ls, err := secrets.List()
		cobra.CheckErr(err)

		for i, l := range ls {
			fmt.Println(i, ":", l)
		}
	},
}

var secretsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add secret",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("secrets called")
	},
}

var secretsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get secret",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("secrets called")
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsAddCmd)
	secretsCmd.AddCommand(secretsGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
