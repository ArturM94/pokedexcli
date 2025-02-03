package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ArturM94/pokedexcli/internal/pokeapi"
	"github.com/ArturM94/pokedexcli/internal/pokecache"
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

func commandExit(cache *pokecache.Cache, config *cliConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cache *pokecache.Cache, config *cliConfig) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()

	return nil
}

func commandMap(cache *pokecache.Cache, config *cliConfig) error {
	locations, err := pokeapi.GetLocationAreas(cache, config.Next)
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cache *pokecache.Cache, config *cliConfig) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locations, err := pokeapi.GetLocationAreas(cache, config.Previous)
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*pokecache.Cache, *cliConfig) error
}

type cliConfig struct {
	Next     *string
	Previous *string
}

var commands map[string]cliCommand

func main() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Paginates over Pokemon maps",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show previous page of Pokemon maps",
			callback:    commandMapb,
		},
	}

	var config cliConfig
	cache := pokecache.NewCache(5 * time.Second)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := commands[commandName]
		if exists {
			err := command.callback(cache, &config)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
