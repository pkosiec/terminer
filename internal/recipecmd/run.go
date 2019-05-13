package recipecmd

import (
	"github.com/pkosiec/terminer/internal/printer"
	"github.com/pkosiec/terminer/pkg/installer"
	"github.com/pkosiec/terminer/pkg/recipe"
	"github.com/spf13/cobra"
	"net/http"
)

type RunType string

var (
	Install  RunType = "Install"
	Rollback RunType = "Rollback"
)

func Run(runType RunType) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		p := printer.New()
		i, err := loadRecipeAndSetupInstaller(args, URL, FilePath, p)
		if err != nil {
			return err
		}

		err = func() error {
			switch runType {
			case Install:
				return i.Install()
			case Rollback:
				return i.Rollback()
			}

			return i.Install()
		}()
		p.Result(err)

		return nil
	}
}

func loadRecipeAndSetupInstaller(recipeNames []string, URL, filePath string, p printer.Printer) (*installer.Installer, error) {
	var r *recipe.Recipe
	var err error

	if len(recipeNames) > 0 && recipeNames[0] != "" {
		r, err = recipe.FromRepository(recipeNames[0], http.DefaultClient)
	} else if URL != "" {
		r, _, err = recipe.FromURL(URL, http.DefaultClient)
	} else if filePath != "" {
		r, err = recipe.FromPath(filePath)
	}

	if err != nil {
		return nil, err
	}

	i, err := installer.New(r, p)
	if err != nil {
		return nil, err
	}

	return i, nil
}
