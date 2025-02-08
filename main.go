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

func commandExit(ctx *commandContext) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(ctx *commandContext) error {
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

func commandMap(ctx *commandContext) error {
	locations, err := pokeapi.GetLocationAreas(ctx.Cache, ctx.Config.Next)
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	ctx.Config.Next = locations.Next
	ctx.Config.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(ctx *commandContext) error {
	if ctx.Config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locations, err := pokeapi.GetLocationAreas(ctx.Cache, ctx.Config.Previous)
	if err != nil {
		return fmt.Errorf("error getting location areas: %w", err)
	}

	ctx.Config.Next = locations.Next
	ctx.Config.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(ctx *commandContext) error {
	fmt.Println("Exploring " + ctx.LocationName + "...")

	locationDetails, err := pokeapi.GetLocationAreaDetails(ctx.Cache, ctx.LocationName)
	if err != nil {
		return fmt.Errorf("error getting location detals: %w", err)
	}

	if len(locationDetails.PokemonEncounters) == 0 {
		fmt.Println("Pokemon not found")
		
		return nil
	}

	fmt.Println("Found Pokemon:")

	for _, pokemonEncounter := range locationDetails.PokemonEncounters {
		fmt.Println(" - " + pokemonEncounter.Pokemon.Name)
	}

	return nil
}

type commandContext struct {
	Cache            *pokecache.Cache
	Config           *cliConfig
	LocationName string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandContext) error
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
		"explore": {
			name:        "explore",
			description: "Shows all Pokemons in the area",
			callback:    commandExplore,
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

		var locationName string
		if len(words) == 2 {
			locationName = words[1]
		}

		command, exists := commands[commandName]
		if exists {
			ctx := &commandContext{cache, &config, locationName}
			err := command.callback(ctx)
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
