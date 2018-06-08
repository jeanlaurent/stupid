package main

import (
	"io"
	"io/ioutil"
	"os"
)

func silence() error {
	_, err := io.Copy(ioutil.Discard, os.Stdin)
	return err
}
