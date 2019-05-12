package recipecmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ValidateArgs(_ *cobra.Command, args []string) error {
	if (len(args) == 0 || len(args) > 1 ) && URL == "" && FilePath == "" {
		return errors.New(`This command requires single recipe name from the official repository.
You can also use additional flags to load recipe from disk or URL.
`)
	}

	return nil
}
