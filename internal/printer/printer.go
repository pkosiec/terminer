package printer

import (
	"fmt"
	"github.com/pkosiec/terminer/internal/metadata"
	"github.com/pkosiec/terminer/pkg/recipe"
	"log"
	"github.com/fatih/color"
)

func AppVersion() {
	appName := color.New(color.Bold).Sprint(metadata.AppName)
	fmt.Printf("%s %s\n", appName, metadata.Version)

	url := color.New(color.Underline).Sprint(metadata.URL)
	fmt.Printf("URL: %s\n", url)
}

func RecipeInfo(r *recipe.Recipe, operationType string) {
	color.New()
	log.Printf("%s recipe %s", operationType, r.Metadata.Name)
	log.Printf("Description: %s", r.Metadata.Description)
}

func Stage(s recipe.Stage, stageIndex, stagesLen int) {
	log.Printf("[STAGE %d/%d] %s\n", stageIndex+1, stagesLen, s.Metadata.Name)
	//TODO: Show description and URL
}

func Step(stepMetadata recipe.UnitMetadata, stepCommand string, stepIndex, stepsLen int) {
	log.Printf(">> [STEP %d/%d] %s\n", stepIndex+1, stepsLen, stepMetadata.Name)
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
