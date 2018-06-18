package main

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestMkDir(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("existing-dir"),
		fs.WithDir("existing-sub-dir"))
	defer rootDirectory.Remove()

	err := mkDir([]string{
		filepath.Join(rootDirectory.Path(), "some-dir"),
		filepath.Join(rootDirectory.Path(), "existing-dir"),
		filepath.Join(rootDirectory.Path(), "sub-dir", "another-dir"),
		filepath.Join(rootDirectory.Path(), "existing-sub-dir", "another-dir"),
	})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("some-dir"),
		fs.WithDir("existing-dir"),
		fs.WithDir("sub-dir",
			fs.WithDir("another-dir")),
		fs.WithDir("existing-sub-dir",
			fs.WithDir("another-dir")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
