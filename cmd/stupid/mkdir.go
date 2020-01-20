package main

import (
	"fmt"
	"os"
)

func mkDir(sources []string) error {
	for _, source := range sources {
		var err error
		source, err = expand(source)
		if err != nil {
			return err
		}
		fmt.Printf("Creating [%v]\n", source)
		if err := os.MkdirAll(source, 0755); err != nil {
			return err
		}
	}
	return nil
}
