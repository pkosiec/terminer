package recipe

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"runtime"
)

type Recipe struct {
	Name   string
	Description string
	OS     string `yaml:"os"`
	Stages []Stage
}

type Stage struct {
	Name    string
	//Website string
	Steps   []Step
}

type Step struct {
	Name     string
	//Website  string
	Command  string `yaml:"cmd"`
	Rollback *string
}

func Load(path string) (*Recipe, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading file %s", path)
	}

	var recipe *Recipe
	err = yaml.Unmarshal(yamlFile, &recipe)
	if err != nil {
		return nil, errors.Wrapf(err, "while unmarshalling recipe from the file %s", path)
	}

	return recipe, nil
}

func ValidateCompatibility(recipe Recipe) error {
	os := runtime.GOOS
	if recipe.OS != os {
		return fmt.Errorf("the recipe %s requires %s OS. It is incompatible with your operating system %s", recipe.Name, recipe.OS, os)
	}

	return nil
}
