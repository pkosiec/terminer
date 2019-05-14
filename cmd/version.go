package cmd

import (
	"github.com/pkosiec/terminer/internal/metadata"
	"github.com/pkosiec/terminer/internal/printer"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the application version",
	Run:   printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(_ *cobra.Command, _ []string) {
	printer.New().AppInfo(metadata.AppName, metadata.Version, metadata.URL)
}
