package main

import (
	"fmt"
	"os"
)

func remove(sources []string) error {
	sources, err := glob(sources, false)
	if err != nil {
		return err
	}
	for _, source := range sources {
		fmt.Printf("Removing [%v]\n", source)
		if err = os.RemoveAll(source); err != nil {
			return err
		}
	}
	return nil
}
