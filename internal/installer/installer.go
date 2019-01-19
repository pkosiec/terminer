package installer

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/pkosiec/terminer/internal/recipe"
	"github.com/pkosiec/terminer/internal/sh"
)

type Installer struct {
	r *recipe.Recipe
}

func New(r *recipe.Recipe) (*Installer, error) {
	if r == nil {
		return nil, errors.New("")
	}

	if err := recipe.ValidateCompatibility(*r); err != nil {
		return nil, fmt.Errorf("Recipe is incompatible with your OS %s. It requires %s operating system", r.OS)
	}

	return &Installer{
		r: r,
	}, nil
}

func (i *Installer) Install() error {
	log.Printf("Installing recipe %s", i.r.Name)
	log.Printf("Description: %s", i.r.Description)
	stages := i.r.Stages
	stagesCount := len(stages)

	for stageNo, stage := range stages {
		log.Println("[%s/%s] STAGE %s", stageNo+1, stagesCount, stage.Name)

		stepsCount := len(stage.Steps)
		for stepNo, step := range stage.Steps {
			log.Println(">> [%s/%s] STEP %s", stepNo+1, stepsCount, step.Name)
			res, err := sh.Exec(step.Command)
			if err != nil {
				return errors.Wrapf(err, "while executing command from Stage %s, Step %s", stage.Name, step.Name)
			}

			log.Printf("%s\n", res)
		}
	}

	return nil
}

func (i *Installer) Rollback() {

}
