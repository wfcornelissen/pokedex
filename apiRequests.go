package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiBaseURL = "https://pokeapi.co/api/v2/location-area/"
)

func apiRequest(URL string) error {
	if URL == "" {
		URL = apiBaseURL
	}
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &Config)
	if err != nil {
		return err
	}

	for _, area := range Config.Results {
		fmt.Printf("- %s\n", area.Name)
	}
	return nil
}
