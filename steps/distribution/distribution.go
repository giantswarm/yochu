package distribution

import (
	"os"

	"github.com/giantswarm/yochu/fs"
)

var (
	vLogger = func(f string, v ...interface{}) {}

	fileMode = os.FileMode(0755)
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, distributionPath string) error {
	vLogger("\n# call distribution.Setup()")

	if err := fsc.MkdirAll(distributionPath, fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func Teardown(fsc *fs.FsClient, distributionPath string) error {
	vLogger("\n# call distribution.Teardown()")

	if err := fsc.Remove(distributionPath); err != nil {
		return maskAny(err)
	}

	return nil
}
