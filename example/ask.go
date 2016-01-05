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

	query := "What is your name?"
	ans, err := ui.Ask(query, &input.Options{
		Default: "tcnksm",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Answer is %s\n", ans)
}
