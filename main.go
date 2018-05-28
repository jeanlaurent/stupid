package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func main() {
	args := os.Args[1:]
	action := args[0]

	switch action {
	case "help":
		printUsage()
	case "home":
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-3)
		}
		fmt.Print(home)
	case "rm":
		checkArgumentsOrQuit(args, 2)
		remove(args[1])
	case "cp":
		checkArgumentsOrQuit(args, 3)
		copy(args[1], args[2])
	default:
		fmt.Printf("I don't what %v means\n", action)
		printUsage()
		os.Exit(-3)
	}
}

func remove(source string) {
	fmt.Printf("removing %v\n", source)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		fmt.Printf("source [%v] does not exist, doing nothing\n", source)
		return
	}
	err = os.RemoveAll(source)
	if err != nil {
		fmt.Println(err)
	}
}

func copy(source, destination string) {
	sourceFile, err := os.Stat(source)
	if os.IsNotExist(err) {
		fmt.Printf("source [%v] does not exist\n", source)
		os.Exit(-2)
	}
	if sourceFile.Mode().IsDir() {
		err := copyDirectory(source, destination)
		if err != nil {
			fmt.Println(err)
			os.Exit(-3)
		}
	} else {
		destinationPath := filepath.Dir(destination)
		_, err := os.Stat(destinationPath)
		if os.IsNotExist(err) {
			sourcePath := filepath.Dir(source)
			sourceDir, err := os.Stat(sourcePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(-3)
			}
			os.MkdirAll(destinationPath, sourceDir.Mode())
		}
		err = copyFile(source, destination)
		if err != nil {
			fmt.Println(err)
			os.Exit(-3)
		}
	}
}

func copyFile(src, dst string) error {
	var err error
	var sourceFD *os.File
	var destinationFD *os.File
	var sourceFileInfo os.FileInfo

	fmt.Printf("Copy file [%v] to [%v]\n", src, dst)
	if sourceFD, err = os.Open(src); err != nil {
		return err
	}
	defer sourceFD.Close()

	if destinationFD, err = os.Create(dst); err != nil {
		return err
	}
	defer destinationFD.Close()

	if _, err = io.Copy(destinationFD, sourceFD); err != nil {
		return err
	}
	if sourceFileInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, sourceFileInfo.Mode())
}

func copyDirectory(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	fmt.Printf("Copy dir [%v] to [%v]\n", src, dst)
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDirectory(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func checkArgumentsOrQuit(args []string, max int) {
	if len(args) < max {
		fmt.Println("not enough argument, I'm the stupid one, you fix it.")
		printUsage()
		os.Exit(-1)
	}
}

func printUsage() {
	fmt.Println("I'm stupidly copying or removing files.")
	fmt.Println("* stupid home")
	fmt.Println("* stupid cp src dst")
	fmt.Println("* stupid rm src")
}
