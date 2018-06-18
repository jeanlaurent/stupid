package main

import (
	"fmt"
	"os"
)

func mkDir(sources []string) error {
	for _, source := range sources {
		fmt.Printf("Creating [%v]\n", source)
		if err := os.MkdirAll(source, 0755); err != nil {
			return err
		}
	}
	return nil
}
