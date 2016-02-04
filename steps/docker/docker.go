package docker

import (
	"bufio"
	"os"
	"strings"

	"github.com/giantswarm/yochu/fetchclient"
	fsPkg "github.com/giantswarm/yochu/fs"
	systemdPkg "github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

type dockerOptions struct {
	DockerPath      string
	PrivateRegistry []string
	StorageEngine   string
	UseIPTables     bool
}

const (
	dockerFolder = "/var/lib/docker"

	btrfsFilesystemType   = "btrfs"
	overlayFilesystemType = "overlay"
)

var (
	vLogger = func(f string, v ...interface{}) {}

	serviceName = "docker.service"
	socketName  = "docker-tcp.socket"

	serviceTemplate = "templates/docker.service.tmpl"
	socketTemplate  = "templates/docker-tcp.socket.tmpl"

	servicePath = "/etc/systemd/system/" + serviceName
	socketPath  = "/etc/systemd/system/" + socketName

	fileMode = os.FileMode(0755)
	services = []string{serviceName, socketName}
	paths    = []string{servicePath, socketPath}
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fsPkg.FsClient, sc *systemdPkg.SystemdClient, fc fetchclient.FetchClient, distributionPath, dockerVersion string, privateRegistry []string, useIPTables, restartDaemon bool) error {
	vLogger("\n# call dockerPkg.Setup()")

	dockerRaw, err := fc.Get("docker/" + dockerVersion + "/docker")
	if err != nil {
		return Mask(err)
	}

	if err := fsc.Write(distributionPath+"/docker", dockerRaw, fileMode); err != nil {
		return Mask(err)
	}

	err = createDockerService(fsc, privateRegistry, useIPTables)
	if err != nil {
		return Mask(err)
	}

	dockerTcpSocket, err := templates.Asset(socketTemplate)
	if err != nil {
		return Mask(err)
	}

	// write docker-tcp.socket unit to host
	if err := fsc.Write(socketPath, dockerTcpSocket, fileMode); err != nil {
		return Mask(err)
	}

	// reload unit files, that is, `systemctl daemon-reload`
	if err := sc.Reload(); err != nil {
		return Mask(err)
	}

	if restartDaemon {
		// start docker-tcp.socket unit
		if err := sc.Start(socketName); err != nil {
			return Mask(err)
		}

		// start docker.service unit
		if err := sc.Start(serviceName); err != nil {
			// If there is a dependency error, we just log it. This only happens in case
			// the provisioner is restarted. Then systemd throws an error when starting
			// docker, even though the only dependency (docker-tcp.socket) does not
			// fail.
			if systemdPkg.IsJobDependency(err) {
				vLogger(err.Error())
			} else {
				return Mask(err)
			}
		}
	}

	return nil
}

func Split(s string, d string) (lst []string) {
	lst = strings.Split(s, d)
	return
}

func createDockerService(fsc *fsPkg.FsClient, privateRegistry []string, useIPTables bool) error {
	opts := dockerOptions{
		PrivateRegistry: privateRegistry,
		StorageEngine:   getStorageEngine(dockerFolder),
		UseIPTables:     useIPTables,
	}

	b, err := templates.Render(serviceTemplate, opts)
	if err != nil {
		return Mask(err)
	}

	if err := fsc.Write(servicePath, b.Bytes(), fileMode); err != nil {
		return Mask(err)
	}

	return nil
}

func getStorageEngine(path string) string {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		panic(err)
	}

	bestOptionLen := 0
	bestOptionFilesystem := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		mountPoint := fields[1]
		fs := fields[2]
		if strings.HasPrefix(path, mountPoint) && len(mountPoint) > bestOptionLen {
			bestOptionLen = len(mountPoint)
			bestOptionFilesystem = fs
		}
	}
	if bestOptionLen > 0 {
		if bestOptionFilesystem == "btrfs" {
			return "btrfs"
		}
		return "overlay"
	}
	panic("/proc/mounts doesnt have a rootfs?")
}

func Teardown(fsc *fsPkg.FsClient, sc *systemdPkg.SystemdClient, stopDaemon bool) error {
	vLogger("\n# call dockerPkg.Teardown()")

	for _, s := range services {
		exists, err := sc.Exists(s)
		if err != nil {
			return Mask(err)
		}

		if !exists || !stopDaemon {
			continue
		}

		if err := sc.Stop(s); err != nil {
			return Mask(err)
		}
	}

	for _, p := range paths {
		if err := fsc.Remove(p); err != nil {
			return Mask(err)
		}
	}

	// reload unit files, that is, `systemctl daemon-reload`
	if err := sc.Reload(); err != nil {
		return Mask(err)
	}

	return nil
}
