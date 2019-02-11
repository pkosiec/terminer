package printer

import (
	"github.com/pkosiec/terminer/internal/metadata"
	"github.com/pkosiec/terminer/internal/recipe"
	"log"
)

func AppVersion() {
	log.Printf("%s v. %s\n%s", metadata.AppName, metadata.Version, metadata.URL)
}

func RecipeInfo(r *recipe.Recipe, operationType string) {
	log.Printf("%s recipe %s", operationType, r.Name)
	log.Printf("Description: %s", r.Description)
}

func Stage(s recipe.Stage, stageIndex, stagesLen int) {
	log.Printf("[STAGE %d/%d] %s\n", stageIndex+1, stagesLen, s.Name)
	//TODO: Show description and URL
}

func Step(stepName, stepCommand string, stepIndex, stepsLen int) {
	log.Printf(">> [STEP %d/%d] %s\n", stepIndex+1, stepsLen, stepName)
	//TODO: Show description and URL
	log.Printf(">> Command: %s", stepCommand)
}

func StepOutput(output string) {
	if output == "" {
		return
	}

	log.Printf("Output: %s\n", output)
}

func Error(err error) {
	log.Println(err)
}