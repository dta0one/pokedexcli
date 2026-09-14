package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dta0one/pokedexcli/internal/pokecache"
)

type config struct {
	cache               pokecache.Cache
	nextLocationAreaURL *string
	prevLocationAreaURL *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error // 👈 Add []string for args
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:] // 👈 Extract any arguments passed after the command
		}

		cmd, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(cfg, args)
		if err != nil {
			fmt.Println(err)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config, args []string) error {
	resp, err := fetchLocationAreas(cfg, cfg.nextLocationAreaURL)
	if err != nil {
		return err
	}

	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.prevLocationAreaURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	resp, err := fetchLocationAreas(cfg, cfg.prevLocationAreaURL)
	if err != nil {
		return err
	}

	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location area for Pokemon",
			callback:    commandExplore,
		},
	}
}

func commandExplore(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("explore requires exactly one location area name")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	resp, err := fetchLocationArea(cfg, areaName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
