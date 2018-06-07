package main

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestCopyFileToNonExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "foo.txt")}, filepath.Join(rootDirectory.Path(), "bar.txt"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "foo"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileToExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "foo.txt")}, filepath.Join(rootDirectory.Path(), "bar.txt"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "foo"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileToNonExistingDir(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "foo.txt")}, filepath.Join(rootDirectory.Path(), "destination")+"/")
	assert.NilError(t, err)

	info, err := os.Stat(rootDirectory.Path())
	assert.NilError(t, err)
	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("destination",
			fs.WithMode(info.Mode()),
			fs.WithFile("foo.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileToExistingDir(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("destination"))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "foo.txt")}, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("destination",
			fs.WithFile("foo.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyMultipleFilesToFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	err := copy(
		[]string{
			filepath.Join(rootDirectory.Path(), "foo.txt"),
			filepath.Join(rootDirectory.Path(), "foo.txt"),
		}, filepath.Join(rootDirectory.Path(), "bar.txt"))
	assert.Error(t, err, "Only one source file allowed when destination is a file")
}

func TestCopyMultipleFilesWithNonExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing")
	err := copy(
		[]string{
			src,
			filepath.Join(rootDirectory.Path(), "foo.txt"),
		}, filepath.Join(rootDirectory.Path(), "destination"))
	assert.Error(t, err, "Source ["+src+"] does not exist")
}

func TestCopyMultipleFilesWithNonExistingGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing*")
	err := copy(
		[]string{
			src,
			filepath.Join(rootDirectory.Path(), "foo.txt"),
		}, filepath.Join(rootDirectory.Path(), "destination"))
	assert.Error(t, err, "Source ["+src+"] does not exist")
}

func TestCopyNonExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root")
	defer rootDirectory.Remove()

	err := copy(
		[]string{filepath.Join(rootDirectory.Path(), "non-existing")}, filepath.Join(rootDirectory.Path(), "bar.txt"))
	assert.ErrorContains(t, err, "does not exist")
}

func TestCopyFileToNonExistingPathFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "foo.txt")}, filepath.Join(rootDirectory.Path(), "bar", "bar.txt"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("bar", fs.WithMode(0700),
			fs.WithFile("bar.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileOverItself(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	f := filepath.Join(rootDirectory.Path(), "foo.txt")
	err := copy([]string{f}, f)
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithFile("foo.txt", "foo"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeToNonExistingPath(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "source", "bar")}, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))),
		fs.WithDir("destination",
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	err := copy([]string{filepath.Join(rootDirectory.Path(), "source", "*")}, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar")),
		),
		fs.WithDir("destination",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeWithEmptyGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo")))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing*")
	err := copy([]string{src}, filepath.Join(rootDirectory.Path(), "destination"))
	assert.Error(t, err, "Source ["+src+"] does not exist")

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
