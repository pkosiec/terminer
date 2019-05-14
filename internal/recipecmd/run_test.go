package recipecmd_test

import (
	"github.com/pkosiec/terminer/internal/recipecmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const ValidRecipePath = "./testdata/valid-recipe.yaml"
const InvalidRecipePath = "./testdata/invalid-recipe.yaml"
const EmptyRecipePath = "./testdata/empty-recipe.yaml"
const FailingRecipePath = "./testdata/failing-recipe.yaml"

func TestRun(t *testing.T) {
	filePathBak := recipecmd.FilePath
	urlBak := recipecmd.URL

	t.Run("Install", func(t *testing.T) {
		installFn := recipecmd.Run(recipecmd.Install)

		t.Run("Valid recipe from path", func(t *testing.T) {
			recipecmd.FilePath = ValidRecipePath
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			assert.NoError(t, err)
		})

		t.Run("Valid recipe from URL", func(t *testing.T) {
			server := setupRemoteRecipeServer(t, "./testdata/valid-recipe.yaml")
			defer server.Close()

			recipecmd.FilePath = ""
			recipecmd.URL = server.URL
			err := installFn(nil, []string{})

			require.NoError(t, err)
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

			// should not exit with error. It should print error instead
			// TODO: Test printing error
			assert.NoError(t, err)
		})

		t.Run("Empty Recipe", func(t *testing.T) {
			path := EmptyRecipePath

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := installFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Too many parameters", func(t *testing.T) {
			recipecmd.FilePath = ""
			recipecmd.URL = ""
			err := installFn(nil, []string{"test", "test2"})

			assert.Error(t, err)
		})
	})

	t.Run("Rollback", func(t *testing.T) {
		rollbackFn := recipecmd.Run(recipecmd.Rollback)

		t.Run("Valid recipe", func(t *testing.T) {
			recipecmd.FilePath = ValidRecipePath
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

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

			// should not exit with error. It should print error instead
			// TODO: Test printing error
			assert.NoError(t, err)
		})

		t.Run("Empty Recipe", func(t *testing.T) {
			path := EmptyRecipePath

			recipecmd.FilePath = path
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{})

			assert.Error(t, err)
		})

		t.Run("Too many parameters", func(t *testing.T) {
			recipecmd.FilePath = ""
			recipecmd.URL = ""
			err := rollbackFn(nil, []string{"test", "test2"})

			assert.Error(t, err)
		})
	})

	recipecmd.FilePath = filePathBak
	recipecmd.URL = urlBak
}

func setupRemoteRecipeServer(t *testing.T, recipePath string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		yamlFile, err := ioutil.ReadFile(recipePath)
		require.NoError(t, err)
		_, err = w.Write(yamlFile)
		require.NoError(t, err)
	}))

	return server
}
