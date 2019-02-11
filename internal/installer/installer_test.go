package installer_test

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/installer"
	"github.com/pkosiec/terminer/internal/recipe"
	"github.com/pkosiec/terminer/internal/shell/automock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Empty recipe", func(t *testing.T) {
		var r *recipe.Recipe

		_, err := installer.New(r)

		require.Error(t, err)
	})

	t.Run("Invalid recipe", func(t *testing.T) {
		r := fixRecipe("testos")

		_, err := installer.New(r)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid operating system")
	})

	t.Run("Valid recipe", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		_, err := installer.New(r)

		require.NoError(t, err)
	})
}

func TestInstaller_Install(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Shell{}
		shImpl.On("Exec", "sh", "echo \"C1/1\"").Return("", nil).Once()
		shImpl.On("Exec", "sh", "echo \"C2/1\"").Return("", nil).Once()
		shImpl.On("Exec", "sh", "echo \"C1/2\"").Return("", nil).Once()
		shImpl.On("Exec", "sh", "echo \"C2/2\"").Return("", nil).Once()
		defer shImpl.AssertExpectations(t)
		i.SetShell(shImpl)

		err = i.Install()
		require.NoError(t, err)
	})

	// Should exit after failed step
	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("Test Err")
		r := fixRecipe(runtime.GOOS)

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Shell{}
		shImpl.On("Exec", "sh", "echo \"C1/1\"").Return("", testErr).Once()
		defer shImpl.AssertExpectations(t)
		i.SetShell(shImpl)

		err = i.Install()
		require.Error(t, err)
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestInstaller_Rollback(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Shell{}
		shImpl.On("Exec", "sh","echo \"R2/2\"").Return("", nil).Once()
		shImpl.On("Exec", "sh","echo \"R1/2\"").Return("", nil).Once()
		shImpl.On("Exec", "sh","echo \"R2/1\"").Return("", nil).Once()
		shImpl.On("Exec", "sh","echo \"R1/1\"").Return("", nil).Once()
		defer shImpl.AssertExpectations(t)
		i.SetShell(shImpl)

		err = i.Rollback()
		require.NoError(t, err)
	})

	// Should not exit after failed step
	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("Test Err")
		r := fixRecipe(runtime.GOOS)

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Shell{}
		shImpl.On("Exec", "sh", "echo \"R2/2\"").Return("", testErr).Once()
		shImpl.On("Exec", "sh", "echo \"R1/2\"").Return("", testErr).Once()
		shImpl.On("Exec", "sh","echo \"R2/1\"").Return("", nil).Once()
		shImpl.On("Exec", "sh","echo \"R1/1\"").Return("", nil).Once()
		defer shImpl.AssertExpectations(t)
		i.SetShell(shImpl)

		err = i.Rollback()
		require.NoError(t, err)
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
						Command:     "echo \"C1/1\"",
						Rollback:    "echo \"R1/1\"",
					},
					{
						Name:        "Step 2",
						ReadMoreURL: "https://step2.stage1.example.com",
						Command:     "echo \"C2/1\"",
						Rollback:    "echo \"R2/1\"",
					},
				},
			},
			{
				Name:        "Stage 2",
				Description: "Stage 2 description",
				ReadMoreURL: "https://stage2.example.com",
				Steps: []recipe.Step{
					{
						Name:        "Step 1",
						ReadMoreURL: "https://step2.stage2.example.com",
						Command:     "echo \"C1/2\"",
						Rollback:    "echo \"R1/2\"",
					},
					{
						Name:        "Step 2",
						ReadMoreURL: "https://step2.stage2.example.com",
						Command:     "echo \"C2/2\"",
						Rollback:    "echo \"R2/2\"",
					},
				},
			},
		},
	}
}
