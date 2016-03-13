package rkt

import (
	"os"

	"github.com/giantswarm/yochu/fetchclient"
	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

type rktOptions struct {
	UseOverlay bool
}

var (
	vLogger = func(f string, v ...interface{}) {}

	rktGarbageServiceName  = "rkt-gc.service"
	rktGarbageTimerName    = "rkt-gc.timer"
	rktMetadataServiceName = "rkt-metadata.service"
	rktMetadataSocketName  = "rkt-metadata.socket"

	rktSystemdPath             = "/etc/systemd/system/"
	rktGarbageServiceTemplate  = "templates/rkt-gc.service.tmpl"
	rktGarbageTimerTemplate    = "templates/rkt-gc.timer.tmpl"
	rktMetadataServiceTemplate = "templates/rkt-metadata.service.tmpl"
	rktMetadataSocketTemplate  = "templates/rkt-metadata.socket.tmpl"

	units = []string{rktGarbageServiceName, rktGarbageTimerName, rktMetadataServiceName, rktMetadataSocketName}

	binaryFileMode = os.FileMode(0755)
	unitFileMode   = os.FileMode(0644)
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, sc *systemd.SystemdClient, fc fetchclient.FetchClient, distributionPath, rktVersion string, startDaemon, useOverlay bool) error {
	vLogger("\n# call rkt.Setup()")

	rktRaw, err := fc.Get("rkt/" + rktVersion + "/rkt")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/rkt", rktRaw, binaryFileMode); err != nil {
		return maskAny(err)
	}

	err = createRktGarbageService(fsc, useOverlay)
	if err != nil {
		return maskAny(err)
	}

	err = createRktMetadataService(fsc, useOverlay)
	if err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if startDaemon {
		if err := sc.Start(rktGarbageTimerName); err != nil {
			return maskAny(err)
		}
		if err := sc.Start(rktMetadataServiceName); err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func createRktGarbageService(fsc *fs.FsClient, useOverlay bool) error {
	opts := rktOptions{
		UseOverlay: useOverlay,
	}

	b, err := templates.Render(rktGarbageServiceTemplate, opts)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(rktSystemdPath+rktGarbageServiceName, b.Bytes(), unitFileMode); err != nil {
		return maskAny(err)
	}

	rktGarbageTimer, err := templates.Asset(rktGarbageTimerTemplate)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(rktSystemdPath+rktGarbageTimerName, rktGarbageTimer, unitFileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func createRktMetadataService(fsc *fs.FsClient, useOverlay bool) error {
	opts := rktOptions{
		UseOverlay: useOverlay,
	}

	b, err := templates.Render(rktMetadataServiceTemplate, opts)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(rktSystemdPath+rktMetadataServiceName, b.Bytes(), unitFileMode); err != nil {
		return maskAny(err)
	}

	rktMetadataSocket, err := templates.Asset(rktMetadataSocketTemplate)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(rktSystemdPath+rktMetadataSocketName, rktMetadataSocket, unitFileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func Teardown(fsc *fs.FsClient, sc *systemd.SystemdClient, distributionPath string, stopDaemon bool) error {
	vLogger("\n# call rkt.Teardown()")

	for _, u := range units {
		exists, err := sc.Exists(u)
		if err != nil {
			return maskAny(err)
		}

		if !exists || !stopDaemon {
			continue
		}

		if err := sc.Stop(u); err != nil {
			return maskAny(err)
		}

		if err := fsc.Remove(rktSystemdPath + u); err != nil {
			return maskAny(err)
		}
	}

	// reload unit files, that is, `systemctl daemon-reload`
	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if err := fsc.Remove(distributionPath + "/rkt"); err != nil {
		return maskAny(err)
	}

	return nil
}
