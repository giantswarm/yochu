// Package fs provides a client for writing and removing files on the local filesystem.
package fs

import (
	"io/ioutil"
	"os"
	"path"
)

var vLogger = func(f string, v ...interface{}) {}

// Configure sets the logger for this package.
func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

// FsClient is a filesystem client.
type FsClient struct{}

// NewFsClient returns a new FsClient.
func NewFsClient() (*FsClient, error) {
	return &FsClient{}, nil
}

// Write creates any necessary parent directories, and then writes the file to destination.
// See os.MkdirAll and ioutil.WriteFile.
func (fsc *FsClient) Write(destination string, raw []byte, perm os.FileMode) error {
	vLogger("  call FsClient.Write(path, data): %s, <blob:%d>", destination, len(raw))

	parent := path.Dir(destination)
	if err := os.MkdirAll(parent, perm); err != nil && !os.IsExist(err) {
		return Mask(err)
	}

	if err := ioutil.WriteFile(destination, raw, perm); err != nil {
		return Mask(err)
	}

	return nil
}

// MkdirAll creates all folders (including the last) for the given path.
func (fsc *FsClient) MkdirAll(path string, perm os.FileMode) error {
	vLogger("  call FsClient.MkdirAll(path, perm): %s", path)

	if err := os.MkdirAll(path, perm); err != nil {
		return Mask(err)
	}

	return nil
}

// Symlink creates a symlink at newname pointing at oldname.
func (fsc *FsClient) Symlink(oldname, newname string) error {
	vLogger("  call FsClient.Symlink(oldname, newname): %s, %s", oldname, newname)

	if err := os.Symlink(oldname, newname); err != nil {
		return Mask(err)
	}

	return nil
}

// Remove removes the destination path, and any children.
// See os.RemoveAll.
func (fsc *FsClient) Remove(destination string) error {
	vLogger("  call FsClient.Remove(path): %s", destination)

	if err := os.RemoveAll(destination); err != nil {
		return Mask(err)
	}

	return nil
}
