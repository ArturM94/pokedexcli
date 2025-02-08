package main

import (
	"bufio"
	"fmt"
	"math/rand"
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

func commandCatch(ctx *commandContext) error {
	fmt.Println("Throwing a Pokeball at " + ctx.PokemonName + "...")

	pokemon, err := pokeapi.GetPokemon(ctx.Cache, ctx.PokemonName)
	if err != nil {
		return fmt.Errorf("error catching pokemon: %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	chance := rand.Intn(100)
	catchRate := 100 - (pokemon.BaseExperience / 3)

	if chance < catchRate {
		fmt.Println(pokemon.Name + " was caught!")
		ctx.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Println(pokemon.Name + " escaped!")
	}

	return nil
}

type commandContext struct {
	Cache        *pokecache.Cache
	Config       *cliConfig
	Pokedex      map[string]*pokeapi.GetPokemonResponse
	LocationName string
	PokemonName  string
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
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
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
		var pokemonName string
		pokedex := make(map[string]*pokeapi.GetPokemonResponse)

		command, exists := commands[commandName]
		if exists {
			if command.name == "explore" && len(words) == 2 {
				locationName = words[1]
			}
			if command.name == "catch" && len(words) == 2 {
				pokemonName = words[1]
			}

			ctx := &commandContext{cache, &config, pokedex, locationName, pokemonName}
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
