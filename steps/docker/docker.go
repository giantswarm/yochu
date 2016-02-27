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
	UseTypeNotify   bool
	DockerExecArgs  []string
}

const (
	dockerFolder = "/var/lib/docker"

	dockerDaemonV1_10Arg = "daemon"
	dockerDaemonArg      = "-d"
	dockerIccEnabledArg  = "--icc=true"

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
		return maskAny(err)
	}

	if err := fsc.Write(distributionPath+"/docker", dockerRaw, fileMode); err != nil {
		return maskAny(err)
	}

	err = createDockerService(fsc, dockerVersion, privateRegistry, useIPTables)
	if err != nil {
		return maskAny(err)
	}

	dockerTcpSocket, err := templates.Asset(socketTemplate)
	if err != nil {
		return maskAny(err)
	}

	// write docker-tcp.socket unit to host
	if err := fsc.Write(socketPath, dockerTcpSocket, fileMode); err != nil {
		return maskAny(err)
	}

	// reload unit files, that is, `systemctl daemon-reload`
	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	if restartDaemon {
		// start docker-tcp.socket unit
		if err := sc.Start(socketName); err != nil {
			return maskAny(err)
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
				return maskAny(err)
			}
		}
	}

	return nil
}

func Split(s string, d string) (lst []string) {
	lst = strings.Split(s, d)
	return
}

func createDockerService(fsc *fsPkg.FsClient, dockerVersion string, privateRegistry []string, useIPTables bool) error {
	opts := dockerOptions{
		PrivateRegistry: privateRegistry,
		StorageEngine:   getStorageEngine(dockerFolder),
		UseIPTables:     useIPTables,
		DockerExecArgs:  make([]string, 0),
	}

	options := addVersionSpecificArguments(&opts, dockerVersion)
	opts = *options

	b, err := templates.Render(serviceTemplate, opts)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(servicePath, b.Bytes(), fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func addVersionSpecificArguments(opts *dockerOptions, dockerVersion string) *dockerOptions {
	opts.DockerExecArgs = append(opts.DockerExecArgs, dockerDaemonArg)

	if normalizeVersion(dockerVersion) >= normalizeVersion("1.9") {
		opts.UseTypeNotify = true
		if normalizeVersion(dockerVersion) >= normalizeVersion("1.10") {
			opts.DockerExecArgs = append(opts.DockerExecArgs, dockerDaemonV1_10Arg)
			opts.DockerExecArgs = append(opts.DockerExecArgs[:0], opts.DockerExecArgs[1:]...)
		}
		opts.DockerExecArgs = append(opts.DockerExecArgs, dockerIccEnabledArg)
	}
	return opts
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
			return maskAny(err)
		}

		if !exists || !stopDaemon {
			continue
		}

		if err := sc.Stop(s); err != nil {
			return maskAny(err)
		}
	}

	for _, p := range paths {
		if err := fsc.Remove(p); err != nil {
			return maskAny(err)
		}
	}

	// reload unit files, that is, `systemctl daemon-reload`
	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	return nil
}

func normalizeVersion(version string) string {
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			vLogger("unable to normalize this version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}
