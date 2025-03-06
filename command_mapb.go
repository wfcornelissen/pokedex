package main

func commandMapb(config *config) error {
	err := apiRequest(config.Previous)
	if err != nil {
		return err
	}
	return nil
}
