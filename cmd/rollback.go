package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback [file path or URL]",
	Short: "Rollbacks a recipe from given path or URL",
	Long: `Rollback command rollbacks a recipe from a local or remote file.
Provide a relative or absolute path to a YAML file with recipe
or an URL to download it.

Examples:
	terminer rollback ./recipe.yaml
	terminer rollback /Users/sample-user/recipe.yaml
	terminer rollback https://example.com/recipe.yaml
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

		return i.Rollback()
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
