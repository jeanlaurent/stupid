package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func tarFiles(dst string, srcs ...string) error {
	srcs, err := glob(srcs, true)
	if err != nil {
		return err
	}
	var w io.Writer
	if dst == "-" {
		w = os.Stdout
	} else {
		dst, err = expand(dst)
		if err != nil {
			return err
		}
		fmt.Printf("Taring [%v]\n", dst)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
		ext := filepath.Ext(dst)
		if ext == ".gz" || ext == ".tgz" {
			gz := gzip.NewWriter(w)
			defer gz.Close()
			w = gz
		}
	}
	tw := tar.NewWriter(w)
	defer tw.Close()
	for _, src := range srcs {
		dir := filepath.Dir(src)
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
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
	var r io.Reader
	if src == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
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
