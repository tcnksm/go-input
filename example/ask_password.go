package main

import (
	"log"
	"os"

	"github.com/tcnksm/go-input"
)

func main() {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "What is your password?"
	name, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Default:  os.Getenv("NAME"),
		Required: true,
		Mask:     true,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Answer is %s\n", name)
}
