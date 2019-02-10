package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [file path or URL]",
	Short: "Installs a recipe from given path or URL",
	Long: `Install command installs a recipe from a local or remote file.
Provide a relative or absolute path to a YAML file with recipe
or an URL to download it.

Examples:
	terminer install ./recipe.yaml
	terminer install /Users/sample-user/recipe.yaml
	terminer install https://example.com/recipe.yaml
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires one argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		i, err := SetupInstaller(args[0])
		if err != nil {
			return err
		}

		return i.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
