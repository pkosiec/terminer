package metadata

// AppName stores application name
const AppName string = "Terminer"

// Version stores application name
var Version = "unreleased"

// URL is an address of the application's repository
const URL string = "https://github.com/pkosiec/terminer"

// RepositoryDetails contains all details about recipes repository
type RepositoryDetails struct {
	Owner           string
	Name            string
	BranchName      string
	RecipeDirectory string
}

// Repository holds details of the official recipes repository
var Repository = RepositoryDetails{
	Owner:           "pkosiec",
	Name:            "terminer",
	BranchName:      "master",
	RecipeDirectory: "recipes",
}
