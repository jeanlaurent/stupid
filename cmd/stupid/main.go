package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	args := os.Args[1:]
	action := args[0]
	var err error
	switch action {
	case "help":
	case "--help":
		printUsage()
	case "cp":
		checkArguments(args, 3)
		err = copy(args[1:len(args)-1], args[len(args)-1])
	case "date":
		fmt.Print(time.Now().Format(time.RFC3339))
	case "home":
		home, err := homedir.Dir()
		if err == nil {
			fmt.Print(home)
		}
	case "mkdir":
		checkArguments(args, 2)
		err = mkDir(args[1:])
	case "rm":
		checkArguments(args, 2)
		err = remove(args[1:])
	case "silence":
		err = silence()
	case "tar":
		checkArguments(args, 3)
		err = tarFiles(args[len(args)-1], args[1:len(args)-1]...)
	case "untar":
		checkArguments(args, 3)
		err = untar(args[1], args[2])
	default:
		fmt.Fprintln(os.Stderr, "I don't know what", action, "means")
		printUsage()
		os.Exit(-3)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func checkArguments(args []string, max int) {
	if len(args) < max {
		fmt.Fprintln(os.Stderr, "Not enough arguments, I'm the stupid one, you fix it")
		printUsage()
		os.Exit(-1)
	}
}

func printUsage() {
	fmt.Println("I'm stupidly manipulating files and directories")
	fmt.Println("* stupid cp SRCS DST")
	fmt.Println("* stupid date")
	fmt.Println("* stupid home")
	fmt.Println("* stupid rm SRCS")
	fmt.Println("* stupid silence")
	fmt.Println("* stupid tar SRCS DST")
	fmt.Println("* stupid untar SRC DST")
}

func glob(sources []string, fail bool) ([]string, error) {
	var paths []string
	for _, source := range sources {
		var err error
		source, err = expand(source)
		if err != nil {
			return nil, err
		}
		matches, err := filepath.Glob(source)
		if err != nil {
			return nil, err
		}
		if matches != nil {
			paths = append(paths, matches...)
		} else if fail {
			return nil, fmt.Errorf("Source [%v] does not exist", source)
		} else {
			fmt.Printf("Source [%v] does not exist, doing nothing\n", source)
		}
	}
	if fail && len(paths) == 0 {
		return nil, fmt.Errorf("No source files")
	}
	return paths, nil
}
