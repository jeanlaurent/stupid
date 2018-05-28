package main

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestRemoveEmptyDirectory(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("toBeDeleted"),
		fs.WithDir("remaining"))
	defer rootDirectory.Remove()

	err := remove(filepath.Join(rootDirectory.Path(), "toBeDeleted"))
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithDir("remaining"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveFullDirectory(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("toBeDeleted",
			fs.WithFile("foo", "foobar")),
		fs.WithDir("remaining"))
	defer rootDirectory.Remove()

	err := remove(filepath.Join(rootDirectory.Path(), "toBeDeleted"))
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithDir("remaining"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveNonExistingDirectory(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("remaining"))
	defer rootDirectory.Remove()

	err := remove(filepath.Join(rootDirectory.Path(), "nonExisting"))
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithDir("remaining"))
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
