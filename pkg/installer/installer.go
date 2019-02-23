package installer

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/printer"
	"github.com/pkosiec/terminer/pkg/recipe"
	"github.com/pkosiec/terminer/pkg/shell"
)

// Installer provides an ability to install recipes
type Installer struct {
	r  *recipe.Recipe
	sh shell.Shell
}

// New creates a new instance of Installer.
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

// Install installs a recipe by executing all steps in all stages
func (installer *Installer) Install() error {
	stages := installer.r.Stages
	stagesLen := len(stages)

	p := printer.NewPrinter(stagesLen, printer.OperationInstall)
	p.Recipe(installer.r.Metadata)

	for stageIndex, stage := range stages {
		p.Stage(stageIndex, stage)

		stepsLen := len(stage.Steps)
		for stepIndex, step := range stage.Steps {
			p.Step(step.Metadata, step.Execute.Run, stepIndex, stepsLen)

			res, err := installer.sh.Exec(step.Execute)
			p.StepOutput(res)
			if err != nil {
				return errors.Wrapf(err, "while executing command from Stage '%s', Step '%s'", stage.Metadata.Name, step.Metadata.Name)
			}
		}
	}



	return nil
}

// Rollback reverts a recipe by executing all steps in all stages in reverse order
func (installer *Installer) Rollback() error {
	stages := installer.r.Stages
	stagesLen := len(stages)

	p := printer.NewPrinter(stagesLen, printer.OperationRollback)
	p.Recipe(installer.r.Metadata)

	for i := stagesLen; i > 0; i-- {
		stage := stages[i-1]
		stageIndex := stagesLen - i

		p.Stage(stageIndex, stage)

		stepsLen := len(stage.Steps)
		for j := stepsLen; j > 0; j-- {
			step := stage.Steps[j-1]
			stepIndex := stepsLen - j

			p.Step(step.Metadata, step.Rollback.Run, stepIndex, stepsLen)

			res, err := installer.sh.Exec(step.Rollback)
			p.StepOutput(res)
			if err != nil {
				// Print error and continue
				wrappedErr := errors.Wrapf(err, "while executing command from Stage %s, Step %s", stage.Metadata.Name, step.Metadata.Name)
				p.StepError(wrappedErr)
			}
		}
	}

	return nil
}
