package main

func commandMap(config *config) error {
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
