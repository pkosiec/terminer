package cmd

import (
	"github.com/spf13/cobra"
)

func PrintVersion(_ *cobra.Command, _ []string) {
	printVersion(nil, nil)
}
