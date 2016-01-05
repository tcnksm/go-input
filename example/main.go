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

	query := "Which language do you prefer to use?"
	ans, err := ui.Select(query, []string{"go", "Go", "golang"}, &input.Options{
		Default: "Go",
		Loop:    true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print the answer
	log.Println("Selected answer is", ans)
}
