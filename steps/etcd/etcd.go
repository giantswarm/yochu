package etcd

import (
	"os"

	"github.com/giantswarm/yochu/fetchclient"
	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

var (
	vLogger = func(f string, v ...interface{}) {}

	etcdServiceName = "etcd2.service"

	etcdServicePath     = "/etc/systemd/system/etcd2.service"
	etcdServiceTemplate = "templates/etcd2.service.tmpl"

	fileMode = os.FileMode(0755)
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, sc *systemd.SystemdClient, fc fetchclient.FetchClient, distributionPath, etcdVersion string, startDaemon bool) error {
	vLogger("\n# call etcd.Setup()")

	etcdRaw, err := fc.Get("etcd/" + etcdVersion + "/etcd")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/etcd2", etcdRaw, fileMode); err != nil {
		return maskAny(err)
	}

	etcdctlRaw, err := fc.Get("etcd/" + etcdVersion + "/etcdctl")
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/etcdctl", etcdctlRaw, fileMode); err != nil {
		return maskAny(err)
	}

	etcdServiceRaw, err := templates.Asset(etcdServiceTemplate)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(etcdServicePath, etcdServiceRaw, fileMode); err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if startDaemon {
		if err := sc.Start(etcdServiceName); err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func Teardown(fsc *fs.FsClient, sc *systemd.SystemdClient, distributionPath string, stopDaemon bool) error {
	vLogger("\n# call etcd.Teardown()")

	exists, err := sc.Exists(etcdServiceName)
	if err != nil {
		return maskAny(err)
	}

	if exists && stopDaemon {
		if err := sc.Stop(etcdServiceName); err != nil {
			return maskAny(err)
		}
	}

	if err := fsc.Remove(distributionPath + "/etcd2"); err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	return nil
}
