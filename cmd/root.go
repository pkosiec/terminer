package cmd

import (
	"fmt"
	"github.com/pkosiec/terminer/internal/installer"
	"github.com/pkosiec/terminer/internal/path"
	"github.com/pkosiec/terminer/internal/recipe"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "terminer",
	Short: "Upgrade your terminal experience",
	Long: `Terminer is an cross-platform installer for terminal presets.
Install Fish or ZSH shell packed with useful plugins and
sleek prompts. Use one of starter recipes or make yours.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SetupInstaller(filePath string) (*installer.Installer, error) {
	var r *recipe.Recipe
	var err error

	if path.IsURL(filePath) {
		r, err = recipe.FromURL(filePath)
	} else {
		r, err = recipe.FromPath(filePath)
	}

	if err != nil {
		return nil, err
	}

	i, err := installer.New(r)
	if err != nil {
		return nil, err
	}

	return i, nil
}
