package cmd

import (
	"github.com/pkosiec/terminer/internal/recipecmd"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [recipe name]",
	Short: "Installs a recipe from official repository, given path or URL",
	Long: `Install command installs a recipe from the official recipe repository.
You can use additional flags to install a recipe from a local or remote file.

Examples:
	terminer install zsh-starter
	terminer install -f ./recipe.yaml
	terminer install --file /Users/sample-user/recipe.yml
	terminer install -u https://example.com/recipe.yaml
	terminer install --url http://foo.bar/recipe.yml
`,
	Args:                  recipecmd.ValidateArgs,
	RunE:                  recipecmd.Run(recipecmd.Install),
	DisableFlagsInUseLine: true,
}

func init() {
	recipecmd.SupportFlags(installCmd)
	rootCmd.AddCommand(installCmd)
}
