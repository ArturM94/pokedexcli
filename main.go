package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	splitted := strings.Split(trimmed, " ")
	var filtered []string

	for _, s := range splitted {
		if s == "" {
			continue
		}

		filtered = append(filtered, s)
	}

	return filtered
}

func main() {
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			break
		}

		text := scanner.Text()
		formatted := strings.ToLower(text)
		cleanedText := cleanInput(formatted)

		fmt.Println("Your command was:", cleanedText[0])
	}
}
