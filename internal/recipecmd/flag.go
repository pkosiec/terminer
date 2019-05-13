package recipecmd

import "github.com/spf13/cobra"

var URL string
var FilePath string

func SupportFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&URL, "url", "u", "", "Recipe URL")
	cmd.Flags().StringVarP(&FilePath, "filepath", "f", "", "Recipe file path")
}
