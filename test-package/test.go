package test_package

import "fmt"

func Something() error {
	err := GetError()

	if err != nil {
		return fmt.Errorf("something went wrong")
	}

	return nil
}

func GetError() error {
	return fmt.Errorf("testing error")
}
