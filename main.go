package main

import (
	"log"

	"github.com/pkosiec/terminer/internal/installer"
	"github.com/pkosiec/terminer/internal/recipe"
)

func main() {
	path := "./recipes/fish.yaml"

	r, err := recipe.From(path)
	exitOnError(err)

	i, err := installer.New(r)
	exitOnError(err)

	err = i.Install()
	exitOnError(err)
}

func exitOnError(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}
