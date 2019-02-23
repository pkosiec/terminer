package cmd

import (
	"github.com/pkosiec/terminer/internal/printer"
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
	Args: validateInstallRollbackArgs,
	RunE: runRollback,
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, args []string) error {
	p := printer.New()
	i, err := setupInstaller(args[0], p)
	if err != nil {
		return err
	}

	err = i.Rollback()
	p.Result(err)

	return nil
}
