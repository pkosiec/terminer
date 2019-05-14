package recipecmd

import "github.com/spf13/cobra"

// URL is a variable which stores an address to a recipe given by user
var URL string

// FilePath is a variable which stores a file path to a recipe given by user
var FilePath string

// SupportFlags sets required flags for recipe operations
func SupportFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&URL, "url", "u", "", "Recipe URL")
	cmd.Flags().StringVarP(&FilePath, "filepath", "f", "", "Recipe file path")
}
