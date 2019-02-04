package recipe

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Recipe struct {
	Name        string
	Description string
	OS          string `yaml:"os"`
	Stages      []Stage
}

type Stage struct {
	Name        string
	Description string
	ReadMoreURL string `yaml:"url"`
	Steps       []Step
}

type Step struct {
	Name        string
	ReadMoreURL string `yaml:"url"`
	Command     string `yaml:"cmd"`
	Rollback    string
}

func FromPath(path string) (*Recipe, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading file %s", path)
	}

	recipe, err := unmarshallRecipe(yamlFile)
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

	recipe, err := unmarshallRecipe(bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading recipe from URL %s", url)
	}

	return recipe, nil
}

func unmarshallRecipe(bytes []byte) (*Recipe, error) {
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
			return errors.Wrapf(err, "while validating stage %d (%s)", stageNo+1, stage.Name)
		}
	}

	return nil
}

func (r *Recipe) validateSteps(stage Stage) error {
	if len(stage.Steps) == 0 {
		return errors.New("No steps defined")
	}

	for stepNo, step := range stage.Steps {
		if step.Command == "" {
			return fmt.Errorf("No command defined in step %d (%s)", stepNo+1, step.Name)
		}
	}

	return nil
}
