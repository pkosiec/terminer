package printer

import (
	"fmt"
	"github.com/pkosiec/terminer/pkg/shared"

	"github.com/fatih/color"
	"github.com/pkosiec/terminer/pkg/recipe"
)

// Printer is an interface of a module, which outputs text to the standard output
//go:generate mockery -name=Printer -output=automock -outpkg=automock -case=underscore
type Printer interface {
	SetContext(operation shared.Operation, stagesCount int)
	Recipe(r recipe.UnitMetadata)
	Stage(stageIndex int, s recipe.Stage)
	Step(stepIndex, steps int, s recipe.UnitMetadata)
	Command(cmd string)
	ExecOutput(output string)
	ExecError(output string)
}

type printer struct {
	operation   shared.Operation
	stages      int
	indentation string
}

// New creates a new Printer
func New() *printer {
	return &printer{}
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

func (p *printer) SetContext(operation shared.Operation, stagesCount int) {
	p.operation = operation
	p.stages = stagesCount
	p.indentation = stagesIndentation(stagesCount)
}

func (p *printer) Recipe(r recipe.UnitMetadata) {
	c := color.New(color.Bold, color.FgBlue)

	c.DisableColor()
	_, _ = c.Printf("Starting %s...\n\n", p.operation)

	c.EnableColor()
	_, _ = c.Printf("%s\n", r.Name)
	c.DisableColor()

	p.descriptionAndURL(r, "")
}

func (p *printer) Stage(stageIndex int, s recipe.Stage) {
	c := color.New(color.Bold, color.FgBlue)

	var name string
	if p.operation == shared.OperationRollback {
		name = fmt.Sprintf("Reverting '%s'", s.Metadata.Name)
	} else {
		name = s.Metadata.Name
	}

	stageCounter := fmt.Sprintf("[%d/%d] ", stageIndex+1, p.stages)
	_, _ = c.Printf("\n%s%s\n", stageCounter, name)

	p.descriptionAndURL(s.Metadata, p.indentation)
}

func (p *printer) Step(stepIndex, steps int, s recipe.UnitMetadata) {
	c := color.New(color.Bold, color.FgCyan)

	fmt.Printf("\n")

	var stepCounter string
	if steps > 1 {
		stepCounter = fmt.Sprintf("[%d/%d] ", stepIndex+1, steps)
	}

	if s.Name != "" {
		_, _ = c.Printf("%s%s%s\n", p.indentation, stepCounter, s.Name)
	}

	p.descriptionAndURL(s, p.indentation)
}

func (p *printer) Command(cmd string) {
	header := color.New(color.Faint, color.Bold)
	_, _ = header.Printf("%sCommand: ", p.indentation)
	_, _ = color.New(color.Faint).Printf("%s\n", cmd)
}

func (p *printer) ExecOutput(output string) {
	p.stepOutput(output, color.New(color.Faint))
}

func (p *printer) ExecError(output string) {
	p.stepOutput(output, color.New(color.Faint, color.FgRed))
}

func (p *printer) AppInfo(appName, version, url string) {
	appNameFmt := color.New(color.Bold).Sprint(appName)
	fmt.Printf("%s %s\n", appNameFmt, version)
	fmt.Printf("URL: %s\n", url)
}

func (p *printer) Result(err error) {
	result := color.New(color.Bold)
	_, _ = result.Printf("\n")

	if err != nil {
		_, _ = result.Add(color.FgRed).Printf("Error:\n")
		_, _ = color.New(color.FgRed).Printf(err.Error())
		return
	}

	_, _ = result.Add(color.FgGreen).Println("Success")
}

func (p *printer) descriptionAndURL(m recipe.UnitMetadata, indentation string) {
	if m.Description != "" {
		fmt.Printf("%s%s\n", indentation, m.Description)
	}

	if m.URL != "" {
		fmt.Printf("%s%s\n", indentation, m.URL)
	}
}

func (p *printer) stepOutput(output string, outputFormatter *color.Color) {
	if output == "" {
		return
	}

	_, _ = outputFormatter.Printf("%s%s\n", p.indentation, output)
}
