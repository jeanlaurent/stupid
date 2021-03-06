package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func copy(sources []string, destination string) error {
	sources, err := glob(sources, true)
	if err != nil {
		return err
	}
	toFile := true
	destination, err = expand(destination)
	if err != nil {
		return err
	}
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
	if same, err := sameFile(src, dst); err != nil || same {
		return err
	}
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

func sameFile(src, dst string) (bool, error) {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return false, err
	}
	absDst, err := filepath.Abs(dst)
	if err != nil {
		return false, err
	}
	return absSrc == absDst, nil
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
