package rkt

import (
	"os"

	"github.com/giantswarm/yochu/fetchclient"
	"github.com/giantswarm/yochu/fs"
)

var (
	vLogger = func(f string, v ...interface{}) {}

	fileMode = os.FileMode(0755)
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, fc fetchclient.FetchClient, distributionPath, rktVersion string) error {
	vLogger("\n# call rkt.Setup()")

	rktRaw, err := fc.Get("rkt/" + rktVersion + "/rkt")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/rkt", rktRaw, fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func Teardown(fsc *fs.FsClient, distributionPath string) error {
	vLogger("\n# call rkt.Teardown()")

	if err := fsc.Remove(distributionPath + "/rkt"); err != nil {
		return maskAny(err)
	}

	return nil
}
