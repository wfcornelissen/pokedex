package main

func commandMapb() error {
	err := apiRequest()
	if err != nil {
		return err
	}
	return nil
}
