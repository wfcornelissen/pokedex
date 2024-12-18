package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

var commands map[string]cliCommand
var pag = 0
var cfg Config = Config{}

const locationURL = "https://pokeapi.co/api/v2/location-area/?limit=20"

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpSectionWrapper,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}

		cleanText := cleanInput(text)
		commandName := cleanText[0]
		if command, exists := commands[commandName]; exists {
			err := command.callback(&cfg)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func helpSectionWrapper(cfg *Config) error {
	return helpSection(commands)
}

func helpSection(list map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for command := range list {
		fmt.Println(command+": ", list[command].description)
	}
	return nil
}

func commandMap(cfg *Config) error {
	var linkToGet string
	if cfg.Next == "" {
		linkToGet = locationURL
	} else {
		linkToGet = cfg.Next
	}
	pag += 20
	cfg.Previous = cfg.Next
	cfg.Next = locationURL + "&offset=" + strconv.Itoa(pag)

	//HTTP GET request
	res, err := http.Get(linkToGet)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create var to capture HTTP response into
	var Response LocationResults

	// Capture HTTP response
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&Response)
	if err != nil {
		return err
	}

	// Create a map of id:name for cities
	finalList := make(map[int]string)

	for locationID, locationName := range Response.Results {
		finalList[locationID] = locationName.Name
	}

	// Debug print
	for city := range finalList {
		fmt.Println(finalList[city])
	}

	return nil
}

func commandMapb(cfg *Config) error {
	pag -= 40
	if cfg.Next == "https://pokeapi.co/api/v2/location-area/?limit=5" {
		fmt.Println("you're on the first page")
		return nil
	}
	commandMap(cfg)
	return nil
}

type LocationResults struct {
	Results []Location
}

type Location struct {
	Id   int
	Name string
}

type Config struct {
	Next     string
	Previous string
}
