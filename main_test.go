package main

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"

	"gotest.tools/fs"
)

func TestRemoveAnEmptyDirectory(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("toBeDeleted"),
		fs.WithDir("remaining"))
	defer rootDirectory.Remove()

	remove(filepath.Join(rootDirectory.Path(), "toBeDeleted"))

	expected := fs.Expected(t, fs.WithDir("remaining"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveAFullDirectory(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("toBeDeleted",
			fs.WithFile("foo", "foobar")),
		fs.WithDir("remaining"))
	defer rootDirectory.Remove()

	remove(filepath.Join(rootDirectory.Path(), "toBeDeleted"))

	expected := fs.Expected(t, fs.WithDir("remaining"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveNonExistingDirectory(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("remaining"))
	defer rootDirectory.Remove()

	remove(filepath.Join(rootDirectory.Path(), "nonExisting"))

	expected := fs.Expected(t, fs.WithDir("remaining"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
func TestCopyFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo\n"))
	defer rootDirectory.Remove()

	copy(filepath.Join(rootDirectory.Path(), "foo.txt"), filepath.Join(rootDirectory.Path(), "bar.txt"))

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo\n"),
		fs.WithFile("bar.txt", "foo\n"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileInNonExistingPath(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo\n"))
	defer rootDirectory.Remove()

	copy(filepath.Join(rootDirectory.Path(), "foo.txt"), filepath.Join(rootDirectory.Path(), "bar", "bar.txt"))

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

	copy(filepath.Join(rootDirectory.Path(), "source"), filepath.Join(rootDirectory.Path(), "destination"))

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
