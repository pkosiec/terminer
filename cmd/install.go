package cmd

import (
	"github.com/pkosiec/terminer/internal/printer"
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
	Args: validateInstallRollbackArgs,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	p := printer.New()
	i, err := setupInstaller(args[0], p)
	if err != nil {
		return err
	}

	err = i.Install()
	p.Result(err)

	return nil
}
