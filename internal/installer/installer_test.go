package installer_test

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/installer"
	"github.com/pkosiec/terminer/internal/recipe"
	"github.com/pkosiec/terminer/internal/sh/automock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func TestInstaller_Install(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := fixRecipe()

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Sh{}
		shImpl.On("Exec", "echo \"C1/1\"").Return("", nil).Once()
		shImpl.On("Exec", "echo \"C2/1\"").Return("", nil).Once()
		shImpl.On("Exec", "echo \"C1/2\"").Return("", nil).Once()
		shImpl.On("Exec", "echo \"C2/2\"").Return("", nil).Once()
		defer shImpl.AssertExpectations(t)
		i.SetSh(shImpl)

		err = i.Install()
		require.NoError(t, err)
	})

	// Should exit after failed step
	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("Test Err")
		r := fixRecipe()

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Sh{}
		shImpl.On("Exec", "echo \"C1/1\"").Return("", testErr).Once()
		defer shImpl.AssertExpectations(t)
		i.SetSh(shImpl)

		err = i.Install()
		require.Error(t, err)
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestInstaller_Rollback(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := fixRecipe()

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Sh{}
		shImpl.On("Exec", "echo \"R2/2\"").Return("", nil).Once()
		shImpl.On("Exec", "echo \"R1/2\"").Return("", nil).Once()
		shImpl.On("Exec", "echo \"R2/1\"").Return("", nil).Once()
		shImpl.On("Exec", "echo \"R1/1\"").Return("", nil).Once()
		defer shImpl.AssertExpectations(t)
		i.SetSh(shImpl)

		err = i.Rollback()
		require.NoError(t, err)
	})

	// Should exit after failed step
	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("Test Err")
		r := fixRecipe()

		i, err := installer.New(r)
		require.NoError(t, err)

		shImpl := &automock.Sh{}
		shImpl.On("Exec", "echo \"R2/2\"").Return("", testErr).Once()
		defer shImpl.AssertExpectations(t)
		i.SetSh(shImpl)

		err = i.Rollback()
		require.Error(t, err)
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func fixRecipe() *recipe.Recipe {
	return &recipe.Recipe{
		OS:          runtime.GOOS,
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
