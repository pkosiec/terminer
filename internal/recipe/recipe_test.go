package recipe_test

import (
	"github.com/pkosiec/terminer/internal/recipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestFromPath(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := fixRecipe("testos")

		r, err := recipe.FromPath("./fixture/valid-recipe.yaml")

		require.NoError(t, err)
		assert.Equal(t, expected, r)
	})

	t.Run("Invalid Path", func(t *testing.T) {
		_, err := recipe.FromPath("./fixture/no-file-exist.yaml")

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file")
	})

	t.Run("Invalid File", func(t *testing.T) {
		_, err := recipe.FromPath("./fixture/invalid-recipe.yaml")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "while loading recipe")
	})
}

func TestFromURL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := fixRecipe("testos")
		server := setupRemoteRecipeServer(t, "./fixture/valid-recipe.yaml", false)
		defer server.Close()

		r, err := recipe.FromURL(server.URL)

		require.NoError(t, err)
		assert.Equal(t, expected, r)
	})

	t.Run("Not existing path", func(t *testing.T) {
		_, err := recipe.FromURL("http://foo-bar.not-existing.url")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "while requesting")
	})

	t.Run("Server Error", func(t *testing.T) {
		server := setupRemoteRecipeServer(t, "", true)
		defer server.Close()

		_, err := recipe.FromURL(server.URL)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid status code")
	})

	t.Run("Error during reading response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}))
		defer server.Close()

		_, err := recipe.FromURL(server.URL)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "while reading response body")
	})

	t.Run("Empty response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		_, err := recipe.FromURL(server.URL)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Empty body")
	})

	t.Run("Invalid File", func(t *testing.T) {
		server := setupRemoteRecipeServer(t, "./fixture/invalid-recipe.yaml", false)
		defer server.Close()

		_, err := recipe.FromURL(server.URL)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "while loading recipe from URL")
	})
}

func TestValidate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		err := r.Validate()

		assert.NoError(t, err)
	})

	t.Run("Invalid OS", func(t *testing.T) {
		r := fixRecipe("notexistingos")

		err := r.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid operating system")
	})

	t.Run("No stages", func(t *testing.T) {
		r := &recipe.Recipe{
			OS:   runtime.GOOS,
			Name: "Test",
		}

		err := r.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "No stages")
	})

	t.Run("No steps in stage", func(t *testing.T) {
		r := &recipe.Recipe{
			OS:   runtime.GOOS,
			Name: "Test",
			Stages: []recipe.Stage{
				{
					Name: "1",
					Steps: []recipe.Step{
						{
							Name:    "Test",
							Command: "echo \"test\"",
						},
					},
				},
				{
					Name: "2",
				},
			},
		}

		err := r.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "No steps")
	})

	t.Run("No commands in stage", func(t *testing.T) {
		r := &recipe.Recipe{
			OS:   runtime.GOOS,
			Name: "Test",
			Stages: []recipe.Stage{
				{
					Name: "1",
					Steps: []recipe.Step{
						{
							Name: "Test",
						},
					},
				},
			},
		}

		err := r.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "No command")
	})
}

func fixRecipe(os string) *recipe.Recipe {
	return &recipe.Recipe{
		OS:          os,
		Name:        "Recipe",
		Description: "Recipe Description",
		Stages: []recipe.Stage{
			{
				Name:        "Stage 1",
				Description: "Stage 1 description",
				ReadMoreURL: "https://stage1.example.com",
				Steps: []recipe.Step{
					{
						Name:        "Step 1",
						ReadMoreURL: "https://step1.stage1.example.com",
						Command:     "echo \"Step 1 of Stage 1\"",
						Rollback:    "echo \"Rollback of Step 1 of Stage 1\"",
					},
					{
						Name:        "Step 2",
						ReadMoreURL: "https://step2.stage1.example.com",
						Command:     "echo \"Step 2 of Stage 1\"",
						Rollback:    "echo \"Rollback of Step 2 of Stage 1\"",
					},
				},
			},
			{
				Name:        "Stage 2",
				Description: "Stage 2 description",
				ReadMoreURL: "https://stage2.example.com",
				Steps: []recipe.Step{
					{
						Name:        "Step 2",
						ReadMoreURL: "https://step2.stage2.example.com",
						Command:     "echo \"Step 1 of Stage 2\"",
						Rollback:    "echo \"Rollback of Step 1 of Stage 2\"",
					},
					{
						Name:        "Step 2",
						ReadMoreURL: "https://step2.stage2.example.com",
						Command:     "echo \"Step 2 of Stage 2\"",
						Rollback:    "echo \"Rollback of Step 2 of Stage 2\"",
					},
				},
			},
		},
	}
}

func setupRemoteRecipeServer(t *testing.T, recipePath string, returnError bool) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if returnError {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Internal Server Error"))
			require.NoError(t, err)
			return
		}

		yamlFile, err := ioutil.ReadFile(recipePath)
		require.NoError(t, err)
		_, err = w.Write(yamlFile)
		require.NoError(t, err)
	}))

	return server
}

