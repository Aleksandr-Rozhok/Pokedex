package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	type cliCommand struct {
		name        string
		description string
		callback    func() error
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
			err := value.callback()
			if err != nil {
				break
			}
		} else {
			fmt.Println("Unknown command:", input)
		}
	}
}

func commandHelp() error {
	fmt.Printf(`%vWelcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
%v`, "\n", "\n")

	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}
