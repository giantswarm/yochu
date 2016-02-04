package fs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const (
	TempPrefix = "yochu-test"
)

// TestWrite tests the Write function
func TestWrite(t *testing.T) {
	Configure(t.Logf)

	fsClient, err := NewFsClient()
	if err != nil {
		t.Fatal("could not create fsClient: ", err)
	}

	tempDir, err := ioutil.TempDir("", TempPrefix)
	if err != nil {
		t.Fatal("could not create temporary directory: ", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := path.Join(tempDir, "temp-file")
	data := []byte("this is some test data")
	mode := os.FileMode(int(0644))

	if err := fsClient.Write(tempFile, data, mode); err != nil {
		t.Fatal("could not write file: ", err)
	}

	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		t.Fatal("could not stat new file: ", err)
	}
	if fileInfo.Mode() != mode {
		t.Fatal("incorrect mode set on file: ", fileInfo.Mode(), mode)
	}

	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatal("could not open file for reading: ", err)
	}
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal("could not read data from file: ", err)
	}

	if !bytes.Equal(data, fileData) {
		t.Fatal("read data does not match written data")
	}
}

// TestRemove tests the Remove function
func TestRemove(t *testing.T) {
	Configure(t.Logf)

	fsClient, err := NewFsClient()
	if err != nil {
		t.Fatal("could not create fsClient: ", err)
	}

	tempFile, err := ioutil.TempFile("", TempPrefix)
	if err != nil {
		t.Fatal("could not create temp file: ", err)
	}

	if err := fsClient.Remove(tempFile.Name()); err != nil {
		t.Fatal("could not remove temp file: ", err)
	}

	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		t.Fatal("could not determine file has been removed: ", err)
	}
}
