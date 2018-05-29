package main

import (
	"os"
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

func TestRemoveWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := remove([]string{
		filepath.Join(rootDirectory.Path(), "*-dir"),
		filepath.Join(rootDirectory.Path(), "full-dir"),
	})
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

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

func TestCopyMultipleFilesToFileFails(t *testing.T) {
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

func TestCopyFileToNonExistingPath(t *testing.T) {
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

func TestTarTree(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	err := tarFiles(dst, filepath.Join(rootDirectory.Path(), "source", "foo.txt"), filepath.Join(rootDirectory.Path(), "source", "bar"))
	assert.NilError(t, err)
	err = untar(dst, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))),
		fs.WithDir("destination",
			fs.WithFile("dst.tar", "", fs.MatchAnyFileContent),
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar")),
		),
	)
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
				fs.WithFile("bar.txt", "bar")),
		),
	)

	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestGzipTarTree(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar.gz")
	err := tarFiles(dst, filepath.Join(rootDirectory.Path(), "source", "foo.txt"), filepath.Join(rootDirectory.Path(), "source", "bar"))
	assert.NilError(t, err)
	err = untar(dst, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))),
		fs.WithDir("destination",
			fs.WithFile("dst.tar.gz", "", fs.MatchAnyFileContent),
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
