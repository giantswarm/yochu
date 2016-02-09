package overlay

import (
	"os"

	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

var (
	vLogger = func(f string, v ...interface{}) {}

	overlayMount         = "usr-bin.mount"
	overlayMountPath     = "/etc/systemd/system/usr-bin.mount"
	overlayMountTemplate = "templates/usr-bin.mount.tmpl"

	fileMode = os.FileMode(0755)
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, sc *systemd.SystemdClient, distributionPath, overlayWorkdir, overlayMountPoint string) error {
	vLogger("\n# call overlay.Setup()")

	if err := fsc.MkdirAll(overlayWorkdir, fileMode); err != nil {
		return maskAny(err)
	}

	opts := struct {
		OverlayUpperdir string
		OverlayWorkdir  string
		MountPoint      string
	}{
		distributionPath,
		overlayWorkdir,
		overlayMountPoint,
	}

	b, err := templates.Render(overlayMountTemplate, opts)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(overlayMountPath, b.Bytes(), fileMode); err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if err := sc.Start(overlayMount); err != nil {
		return maskAny(err)
	}

	return nil
}

func Teardown(fsc *fs.FsClient, sc *systemd.SystemdClient, distributionPath, overlayWorkdir, overlayMountPoint string) error {
	vLogger("\n# call overlay.Teardown()")

	exists, err := sc.Exists(overlayMount)
	if err != nil {
		return maskAny(err)
	}

	if exists {
		if err := sc.Stop(overlayMount); err != nil {
			return maskAny(err)
		}
	}

	if err := fsc.Remove(overlayMountPath); err != nil && !fs.IsNotExist(err) {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if err := fsc.Remove(overlayWorkdir); err != nil && !fs.IsNotExist(err) {
		return maskAny(err)
	}

	return nil
}
