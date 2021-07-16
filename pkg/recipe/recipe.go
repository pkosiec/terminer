package recipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkosiec/terminer/internal/metadata"
	"github.com/pkosiec/terminer/pkg/path"
	"github.com/pkosiec/terminer/pkg/shell"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// AnyOS is OS string that matches any operating system
const AnyOS = "any"

// UnitMetadata stores metadata for a generic Recipe unit, such as Recipe, Stage or Step
type UnitMetadata struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	URL         string `yaml:"url" json:"url"`
}

// Recipe stores needed steps to install a gjven piece of functionality
type Recipe struct {
	OS       string       `yaml:"os" json:"os"`
	Metadata UnitMetadata `yaml:"metadata" json:"metadata"`
	Stages   []Stage      `yaml:"stages" json:"stages"`
}

// Stage represents a logical part of recipe that consists of steps
type Stage struct {
	Metadata UnitMetadata `yaml:"metadata" json:"metadata"`
	Steps    []Step       `yaml:"steps" json:"steps"`
}

// Step contains data about a single shell command, which can be installed or reverted
type Step struct {
	Metadata UnitMetadata  `yaml:"metadata" json:"metadata"`
	Execute  shell.Command `yaml:"execute" json:"execute"`
	Rollback shell.Command `yaml:"rollback" json:"rollback"`
}

// FromPath creates a Recipe from given file
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

// HTTPClient is an interface that is used for HTTP requests
//go:generate mockery -name=HTTPClient -output=automock -outpkg=automock -case=underscore
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// FromURL downloads a file from given URL and uses it to create a Recipe
func FromURL(url string, httpClient HTTPClient) (*Recipe, int, error) {
	if !path.IsURL(url) {
		return nil, 0, fmt.Errorf("Incorrect recipe URL")
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "while requesting recipe from URL %s", url)
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Println(errors.Wrapf(err, "while closing response body").Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, fmt.Errorf("Invalid status code while downloading file from URL %s: %d. Expected: %d", url, res.StatusCode, http.StatusOK)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, errors.Wrapf(err, "while reading response body while downloading file from URL %s", url)
	}

	if len(bytes) == 0 {
		return nil, res.StatusCode, fmt.Errorf("Empty body while downloading file from URL %s", url)
	}

	recipe, err := unmarshalRecipe(bytes)
	if err != nil {
		return nil, res.StatusCode, errors.Wrapf(err, "while loading recipe from URL %s", url)
	}

	return recipe, res.StatusCode, nil
}

// FromRepository downloads a recipe from official recipes repository
func FromRepository(recipeName string, httpClient HTTPClient) (*Recipe, error) {
	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/%s/%s/%s/%s/%s.yaml",
		metadata.Repository.Owner,
		metadata.Repository.Name,
		metadata.Repository.BranchName,
		metadata.Repository.RecipeDirectory,
		recipeName,
		runtime.GOOS,
	)

	recipeListURL := fmt.Sprintf("https://github.com/%s/%s/tree/%s/%s",
		metadata.Repository.Owner,
		metadata.Repository.Name,
		metadata.Repository.BranchName,
		metadata.Repository.RecipeDirectory,
	)

	var statusCode int
	r, statusCode, err := FromURL(url, httpClient)
	if err != nil {
		if statusCode == http.StatusNotFound {
			return nil, fmt.Errorf("Cannot find recipe `%s` on official repository.\nSee the official list of the recipes on %s\n", recipeName, recipeListURL)
		}

		return nil, errors.Wrapf(err, "Error while finding recipe `%s` on official repository", recipeName)
	}

	return r, nil
}

// Validate checks if the recipe is valid to run on current OS and whether all stages and steps are not empty
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

func validateExtension(path string) error {
	ext := filepath.Ext(path)
	lowercaseExt := strings.ToLower(ext)

	if lowercaseExt != ".yaml" && lowercaseExt != ".yml" && lowercaseExt != ".json" {
		return fmt.Errorf("Invalid file extension `%s`. Expected: yaml, yml or json", ext)
	}

	return nil
}

func unmarshalRecipe(bytes []byte) (*Recipe, error) {
	var recipe *Recipe

	if json.Valid(bytes) {
		err := json.Unmarshal(bytes, &recipe)
		return recipe, err
	}

	err := yaml.Unmarshal(bytes, &recipe)
	return recipe, err
}

func (r *Recipe) validateOS() error {
	os := runtime.GOOS
	if r.OS != os && r.OS != AnyOS {
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
		if len(step.Execute.Run) == 0 {
			return fmt.Errorf("No commands defined in step %d (%s)", stepNo+1, step.Metadata.Name)
		}
	}

	return nil
}
