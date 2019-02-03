package recipe_test

import (
	"github.com/pkosiec/terminer/internal/recipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func TestFrom(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r, err := recipe.From("./fixture/valid-recipe.yaml")
		expected := fixRecipe("testos")
		require.NoError(t, err)
		assert.Equal(t, expected, r)
	})

	t.Run("Invalid Path", func(t *testing.T) {
		_, err := recipe.From("./fixture/no-file-exist.yaml")
		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file")
	})

	t.Run("Invalid File", func(t *testing.T) {
		_, err := recipe.From("./fixture/invalid-recipe.yaml")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "while loading recipe")
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
