package recipe_test

import (
	"github.com/pkosiec/terminer/pkg/recipe"
	"github.com/pkosiec/terminer/pkg/shell"
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

		r, err := recipe.FromPath("./testdata/valid-recipe.yaml")

		require.NoError(t, err)
		assert.Equal(t, expected, r)
	})

	t.Run("Invalid Path", func(t *testing.T) {
		_, err := recipe.FromPath("./testdata/no-file-exist.yaml")

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file")
	})

	t.Run("Invalid extension", func(t *testing.T) {
		_, err := recipe.FromPath("./testdata/valid-recipe.sh")

		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid file extension")
	})

	t.Run("Invalid File", func(t *testing.T) {
		_, err := recipe.FromPath("./testdata/invalid-recipe.yaml")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "while loading recipe")
	})
}

func TestFromURL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := fixRecipe("testos")
		server := setupRemoteRecipeServer(t, "./testdata/valid-recipe.yaml", false)
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
		server := setupRemoteRecipeServer(t, "./testdata/invalid-recipe.yaml", false)
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
			OS: runtime.GOOS,
			Metadata: recipe.UnitMetadata{
				Name: "Test",
			},
		}

		err := r.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "No stages")
	})

	t.Run("No steps in stage", func(t *testing.T) {
		r := &recipe.Recipe{
			OS: runtime.GOOS,
			Metadata: recipe.UnitMetadata{
				Name: "Test",
			},
			Stages: []recipe.Stage{
				{
					Metadata: recipe.UnitMetadata{
						Name: "1",
					},
					Steps: []recipe.Step{
						{
							Metadata: recipe.UnitMetadata{
								Name: "Test",
							},
							Execute: shell.Command{
								Run: "echo \"test\"",
							},
						},
					},
				},
				{
					Metadata: recipe.UnitMetadata{
						Name: "2",
					},
				},
			},
		}

		err := r.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "No steps")
	})

	t.Run("No commands in stage", func(t *testing.T) {
		r := &recipe.Recipe{
			OS: runtime.GOOS,
			Metadata: recipe.UnitMetadata{
				Name: "Test",
			},
			Stages: []recipe.Stage{
				{
					Metadata: recipe.UnitMetadata{
						Name: "1",
					},
					Steps: []recipe.Step{
						{
							Metadata: recipe.UnitMetadata{
								Name: "Test",
							},
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
		OS: os,
		Metadata: recipe.UnitMetadata{
			Name:        "Recipe",
			Description: "Recipe Description",
		},
		Stages: []recipe.Stage{
			{
				Metadata: recipe.UnitMetadata{
					Name:        "Stage 1",
					Description: "Stage 1 description",
					URL:         "https://stage1.example.com",
				},
				Steps: []recipe.Step{
					{
						Metadata: recipe.UnitMetadata{
							Name: "Step 1",
							URL:  "https://step1.stage1.example.com",
						},
						Execute: shell.Command{
							Run: "echo \"Step 1 of Stage 1\"",
						},
						Rollback: shell.Command{
							Run: "echo \"Rollback of Step 1 of Stage 1\"",
						},
					},
					{
						Metadata: recipe.UnitMetadata{
							Name: "Step 2",
							URL:  "https://step2.stage1.example.com",
						},
						Execute: shell.Command{
							Run: "echo \"Step 2 of Stage 1\"",
						},
						Rollback: shell.Command{
							Run: "echo \"Rollback of Step 2 of Stage 1\"",
						},
					},
				},
			},
			{
				Metadata: recipe.UnitMetadata{
					Name:        "Stage 2",
					Description: "Stage 2 description",
					URL:         "https://stage2.example.com",
				},
				Steps: []recipe.Step{
					{
						Metadata: recipe.UnitMetadata{
							Name: "Step 1",
							URL:  "https://step1.stage2.example.com",
						},
						Execute: shell.Command{
							Shell: "sh",
							Run:   "echo \"Step 1 of Stage 2\"",
						},
						Rollback: shell.Command{
							Run: "echo \"Rollback of Step 1 of Stage 2\"",
						},
					},
					{
						Metadata: recipe.UnitMetadata{
							Name: "Step 2",
							URL:  "https://step2.stage2.example.com",
						},
						Execute: shell.Command{
							Run: "echo \"Step 2 of Stage 2\"",
						},
						Rollback: shell.Command{
							Run: "echo \"Rollback of Step 2 of Stage 2\"",
						},
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