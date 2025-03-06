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

// Offset tracks pagination for the API requests
var LocationResp config

func apiRequest() error {
	var URL string
	if LocationResp.Next == "" {
		URL = apiBaseURL
	} else {
		URL = LocationResp.Next
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

	err = json.Unmarshal(body, &LocationResp)
	if err != nil {
		return err
	}
	fmt.Println(LocationResp.Next)
	fmt.Println("Location areas:")
	for _, area := range LocationResp.Results {
		fmt.Printf("- %s\n", area.Name)
	}
	return nil
}
