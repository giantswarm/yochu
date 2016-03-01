package k8s

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

func Setup(fsc *fs.FsClient, fc fetchclient.FetchClient, distributionPath, k8sVersion string) error {
	vLogger("\n# call k8s.Setup()")

	k8sRaw, err := fc.Get("k8s/" + k8sVersion + "/kubectl")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/kubectl", k8sRaw, fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func Teardown(fsc *fs.FsClient, distributionPath string) error {
	vLogger("\n# call k8s.Teardown()")

	if err := fsc.Remove(distributionPath + "/kubectl"); err != nil {
		return maskAny(err)
	}

	return nil
}
