package main

import "fmt"

func commandMap(config *config) error {
	fmt.Println("Location areas:")
	err := apiRequest(config.Next)
	if err != nil {
		return err
	}

	return nil
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
