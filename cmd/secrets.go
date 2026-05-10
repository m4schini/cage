/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cage/agent/secrets"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret"},
	Short:   "Manage secrets (mainly ai api keys) for your cages",
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
	Use:   "add LABEL",
	Short: "add secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Secret: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // move to next line after user presses Enter
		cobra.CheckErr(err)

		cobra.CheckErr(secrets.Store(args[0], string(bytePassword)))
	},
}

var secretsGetCmd = &cobra.Command{
	Use:   "get LABEL",
	Short: "get secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		unredacted, err := cmd.Flags().GetBool("unredacted")
		cobra.CheckErr(err)

		s, err := secrets.Retrieve(args[0])
		cobra.CheckErr(err)

		if !unredacted {
			s = redact(s)
		}

		fmt.Println(s)
	},
}

func redact(s string) string {
	runes := []rune(s) // handle Unicode correctly
	n := len(runes)

	if n <= 4 {
		return s
	}

	// Create a slice of '*' for all but the last 4 characters
	redacted := make([]rune, n)
	for i := 0; i < n-4; i++ {
		redacted[i] = '*'
	}
	copy(redacted[n-4:], runes[n-4:])

	return string(redacted)
}

var secretsDeleteCmd = &cobra.Command{
	Use:     "delete LABEL",
	Aliases: []string{"rm"},
	Short:   "delete secret",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(secrets.Delete(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsAddCmd)
	secretsCmd.AddCommand(secretsGetCmd)
	secretsCmd.AddCommand(secretsDeleteCmd)

	secretsGetCmd.Flags().Bool("unredacted", false, "show secret unredacted")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
