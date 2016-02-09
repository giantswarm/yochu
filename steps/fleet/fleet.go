package fleet

import (
	"os"

	"github.com/giantswarm/yochu/fetchclient"
	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

var (
	vLogger = func(f string, v ...interface{}) {}

	fleetServiceName = "fleet.service"

	fleetServicePath     = "/etc/systemd/system/fleet.service"
	fleetServiceTemplate = "templates/fleet.service.tmpl"

	fileMode = os.FileMode(0755)
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, sc *systemd.SystemdClient, fc fetchclient.FetchClient, distributionPath, fleetVersion string, startDaemon bool) error {
	vLogger("\n# call fleet.Setup()")

	fleetdRaw, err := fc.Get("fleet/" + fleetVersion + "/fleetd")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/fleet", fleetdRaw, fileMode); err != nil {
		return maskAny(err)
	}

	if err := fsc.Symlink(distributionPath+"/fleet", distributionPath+"/fleetd"); err != nil {
		return maskAny(err)
	}

	fleetctlRaw, err := fc.Get("fleet/" + fleetVersion + "/fleetctl")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/fleetctl", fleetctlRaw, fileMode); err != nil {
		return maskAny(err)
	}

	fleetServiceRaw, err := templates.Asset(fleetServiceTemplate)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(fleetServicePath, fleetServiceRaw, fileMode); err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if startDaemon {
		if err := sc.Start(fleetServiceName); err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func Teardown(fsc *fs.FsClient, sc *systemd.SystemdClient, distributionPath string, stopDaemon bool) error {
	vLogger("\n# call fleet.Teardown()")

	exists, err := sc.Exists(fleetServiceName)
	if err != nil {
		return maskAny(err)
	}

	if exists && stopDaemon {
		if err := sc.Stop(fleetServiceName); err != nil {
			return maskAny(err)
		}
	}

	if err := fsc.Remove(distributionPath + "/fleet"); err != nil {
		return maskAny(err)
	}

	if err := fsc.Remove(distributionPath + "/fleetd"); err != nil {
		return maskAny(err)
	}

	if err := fsc.Remove(distributionPath + "/fleetctl"); err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	return nil
}
