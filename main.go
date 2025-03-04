package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		newScanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Pokedex > ")
		newScanner.Scan()

		userInput := newScanner.Text()
		cleanedInput := cleanInput(userInput)

		if command, ok := commands[cleanedInput[0]]; ok {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	splitText := strings.Split(lowerText, " ")

	return splitText
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
}
