package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Aleksandr-Rozhok/internal/PokeAPI"
	"github.com/Aleksandr-Rozhok/internal/Pokecache"
	"os"
	"strings"
	"time"
)

type config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationData struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []LocationResult `json:"results"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type PokemonData struct {
	ID                int                 `json:"id"`
	Name              string              `json:"name"`
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

var arg string

func main() {
	cfg := config{
		Next:     "https://pokeapi.co/api/v2/location-area/?limit=20&offset=20",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Millisecond),
	}

	type cliCommand struct {
		name        string
		description string
		callback    func(*config)
	}

	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations in the Pokemon world",
			callback:    pokeMap,
		},
		"mapb": {
			name:        "map",
			description: "Displays 20 previous locations in the Pokemon world",
			callback:    pokeMapB,
		},
		"explore": {
			name:        "explore",
			description: "Displays 10 pokemons in a given area",
			callback:    explore,
		},
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("pokedex > ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		command := strings.Split(input, " ")

		if value, exists := commands[command[0]]; exists {
			if len(command) > 1 {
				arg = command[1]
			}

			value.callback(&cfg)
		} else {
			fmt.Println("Unknown command:", input)
		}
	}
}

func commandHelp(cfg *config) {
	fmt.Printf(`%vWelcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: Displays 20 locations in the Pokemon world
mapb: Displays 20 previous locations in the Pokemon world
explore: Displays 10 pokemons in a given area
%v`, "\n", "\n")

}

func commandExit(cfg *config) {
	os.Exit(0)

}

func pokeMap(cfg *config) {
	cache, exists := cfg.Cache.Get(cfg.Next)
	fmt.Printf("Cache exists: %v\n", exists)
	if exists {
		var result LocationData
		err := json.Unmarshal(cache, &result)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		cfg.Next = result.Next
		cfg.Previous = result.Previous

		locations := result.Results

		for _, location := range locations {
			fmt.Println(location.Name)
		}
	} else if cfg.Next == "" {
		fmt.Println("Error: You are on the last page")
	} else {
		body := pokeAPI.PokeAPI(cfg.Next)
		cfg.Cache.Add(cfg.Next, body)

		var result LocationData

		err := json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error:", err)
		}

		cfg.Next = result.Next
		cfg.Previous = result.Previous

		locations := result.Results

		for _, location := range locations {
			fmt.Println(location.Name)
		}
	}
}

func pokeMapB(cfg *config) {
	cache, exists := cfg.Cache.Get(cfg.Previous)
	fmt.Printf("Cache exists: %v\n", exists)
	if exists {
		var result LocationData
		err := json.Unmarshal(cache, &result)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		cfg.Next = result.Next
		cfg.Previous = result.Previous

		locations := result.Results

		for _, location := range locations {
			fmt.Println(location.Name)
		}
	} else if cfg.Previous == "" {
		fmt.Println("Error: You are on the first page")
	} else {
		body := pokeAPI.PokeAPI(cfg.Previous)

		cfg.Cache.Add(cfg.Previous, body)

		var result LocationData

		err := json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error:", err)
		}

		cfg.Next = result.Next
		cfg.Previous = result.Previous

		locations := result.Results

		for _, location := range locations {
			fmt.Println(location.Name)
		}
	}
}

func explore(cfg *config) {
	location := arg
	link := "https://pokeapi.co/api/v2/location-area/" + location
	fmt.Printf("Exploring %s...\n", location)

	cache, exists := cfg.Cache.Get(link)
	if exists {
		var result PokemonData
		err := json.Unmarshal(cache, &result)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		pokemons := result.PokemonEncounters
		for _, pokemon := range pokemons {
			fmt.Println(pokemon.Pokemon.Name)
		}
	} else {
		body := pokeAPI.PokeAPI(link)
		cfg.Cache.Add(link, body)

		var result PokemonData

		err := json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error:", err)
		}

		pokemons := result.PokemonEncounters
		for _, pokemon := range pokemons {
			fmt.Println(pokemon.Pokemon.Name)
		}
	}
}
