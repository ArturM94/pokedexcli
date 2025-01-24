package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	splitted := strings.Split(trimmed, " ")
	filtered := []string{}

	for _, s := range splitted {
		if s == "" {
			continue
		}

		filtered = append(filtered, s)
	}

	return filtered
}

func main() {
	fmt.Println("Hello, World!")
}
