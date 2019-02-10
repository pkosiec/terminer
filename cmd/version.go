package cmd

import (
	"fmt"
	"github.com/pkosiec/terminer/internal/metadata"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v. %s\n%s", metadata.AppName, metadata.Version, metadata.URL)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
