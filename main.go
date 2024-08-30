package main

import (
	"bufio"
	"fmt"
	"github.com/Aleksandr-Rozhok/Pokedex/internal/PokeAPI"
	"os"
	"strings"
)

type config struct {
	Next     string
	Previous string
}

func main() {
	cfg := config{
		Next:     "https://pokeapi.co/api/v2/location-area/?limit=20&offset=20",
		Previous: "",
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

		if value, exists := commands[input]; exists {
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
%v`, "\n", "\n")

}

func commandExit(cfg *config) {
	os.Exit(0)

}

func pokeMap(cfg *config) {
	if cfg.Next == "" {
		fmt.Println("Error: You are on the last page")
	} else {
		requestResult := PokeAPI.PokeAPI(cfg.Next)
		cfg.Next = requestResult.Next
		cfg.Previous = requestResult.Previous

		locations := requestResult.Results

		for _, location := range locations {
			fmt.Println(location.Name)
		}
	}
}

func pokeMapB(cfg *config) {
	if cfg.Previous == "" {
		fmt.Println("Error: You are on the first page")
	} else {
		requestResult := PokeAPI.PokeAPI(cfg.Previous)
		cfg.Next = requestResult.Next
		cfg.Previous = requestResult.Previous

		locations := requestResult.Results

		for _, location := range locations {
			fmt.Println(location.Name)
		}
	}
}
