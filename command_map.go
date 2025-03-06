package main

func commandMap() error {
	err := apiRequest()
	if err != nil {
		return err
	}

	return nil
}

type config struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous interface{}    `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
