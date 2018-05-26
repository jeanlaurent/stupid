package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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
	case "home":
		home, err := homedir.Dir()
		if err == nil {
			fmt.Print(home)
		}
	case "rm":
		checkArguments(args, 2)
		err = remove(args[1:])
	case "cp":
		checkArguments(args, 3)
		err = copy(args[1:len(args)-1], args[len(args)-1])
	case "tar":
		checkArguments(args, 2)
		err = tarFiles(args[len(args)-1], args[1:len(args)-1]...)
	case "untar":
		checkArguments(args, 2)
		err = untar(args[1], args[2])
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
	fmt.Println("I'm stupidly manipulating files and directories")
	fmt.Println("* stupid home")
	fmt.Println("* stupid cp SRCS DST")
	fmt.Println("* stupid rm SRCS")
	fmt.Println("* stupid tar SRCS DST")
	fmt.Println("* stupid untar SRC DST")
}

func remove(sources []string) error {
	sources, err := glob(sources)
	if err != nil {
		return err
	}
	for _, source := range sources {
		fmt.Println("Removing", source)
		_, err := os.Stat(source)
		if os.IsNotExist(err) {
			fmt.Printf("Source [%v] does not exist, doing nothing\n", source)
			continue
		}
		if err != nil {
			return err
		}
		if err = os.RemoveAll(source); err != nil {
			return err
		}
	}
	return nil
}

func glob(sources []string) ([]string, error) {
	var paths []string
	for _, source := range sources {
		matches, err := filepath.Glob(source)
		if err != nil {
			return nil, err
		}
		paths = append(paths, matches...)
	}
	return paths, nil
}

func copy(sources []string, destination string) error {
	sources, err := glob(sources)
	if err != nil {
		return err
	}
	toFile := true
	info, err := os.Stat(destination)
	if os.IsNotExist(err) {
		if destination[len(destination)-1] == '/' || len(sources) > 1 {
			toFile = false
		}
	} else if err != nil {
		return err
	} else if info.IsDir() {
		toFile = false
	} else if len(sources) > 1 {
		return fmt.Errorf("Only one source file allowed when destination is a file")
	}
	for _, source := range sources {
		info, err = os.Stat(source)
		if os.IsNotExist(err) {
			return fmt.Errorf("Source [%v] does not exist", source)
		}
		if err != nil {
			return err
		}
		dest := filepath.Join(destination, filepath.Base(source))
		if info.IsDir() {
			fmt.Printf("Copying dir [%v] to [%v]\n", source, dest)
			if err = copyDirectory(source, dest, info.Mode()); err != nil {
				return err
			}
			continue
		}
		if toFile {
			dest = destination
		}
		dirInfo, err := os.Stat(filepath.Dir(source))
		if err != nil {
			return err
		}
		if err = os.MkdirAll(filepath.Dir(dest), dirInfo.Mode()); err != nil {
			return err
		}
		fmt.Printf("Copying file [%v] to [%v]\n", source, dest)
		if err = copyFile(source, dest, info.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
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
	return os.Chmod(dst, mode)
}

func copyDirectory(src string, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(dst, mode); err != nil {
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
			if err = copyDirectory(srcfp, dstfp, info.Mode()); err != nil {
				return err
			}
		} else if err = copyFile(srcfp, dstfp, info.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func tarFiles(dst string, srcs ...string) error {
	fmt.Printf("Taring [%v]\n", dst)
	if err := os.MkdirAll(filepath.Dir(dst), os.ModeDir); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	var w io.Writer
	w = f
	ext := filepath.Ext(dst)
	if ext == ".gz" || ext == ".tgz" {
		gz := gzip.NewWriter(w)
		defer gz.Close()
		w = gz
	}
	tw := tar.NewWriter(w)
	defer tw.Close()
	for _, src := range srcs {
		dir := filepath.Dir(src)
		err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			hdr, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}
			hdr.Name = rel
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(tw, f)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func untar(src, dst string) error {
	fmt.Printf("Untaring [%v] to [%v]\n", src, dst)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := os.MkdirAll(dst, os.ModeDir); err != nil {
		return err
	}
	var r io.Reader
	r = f
	ext := filepath.Ext(src)
	if ext == ".gz" || ext == ".tgz" {
		gz, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer gz.Close()
		r = gz
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		path := filepath.Join(dst, hdr.Name)
		info := hdr.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, tr)
		if err != nil {
			return err
		}
	}
	return nil
}
