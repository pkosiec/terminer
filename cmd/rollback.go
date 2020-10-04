package cmd

import (
	"github.com/pkosiec/terminer/internal/recipecmd"
	"github.com/pkosiec/terminer/pkg/shared"
	"github.com/spf13/cobra"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback [recipe name]",
	Short: "Rollbacks a recipe from official repository, given path or URL",
	Long: `Rollback command uninstalls a recipe from the official recipe repository.
You can use additional flags to rollback a recipe from a local or remote file.`,
Example:
`	terminer rollback zsh-starter
	terminer rollback -f ./recipe.yaml
	terminer rollback --file /Users/sample-user/recipe.yml
	terminer rollback -u https://example.com/recipe.yaml
	terminer rollback --url http://foo.bar/recipe.yml
`,
	Args:                  recipecmd.ValidateArgs,
	RunE:                  recipecmd.Run(shared.OperationRollback),
	DisableFlagsInUseLine: true,
}

func init() {
	recipecmd.SupportFlags(rollbackCmd)
	rootCmd.AddCommand(rollbackCmd)
}
