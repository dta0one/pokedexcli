package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// Wait for user input
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		// Clean input: lowercase and split into words
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		firstWord := words[0]
		fmt.Printf("Your command was: %s\n", firstWord)
	}
}
