package cmd

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestRunRollback(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		replaceOSLineInRecipe(t, ValidRecipePath, RecipeOSPlaceholder, runtime.GOOS)

		err := runRollback(nil, []string{ValidRecipePath})

		replaceOSLineInRecipe(t, ValidRecipePath, runtime.GOOS, RecipeOSPlaceholder)

		assert.NoError(t, err)
	})

	t.Run("Invalid Path", func(t *testing.T) {
		err := runRollback(nil, []string{"./testdata/file.yaml"})

		assert.Error(t, err)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		err := runRollback(nil, []string{"https://example.com/foo/bar"})

		assert.Error(t, err)
	})

	t.Run("Invalid Recipe", func(t *testing.T) {
		err := runRollback(nil, []string{InvalidRecipePath})

		assert.Error(t, err)
	})

}
