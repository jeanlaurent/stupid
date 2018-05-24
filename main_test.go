package main

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"

	"gotest.tools/fs"
)

// fs.WithMode(os.FileMode(0700)),
// fs.WithFile("file1", "content\n")),

func TestItWorks(t *testing.T) {
	directoryName := fs.NewDir(t, "stupid-test", fs.WithDir("toBeDeleted"), fs.WithDir("remaining"))
	defer directoryName.Remove()

	remove(filepath.Join(directoryName.Path(), "toBeDeleted"))

	expected := fs.Expected(t, fs.WithDir("remaining"))

	assert.Assert(t, fs.Equal(directoryName.Path(), expected))
}
