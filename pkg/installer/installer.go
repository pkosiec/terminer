package installer

import (
	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/printer"
	"github.com/pkosiec/terminer/pkg/recipe"
	"github.com/pkosiec/terminer/pkg/shell"
)

// Installer provides an ability to install recipes
type Installer struct {
	r       *recipe.Recipe
	sh      shell.Shell
	printer printer.Printer
}

// New creates a new instance of Installer.
func New(r *recipe.Recipe, p printer.Printer) (*Installer, error) {
	if r == nil {
		return nil, errors.New("Recipe is empty")
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return &Installer{
		r:       r,
		sh:      shell.New(p.Command, p.ExecOutput, p.ExecError),
		printer: p,
	}, nil
}

// Install installs a recipe by executing all steps in all stages
func (installer *Installer) Install() error {
	stagesCount := len(installer.r.Stages)
	installer.printer.SetContext(printer.OperationInstall, stagesCount)

	stages := installer.r.Stages

	installer.printer.Recipe(installer.r.Metadata)

	for stageIndex, stage := range stages {
		installer.printer.Stage(stageIndex, stage)

		stepsLen := len(stage.Steps)
		for stepIndex, step := range stage.Steps {
			installer.printer.Step(stepIndex, stepsLen, step.Metadata)

			err := installer.sh.Exec(step.Execute, true)
			if err != nil {
				return errors.Wrapf(err, "while executing command from Stage '%s', Step '%s'", stage.Metadata.Name, step.Metadata.Name)
			}
		}
	}

	return nil
}

// Rollback reverts a recipe by executing all steps in all stages in reverse order
func (installer *Installer) Rollback() error {
	stagesCount := len(installer.r.Stages)
	installer.printer.SetContext(printer.OperationRollback, stagesCount)

	stages := installer.r.Stages
	stagesLen := len(stages)

	hasErrorOccurred := false

	installer.printer.Recipe(installer.r.Metadata)

	for i := stagesLen; i > 0; i-- {
		stage := stages[i-1]
		stageIndex := stagesLen - i

		installer.printer.Stage(stageIndex, stage)

		stepsLen := len(stage.Steps)
		for j := stepsLen; j > 0; j-- {
			step := stage.Steps[j-1]
			stepIndex := stepsLen - j

			installer.printer.Step(stepIndex, stepsLen, step.Metadata)

			err := installer.sh.Exec(step.Rollback, false)
			if err != nil {
				hasErrorOccurred = true
			}
		}
	}

	if hasErrorOccurred {
		return errors.New("Error(s) received during steps execution. See the logs for details.")
	}

	return nil
}
