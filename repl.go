package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wfcornelissen/pokedex/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
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

type locationAreaExplore struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int           `json:"min_level"`
				MaxLevel        int           `json:"max_level"`
				ConditionValues []interface{} `json:"condition_values"`
				Chance          int           `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type PokemonResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order int `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string      `json:"back_default"`
		BackFemale       interface{} `json:"back_female"`
		BackShiny        string      `json:"back_shiny"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		FrontDefault     string      `json:"front_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontShiny       string      `json:"front_shiny"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       interface{} `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"crystal"`
				Gold struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"gold"`
				Silver struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string      `json:"back_default"`
						BackFemale       interface{} `json:"back_female"`
						BackShiny        string      `json:"back_shiny"`
						BackShinyFemale  interface{} `json:"back_shiny_female"`
						FrontDefault     string      `json:"front_default"`
						FrontFemale      interface{} `json:"front_female"`
						FrontShiny       string      `json:"front_shiny"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	PastAbilities []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
	} `json:"past_abilities"`
}

var library map[string]PokemonResponse
var registry map[string]cliCommand
var mapCount int
var cache *internal.Cache

func init() {
	mapCount = 0
	cache = internal.NewCache(5 * time.Minute)
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
		"explore": {
			name: "explore",
			description: `Allows the user to explore an area by name or id, displaying
			a list of encounterable pokemon in the selected area.`,
			callback: CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: `Attempts to catch specified pokemon.`,
			callback:    CommandCatch,
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

func CheckInput(command []string) {
	opt := ""

	if len(command) == 1 {
		opt = ""
	} else {
		opt = CleanInput(command[1])[0]
	}
	for _, regCmd := range registry {
		if command[0] == regCmd.name {
			err := regCmd.callback(opt)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func CommandExit(option string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(option string) error {
	for _, regCmd := range registry {
		fmt.Printf("Name: %v\n", regCmd.name)
		fmt.Printf("Description: %v\n", regCmd.description)
		fmt.Println("----------------")
	}
	return nil
}

func CommandMap(option string) error {
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
	cacheRes, exists := cache.Get(fullURL)
	if !exists {
		fmt.Println("Cache does not exist.")
		fmt.Printf("Fetching from: %v\n", fullURL)
		res, err := http.Get(fullURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// Read the response body into bytes
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body: %v", err)
		}

		// Cache the raw response body bytes
		cache.Add(fullURL, bodyBytes)

		// Decode the JSON from the bytes we just read
		var result locationAreaResponse
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return err
		}

		for _, location := range result.Results {
			fmt.Println(location.Name)
		}

		return nil
	}

	// Cache hit - unmarshal from cached bytes
	var result locationAreaResponse
	err := json.Unmarshal(cacheRes, &result)
	if err != nil {
		return fmt.Errorf("Error unmarshaling cache: %v", err)
	}

	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func CommandMapB(option string) error {
	if mapCount == 1 {
		fmt.Println("you're on the first page")
		return nil
	}
	mapCount--
	mapCount--
	CommandMap(option)
	return nil
}

func CommandExplore(option string) error {
	exploreURL := "https://pokeapi.co/api/v2/location-area/"
	if option == "" {
		return fmt.Errorf("No option supplied")
	}

	fullURL := exploreURL + option
	fmt.Printf("Exploring %v...\n", option)
	cacheRes, exists := cache.Get(fullURL)
	if !exists {

		res, err := http.Get(fullURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body: %v", err)
		}

		// Cache the raw response body bytes
		cache.Add(fullURL, bodyBytes)

		// Decode the JSON from the bytes we just read
		var result locationAreaExplore
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return err
		}

		for _, encounters := range result.PokemonEncounters {
			fmt.Println(encounters.Pokemon.Name)
		}
		return nil
	}
	var result locationAreaExplore
	err := json.Unmarshal(cacheRes, &result)
	if err != nil {
		return err
	}

	for _, encounters := range result.PokemonEncounters {
		fmt.Println(encounters.Pokemon.Name)
	}

	return nil
}

func CommandCatch(option string) error {
	pokeURL := "https://pokeapi.co/api/v2/pokemon/"
	fullURL := pokeURL + option

	pokeByte, exists := cache.Get(fullURL)
	if !exists {
		res, err := http.Get(fullURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body: %v", err)
		}

		// Cache the raw response body bytes
		cache.Add(fullURL, bodyBytes)

		// Decode the JSON from the bytes we just read
		var result PokemonResponse
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return err
		}
		if CatchPokemon(&result) {
			library[fullURL] = result
		}
		return nil
	}
	var result PokemonResponse
	err := json.Unmarshal(pokeByte, &result)
	if err != nil {
		return err
	}
	if CatchPokemon(&result) {
		library[fullURL] = result
	}

	return nil
}

func CatchPokemon(poke *PokemonResponse) bool {
	fmt.Printf("Throwing a Pokeball at %v...\n", poke.Name)
	baseExp := poke.BaseExperience
	chance := 100 - ((-0.11836734693 * float64(baseExp)) + 102.367346938)

	r := rand.Float64()

	val := r * 100

	if val > chance {
		fmt.Printf("%v was caught!\n", poke.Name)
		return true
	}
	fmt.Printf("%v escaped!\n", poke.Name)
	return false
}
