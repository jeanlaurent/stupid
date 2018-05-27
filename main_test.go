package main

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestRemove(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("file", ""),
		fs.WithFile("remaining-file", ""),
		fs.WithDir("remaining-dir"))
	defer rootDirectory.Remove()

	err := remove([]string{
		filepath.Join(rootDirectory.Path(), "empty-dir"),
		filepath.Join(rootDirectory.Path(), "full-dir"),
		filepath.Join(rootDirectory.Path(), "non-existing"),
		filepath.Join(rootDirectory.Path(), "file"),
	})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("remaining-file", ""),
		fs.WithDir("remaining-dir"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo\n"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "foo.txt"), filepath.Join(rootDirectory.Path(), "bar.txt"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo\n"),
		fs.WithFile("bar.txt", "foo\n"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileInNonExistingPath(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo\n"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "foo.txt"), filepath.Join(rootDirectory.Path(), "bar", "bar.txt"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo\n"),
		fs.WithDir("bar", fs.WithMode(0700),
			fs.WithFile("bar.txt", "foo\n")),
	)
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTree(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n")),
		),
	)
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "source"), filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n")),
		),
		fs.WithDir("destination",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n")),
		),
	)

	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
