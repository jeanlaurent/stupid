package main

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

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

func TestTarTreeWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	err := tarFiles(dst, filepath.Join(rootDirectory.Path(), "source", "*"))
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

func TestTarTreeWithEmptyGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo\n"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar\n"))))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing*")
	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	err := tarFiles(dst, src)
	assert.Error(t, err, "Source ["+src+"] does not exist")
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
