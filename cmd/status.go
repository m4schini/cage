package cmd

import (
	"cage/container/runtime"
	"cage/state"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status summary",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if !state.IsInitialized() {
			fmt.Println("cage is not initialized")
			return
		}
		fmt.Println("Using config:", viper.ConfigFileUsed())
		fmt.Println()

		fmt.Println("Runtimes:")
		runtimes := runtime.Available(ctx)
		for runtime, socket := range runtimes {
			var status string
			switch {
			case socket == nil:
				status = "not found"
				break
			case !socket.Available:
				status = fmt.Sprintf("unavailable (%v)", socket.Host)
				break
			default:
				status = "available"
			}

			fmt.Printf("%v\t%v\n", runtime, status)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
