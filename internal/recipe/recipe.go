package recipe

import (
	"fmt"
	"github.com/pkosiec/terminer/internal/shell"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type UnitMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

type Recipe struct {
	OS       string `yaml:"os"`
	Metadata UnitMetadata
	Stages   []Stage
}

type Stage struct {
	Metadata UnitMetadata
	Steps    []Step
}

type Step struct {
	Metadata UnitMetadata
	Execute  shell.Command
	Rollback shell.Command
}

func FromPath(path string) (*Recipe, error) {
	err := validateExtension(path)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading file from path `%s`", path)
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading file from path `%s`", path)
	}

	recipe, err := unmarshalRecipe(yamlFile)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading recipe from file %s", path)
	}

	return recipe, nil
}

func FromURL(url string) (*Recipe, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "while requesting recipe from URL %s", url)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code while downloading file from URL %s: %d. Expected: %d", url, res.StatusCode, http.StatusOK)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading response body while downloading file from URL %s", url)
	}

	if len(bytes) == 0 {
		return nil, fmt.Errorf("Empty body while downloading file from URL %s", url)
	}

	recipe, err := unmarshalRecipe(bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading recipe from URL %s", url)
	}

	return recipe, nil
}

func validateExtension(path string) error {
	ext := filepath.Ext(path)
	lowercaseExt := strings.ToLower(ext)

	if lowercaseExt != ".yaml" && lowercaseExt != ".yml" {
		return fmt.Errorf("Invalid file extension `%s`. Expected: yaml or yml", ext)
	}

	return nil
}

func unmarshalRecipe(bytes []byte) (*Recipe, error) {
	var recipe *Recipe
	err := yaml.Unmarshal(bytes, &recipe)
	return recipe, err
}

func (r *Recipe) Validate() error {
	err := r.validateOS()
	if err != nil {
		return err
	}

	err = r.validateStages()
	if err != nil {
		return err
	}

	return nil
}

func (r *Recipe) validateOS() error {
	os := runtime.GOOS
	if r.OS != os {
		return fmt.Errorf("Invalid operating system. Required: %s. Actual: %s", r.OS, os)
	}

	return nil
}

func (r *Recipe) validateStages() error {
	if len(r.Stages) == 0 {
		return fmt.Errorf("No stages defined in recipe")
	}

	for stageNo, stage := range r.Stages {
		err := r.validateSteps(stage)
		if err != nil {
			return errors.Wrapf(err, "while validating stage %d (%s)", stageNo+1, stage.Metadata.Name)
		}
	}

	return nil
}

func (r *Recipe) validateSteps(stage Stage) error {
	if len(stage.Steps) == 0 {
		return errors.New("No steps defined")
	}

	for stepNo, step := range stage.Steps {
		if step.Execute.Run == "" {
			return fmt.Errorf("No command defined in step %d (%s)", stepNo+1, step.Metadata.Name)
		}
	}

	return nil
}
