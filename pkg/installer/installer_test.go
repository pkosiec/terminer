package installer_test

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/printer"
	printerAutomock "github.com/pkosiec/terminer/internal/printer/automock"
	"github.com/pkosiec/terminer/pkg/installer"
	"github.com/pkosiec/terminer/pkg/recipe"
	"github.com/pkosiec/terminer/pkg/shell"
	"github.com/pkosiec/terminer/pkg/shell/automock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Empty recipe", func(t *testing.T) {
		var r *recipe.Recipe

		p := &printerAutomock.Printer{}
		_, err := installer.New(r, p)

		require.Error(t, err)
	})

	t.Run("Invalid recipe", func(t *testing.T) {
		r := fixRecipe("testos")

		p := &printerAutomock.Printer{}
		_, err := installer.New(r, p)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid operating system")
	})

	t.Run("Valid recipe", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		p := &printerAutomock.Printer{}
		_, err := installer.New(r, p)

		require.NoError(t, err)
	})
}

func TestInstaller_Install(t *testing.T) {
	printerFn := mock.AnythingOfType("shell.PrintFn")

	t.Run("Success", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		p := &printerAutomock.Printer{}
		p.On("SetContext", printer.OperationInstall, 2).Return().Once()
		p.On("Recipe", r.Metadata).Return().Once()
		defer p.AssertExpectations(t)

		shImpl := &automock.Shell{}
		defer shImpl.AssertExpectations(t)

		for stageIdx, stage := range r.Stages {
			p.On("Stage", stageIdx, stage).Return().Once()

			for stepIdx, step := range stage.Steps {
				p.On("Step", stepIdx, len(stage.Steps), step.Execute.Run, step.Metadata).Return().Once()
				shImpl.On("Exec", fixCommand(step.Execute.Run), printerFn, printerFn).Return( nil).Once()
			}
		}

		i, err := installer.New(r, p)
		require.NoError(t, err)

		i.SetShell(shImpl)

		err = i.Install()
		require.NoError(t, err)
	})

	// Should exit after failed step
	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("Test Err")
		r := fixRecipe(runtime.GOOS)

		p := &printerAutomock.Printer{}
		p.On("SetContext", printer.OperationInstall, 2).Return().Once()
		p.On("Recipe", r.Metadata).Return().Once()
		defer p.AssertExpectations(t)

		stage := r.Stages[0]
		step := stage.Steps[0]
		p.On("Stage", 0, stage).Return().Once()
		p.On("Step", 0, len(stage.Steps), step.Execute.Run, step.Metadata).Return().Once()


		shImpl := &automock.Shell{}
		shImpl.On("Exec", fixCommand("echo \"C1/1\""), printerFn, printerFn).Return(testErr).Once()
		defer shImpl.AssertExpectations(t)

		i, err := installer.New(r, p)
		require.NoError(t, err)

		i.SetShell(shImpl)

		err = i.Install()
		require.Error(t, err)
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestInstaller_Rollback(t *testing.T) {
	printerFn := mock.AnythingOfType("shell.PrintFn")

	t.Run("Success", func(t *testing.T) {
		r := fixRecipe(runtime.GOOS)

		p := &printerAutomock.Printer{}
		p.On("SetContext", printer.OperationRollback, 2).Return().Once()
		p.On("Recipe", r.Metadata).Return().Once()
		defer p.AssertExpectations(t)

		stage := r.Stages[1]
		p.On("Stage", 0, stage).Return().Once()
		p.On("Step", 0, len(stage.Steps), stage.Steps[1].Rollback.Run, stage.Steps[1].Metadata).Return().Once()
		p.On("Step", 1, len(stage.Steps), stage.Steps[0].Rollback.Run, stage.Steps[0].Metadata).Return().Once()

		stage = r.Stages[0]
		p.On("Stage", 1, stage).Return().Once()
		p.On("Step", 0, len(r.Stages[0].Steps), stage.Steps[1].Rollback.Run, stage.Steps[1].Metadata).Return().Once()
		p.On("Step", 1, len(r.Stages[0].Steps), stage.Steps[0].Rollback.Run, stage.Steps[0].Metadata).Return().Once()

		shImpl := &automock.Shell{}
		defer shImpl.AssertExpectations(t)

		for _, stage := range r.Stages {
			for _, step := range stage.Steps {
				shImpl.On("Exec", fixCommand(step.Rollback.Run), printerFn, printerFn).Return( nil).Once()
			}
		}

		i, err := installer.New(r, p)
		require.NoError(t, err)

		i.SetShell(shImpl)

		err = i.Rollback()
		require.NoError(t, err)
	})

	// Should not exit after failed step
	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("Test Err")
		r := fixRecipe(runtime.GOOS)

		p := &printerAutomock.Printer{}
		p.On("SetContext", printer.OperationRollback, 2).Return().Once()
		p.On("Recipe", r.Metadata).Return().Once()

		stage := r.Stages[1]
		p.On("Stage", 0, stage).Return().Once()
		p.On("Step", 0, len(stage.Steps), stage.Steps[1].Rollback.Run, stage.Steps[1].Metadata).Return().Once()
		p.On("Step", 1, len(stage.Steps), stage.Steps[0].Rollback.Run, stage.Steps[0].Metadata).Return().Once()

		stage = r.Stages[0]
		p.On("Stage", 1, stage).Return().Once()
		p.On("Step", 0, len(r.Stages[0].Steps), stage.Steps[1].Rollback.Run, stage.Steps[1].Metadata).Return().Once()
		p.On("Step", 1, len(r.Stages[0].Steps), stage.Steps[0].Rollback.Run, stage.Steps[0].Metadata).Return().Once()
		defer p.AssertExpectations(t)

		shImpl := &automock.Shell{}
		defer shImpl.AssertExpectations(t)

		i, err := installer.New(r, p)
		require.NoError(t, err)

		for _, stage := range r.Stages {
			for _, step := range stage.Steps {
				shImpl.On("Exec", fixCommand(step.Rollback.Run), printerFn, printerFn).Return( testErr).Once()
			}
		}

		i.SetShell(shImpl)

		err = i.Rollback()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error(s) received during steps execution")
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
