package installer

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/printer"
	"github.com/pkosiec/terminer/internal/recipe"
	"github.com/pkosiec/terminer/internal/shell"
)

type Installer struct {
	r  *recipe.Recipe
	sh shell.Shell
}

func New(r *recipe.Recipe) (*Installer, error) {
	if r == nil {
		return nil, errors.New("Recipe is empty")
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return &Installer{
		r:  r,
		sh: shell.New(),
	}, nil
}

func (installer *Installer) Install() error {
	printer.RecipeInfo(installer.r, "Installing")

	stages := installer.r.Stages
	stagesLen := len(stages)

	for stageIndex, stage := range stages {
		printer.Stage(stage, stageIndex, stagesLen)

		stepsLen := len(stage.Steps)
		for stepIndex, step := range stage.Steps {
			printer.Step(step.Name, step.Command, stepIndex, stepsLen)

			stepShell := installer.shellForStep(step)

			res, err := installer.sh.Exec(stepShell, step.Command)
			printer.StepOutput(res)
			if err != nil {
				return errors.Wrapf(err, "while executing command from Stage '%s', Step '%s'", stage.Name, step.Name)
			}
		}
	}

	return nil
}

func (installer *Installer) Rollback() error {
	printer.RecipeInfo(installer.r, "Uninstalling")

	stages := installer.r.Stages
	stagesLen := len(stages)

	for i := stagesLen; i > 0; i-- {
		stage := stages[i-1]
		stageIndex := stagesLen - i
		printer.Stage(stage, stageIndex, stagesLen)

		stepsLen := len(stage.Steps)
		for j := stepsLen; j > 0; j-- {
			step := stage.Steps[j-1]
			stepIndex := stepsLen - j

			printer.Step(step.Name, step.Rollback, stepIndex, stepsLen)

			stepShell := installer.shellForStep(step)
			res, err := installer.sh.Exec(stepShell, step.Rollback)
			printer.StepOutput(res)
			if err != nil {
				// Print error and continue
				wrappedErr := errors.Wrapf(err, "while executing command from Stage %s, Step %s", stage.Name, step.Name)
				printer.Error(wrappedErr)
			}
		}
	}

	return nil
}

func (installer *Installer) shellForStep(step recipe.Step) string {
	if step.Shell == "" {
		return shell.DefaultShell
	}

	return step.Shell
}
