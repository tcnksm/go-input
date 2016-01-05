package main

import (
	"fmt"
	"log"
	"os"

	"encoding/json"

	"github.com/tcnksm/go-input"
)

type User struct {
	Name, Language, Password string
}

func main() {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query1 := "What is your name?"
	name, err := ui.Ask(query1, &input.Options{
		Default: "tcnksm",
	})
	if err != nil {
		log.Fatal(err)
	}

	query2 := "Which language do you prefer to use?"
	lang, err := ui.Select(query2, []string{"go", "Go", "golang"}, &input.Options{
		//Default: "Go",
		Loop: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(&User{
		Name:     name,
		Language: lang,
	}, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

}
