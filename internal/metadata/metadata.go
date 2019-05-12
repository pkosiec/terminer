package metadata

const AppName string = "Terminer"
const Version string = "unreleased"
const URL string = "https://github.com/pkosiec/terminer"

type RepositoryDetails struct {
	Owner string
	Name string
	BranchName string
}


var Repository = RepositoryDetails{
	Owner: "pkosiec",
	Name: "terminer",
	BranchName: "master",
}
