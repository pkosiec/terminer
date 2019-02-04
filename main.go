package main

import (
	"github.com/pkosiec/terminer/cmd"
	"log"
)

func main() {
	cmd.Execute()
}

func exitOnError(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}
