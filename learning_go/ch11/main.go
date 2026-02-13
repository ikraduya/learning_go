package main

import (
	"embed"
	"fmt"
	"os"
	"slices"
	"strings"
)

//go:embed data
var dataFs embed.FS

var availableLanguages = []string{"english", "indonesian", "japanese"}
var availableLanguageStr = fmt.Sprintf("Available languages: %v", availableLanguages)

func printHelp() {
	fmt.Printf("%s needs a <language> argument.\n%s.\n\tFor example: %s japanese\n", os.Args[0], availableLanguageStr, os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	language := strings.ToLower(strings.TrimSpace(os.Args[1]))
	if !slices.Contains(availableLanguages, language) {
		fmt.Printf("%s language not available. %s.\n", language, availableLanguageStr)
		os.Exit(0)
	}

	textFilepath := fmt.Sprintf("data/%s_rights.txt", language)
	content, err := dataFs.ReadFile(textFilepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(content))
}
