package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	date    = "unknown"
	commit  = "none"
	version = "dev"

	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("az-dns version", version)
			fmt.Println("  built at", date)
			fmt.Println("  from commit", commit)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
