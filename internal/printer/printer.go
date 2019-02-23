package printer

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pkosiec/terminer/internal/metadata"
	"github.com/pkosiec/terminer/pkg/recipe"
	"strings"
	"time"
)

func AppInfo() {
	appName := color.New(color.Bold).Sprint(metadata.AppName)
	fmt.Printf("%s %s\n", appName, metadata.Version)

	url := color.New(color.Underline).Sprint(metadata.URL)
	fmt.Printf("URL: %s\n", url)
}

func Error(err error) {
	color.New(color.FgRed, color.Bold).Printf("Error:\n")
	color.New(color.FgRed).Printf(err.Error())
}

// TODO: Use spinner
func Spinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[40], 100*time.Millisecond)
}

type Printer struct {
	stages      int
	indentation string
}

func NewPrinter(stages int) *Printer {
	return &Printer{
		stages:      stages,
		indentation: stagesIndentation(stages),
	}
}

func stagesIndentation(stagesCount int) string {
	var digitsCount int

	for stagesCount > 0 {
		digitsCount++
		stagesCount = stagesCount / 10
	}

	count := digitsCount*2 + 4 // "[" + digits + "/" + digits + "]" + " "

	var indentation string
	for i := 0; i < count; i++ {
		indentation = fmt.Sprintf("%s%s", indentation, " ")
	}

	return indentation
}

func (p *Printer) Recipe(r recipe.UnitMetadata, operationType string) {
	c := color.New(color.Bold, color.FgBlue)

	c.DisableColor()
	c.Printf("Starting %s...\n\n", operationType)

	c.EnableColor()
	c.Printf("%s\n", r.Name)
	c.DisableColor()

	p.descriptionAndURL(r, "")
}

func (p *Printer) Stage(stageIndex int, s recipe.Stage) {
	c := color.New(color.Bold, color.FgBlue)

	stageCounter := fmt.Sprintf("[%d/%d] ", stageIndex+1, p.stages)
	c.Printf("\n%s%s\n", stageCounter, s.Metadata.Name)

	p.descriptionAndURL(s.Metadata, p.indentation)
}

func (p *Printer) Step(s recipe.UnitMetadata, stepCommand string, stepIndex, steps int) {
	c := color.New(color.Bold, color.FgCyan)

	fmt.Printf("\n")

	var stepCounter string
	if steps > 1 {
		stepCounter = fmt.Sprintf("[%d/%d] ", stepIndex+1, steps)
	}

	if s.Name != "" {
		c.Printf("%s%s%s\n", p.indentation, stepCounter, s.Name)
	}

	p.descriptionAndURL(s, p.indentation)

	color.New(color.Faint, color.Bold).Printf("%sCommand: ", p.indentation)
	color.New(color.Faint).Printf("%s\n", stepCommand)
}

func (p *Printer) StepOutput(output string) {
	p.stepOutput(output, color.New(color.Faint))
}

func (p *Printer) StepError(err error) {
	p.stepOutput(err.Error(), color.New(color.FgRed))
}

func (p *Printer) descriptionAndURL(m recipe.UnitMetadata, indentation string) {
	if m.Description != "" {
		fmt.Printf("%s%s\n", indentation, m.Description)
	}

	if m.URL != "" {
		fmt.Printf("%s%s\n", indentation, m.URL)
	}
}

func (p *Printer) stepOutput(output string, outputFormatter *color.Color) {
	if output == "" {
		return
	}

	color.New(color.Faint, color.Bold).Printf("%sResult:\n", p.indentation)

	formattedOutput := strings.Replace(output, "\n", fmt.Sprintf("\n%s", p.indentation), -1)
	outputFormatter.Printf("%s%s", p.indentation, formattedOutput)
}
