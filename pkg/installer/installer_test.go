package installer_test

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/pkg/installer"
	"github.com/pkosiec/terminer/pkg/recipe"
	"github.com/pkosiec/terminer/pkg/shell"
	"github.com/pkosiec/terminer/pkg/shell/automock"
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
		shImpl.On("Exec", fixCommand("echo \"C1/1\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"C2/1\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"C1/2\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"C2/2\"")).Return("", nil).Once()
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
		shImpl.On("Exec", fixCommand("echo \"C1/1\"")).Return("", testErr).Once()
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
		shImpl.On("Exec", fixCommand("echo \"R2/2\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"R1/2\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"R2/1\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"R1/1\"")).Return("", nil).Once()
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
		shImpl.On("Exec", fixCommand("echo \"R2/2\"")).Return("", testErr).Once()
		shImpl.On("Exec", fixCommand("echo \"R1/2\"")).Return("", testErr).Once()
		shImpl.On("Exec", fixCommand("echo \"R2/1\"")).Return("", nil).Once()
		shImpl.On("Exec", fixCommand("echo \"R1/1\"")).Return("", nil).Once()
		defer shImpl.AssertExpectations(t)
		i.SetShell(shImpl)

		err = i.Rollback()
		require.NoError(t, err)
	})
}

func fixCommand(run string) shell.Command {
	return shell.Command{
		Run: run,
	}
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
							Run: "echo \"C1/1\"",
						},
						Rollback: shell.Command{
							Run: "echo \"R1/1\"",
						},
					},
					{
						Metadata: recipe.UnitMetadata{
							Name: "Step 2",
							URL:  "https://step2.stage1.example.com",
						},
						Execute: shell.Command{
							Run: "echo \"C2/1\"",
						},
						Rollback: shell.Command{
							Run: "echo \"R2/1\"",
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
							Run: "echo \"C1/2\"",
						},
						Rollback: shell.Command{
							Run: "echo \"R1/2\"",
						},
					},
					{
						Metadata: recipe.UnitMetadata{
							Name: "Step 2",
							URL:  "https://step2.stage2.example.com",
						},
						Execute: shell.Command{
							Run: "echo \"C2/2\"",
						},
						Rollback: shell.Command{
							Run: "echo \"R2/2\"",
						},
					},
				},
			},
		},
	}
}
