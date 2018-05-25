package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func main() {
	args := os.Args[1:]
	action := args[0]
	var err error
	switch action {
	case "help":
	case "--help":
		printUsage()
	case "home":
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-3)
		}
		fmt.Print(home)
	case "rm":
		checkArguments(args, 2)
		err = remove(args[1])
	case "cp":
		checkArguments(args, 3)
		err = copy(args[1], args[2])
	default:
		fmt.Printf("I don't know what %v means\n", action)
		printUsage()
		os.Exit(-3)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func checkArguments(args []string, max int) {
	if len(args) < max {
		fmt.Println("Not enough arguments, I'm the stupid one, you fix it")
		printUsage()
		os.Exit(-1)
	}
}

func printUsage() {
	fmt.Println("I'm stupidly copying or removing files and directories")
	fmt.Println("* stupid home")
	fmt.Println("* stupid cp SRC DST")
	fmt.Println("* stupid rm SRC")
}

func remove(source string) error {
	fmt.Println("Removing", source)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		fmt.Printf("Source [%v] does not exist, doing nothing\n", source)
		return nil
	}
	return os.RemoveAll(source)
}

func copy(source, destination string) error {
	sourceInfo, err := os.Stat(source)
	if os.IsNotExist(err) {
		return fmt.Errorf("Source [%v] does not exist", source)
	}
	if err != nil {
		return err
	}
	if sourceInfo.Mode().IsDir() {
		return copyDirectory(source, destination)
	}
	sourceDirInfo, err := os.Stat(filepath.Dir(source))
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(destination), sourceDirInfo.Mode())
	if err != nil {
		return err
	}
	return copyFile(source, destination)
}

func copyFile(src, dst string) error {
	fmt.Printf("Copying file [%v] to [%v]\n", src, dst)
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if _, err = io.Copy(destination, source); err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

func copyDirectory(src string, dst string) error {
	fmt.Printf("Copying dir [%v] to [%v]\n", src, dst)
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, info := range infos {
		srcfp := filepath.Join(src, info.Name())
		dstfp := filepath.Join(dst, info.Name())
		if info.IsDir() {
			if err = copyDirectory(srcfp, dstfp); err != nil {
				return err
			}
		} else if err = copyFile(srcfp, dstfp); err != nil {
			return err
		}
	}
	return nil
}
