package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rogpeppe/go-internal/cache"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type locationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var registry map[string]cliCommand
var mapCount int

func init() {
	mapCount = 0
	registry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    CommandExit,
		},
		"help": {
			name:        "help",
			description: "Display a list of commands",
			callback:    CommandHelp,
		},
		"map": {
			name: "map",
			description: `Display a list of 20 locations. Subsequent calls
			will result in the next 20 locations being displayed.`,
			callback: CommandMap,
		},
		"mapb": {
			name: "mapb",
			description: `Display previous list of 20 locations. Subesquent calls
			will result in the previous 20 locations being displayed.`,
			callback: CommandMapB,
		},
	}
}

func CleanInput(text string) []string {
	var result []string

	for _, word := range strings.Fields(text) {
		word = strings.ToLower(word)
		result = append(result, word)
	}
	return result
}

func CheckInput(command string) {
	for _, regCmd := range registry {
		if command == regCmd.name {
			err := regCmd.callback()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}
func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp() error {
	for _, regCmd := range registry {
		fmt.Printf("Name: %v\n", regCmd.name)
		fmt.Printf("Description: %v\n", regCmd.description)
		fmt.Println("----------------")
	}
	return nil
}

func CommandMap() error {
	mapURL := "https://pokeapi.co/api/v2/location-area/"
	fullURL := ""
	if mapCount == 0 {
		fullURL = mapURL
		mapCount++
	} else {
		offset := mapCount * 20
		fullURL = fmt.Sprintf("%s?offset=%d", mapURL, offset)
		mapCount++
	}
	// Check for cache entry
	_, exists := cache.Entry[fullURL]
	if exists {
		fmt.Println("Cache exists")

	}

	fmt.Printf("Fetching from: %v\n", fullURL)
	res, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var result locationAreaResponse
	err = decoder.Decode(&result)
	if err != nil {
		return err
	}

	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func CommandMapB() error {
	if mapCount == 1 {
		fmt.Println("you're on the first page")
		return nil
	}
	mapCount--
	mapCount--
	CommandMap()
	return nil
}
