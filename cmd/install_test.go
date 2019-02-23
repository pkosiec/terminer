package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"
)

const ValidRecipePath = "./testdata/valid-recipe.yaml"
const InvalidRecipePath = "./testdata/invalid-recipe.yaml"
const EmptyRecipePath = "./testdata/empty-recipe.yaml"
const FailingRecipePath = "./testdata/failing-recipe.yaml"
const RecipeOSPlaceholder = "{CURRENT_OS_PLACEHOLDER}"

func TestRunInstall(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		replaceOSLineInRecipe(t, ValidRecipePath, RecipeOSPlaceholder, runtime.GOOS)

		err := runInstall(nil, []string{ValidRecipePath})

		replaceOSLineInRecipe(t, ValidRecipePath, runtime.GOOS, RecipeOSPlaceholder)

		assert.NoError(t, err)
	})

	t.Run("Invalid Path", func(t *testing.T) {
		err := runInstall(nil, []string{"./testdata/file.yaml"})

		assert.Error(t, err)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		err := runInstall(nil, []string{"https://example.com/foo/bar"})

		assert.Error(t, err)
	})

	t.Run("Invalid Recipe", func(t *testing.T) {
		err := runInstall(nil, []string{InvalidRecipePath})

		assert.Error(t, err)
	})

	t.Run("Failing Recipe", func(t *testing.T) {
		err := runInstall(nil, []string{FailingRecipePath})

		assert.Error(t, err)
	})

	t.Run("Empty Recipe", func(t *testing.T) {
		err := runInstall(nil, []string{EmptyRecipePath})

		assert.Error(t, err)
	})
}

func replaceOSLineInRecipe(t *testing.T, path, from, to string) {
	input, err := ioutil.ReadFile(path)
	require.NoError(t, err)

	output := strings.Replace(string(input), fmt.Sprintf("os: %s", from), fmt.Sprintf("os: %s", to), 1)
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
