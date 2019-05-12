package recipecmd_test

import (
	"fmt"
	"github.com/pkosiec/terminer/internal/recipecmd"
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

func TestRun(t *testing.T) {
	filePathBak := recipecmd.FilePath
	urlBak := recipecmd.URL

	t.Run("Install", func(t *testing.T) {
		installFn := recipecmd.Run(recipecmd.Install)

		t.Run("Valid recipe from repository", func(t *testing.T) {
			//TODO: Test it properly
			//recipecmd.FilePath = ""
			//recipecmd.URL = ""
			//
			//err := installFn(nil, []string{"fish-starter"})
			//
			//assert.NoError(t, err)
		})

		t.Run("Valid recipe from URL", func(t *testing.T) {
			//TODO: Test it properly - replace OS
			//url := fmt.Sprintf(
			//	"https://raw.githubusercontent.com/%s/%s/%s/pkg/recipe/testdata/valid-recipe.yaml",
			//	metadata.Repository.Owner,
			//	metadata.Repository.Name,
			//	metadata.Repository.BranchName,
			//)
			//
			//recipecmd.FilePath = ""
			//recipecmd.URL = url
			//
			//err := installFn(nil, []string{})
			//
			//assert.NoError(t, err)
		})

		t.Run("Valid recipe from path", func(t *testing.T) {
			replaceOSLineInRecipe(t, ValidRecipePath, RecipeOSPlaceholder, runtime.GOOS)

			recipecmd.FilePath = ValidRecipePath
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			replaceOSLineInRecipe(t, ValidRecipePath, runtime.GOOS, RecipeOSPlaceholder)

			assert.NoError(t, err)
		})

		t.Run("Invalid Recipe from path", func(t *testing.T) {
			recipecmd.FilePath = InvalidRecipePath
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Invalid path", func(t *testing.T) {
			path := "./testdata/file.yaml"

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Invalid URL", func(t *testing.T) {
			url := "https://example.com/foo/bar"

			recipecmd.FilePath = ""
			recipecmd.URL = url
			err := installFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Failing Recipe", func(t *testing.T) {
			path := FailingRecipePath

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Empty Recipe", func(t *testing.T) {
			path := EmptyRecipePath

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			assert.Error(t, err)
		})
	})

	t.Run("Rollback", func(t *testing.T) {
		rollbackFn := recipecmd.Run(recipecmd.Rollback)

		t.Run("Valid recipe", func(t *testing.T) {
			replaceOSLineInRecipe(t, ValidRecipePath, RecipeOSPlaceholder, runtime.GOOS)

			recipecmd.FilePath = ValidRecipePath
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

			replaceOSLineInRecipe(t, ValidRecipePath, runtime.GOOS, RecipeOSPlaceholder)

			assert.NoError(t, err)
		})

		t.Run("Invalid Recipe", func(t *testing.T) {
			recipecmd.FilePath = InvalidRecipePath
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Invalid Path", func(t *testing.T) {
			path := "./testdata/file.yaml"

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Invalid URL", func(t *testing.T) {
			url := "https://example.com/foo/bar"

			recipecmd.FilePath = ""
			recipecmd.URL = url
			err := rollbackFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Failing Recipe", func(t *testing.T) {
			path := FailingRecipePath

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Empty Recipe", func(t *testing.T) {
			path := EmptyRecipePath

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

			assert.Error(t, err)
		})
	})

	recipecmd.FilePath = filePathBak
	recipecmd.URL = urlBak
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