package templates

import (
	"strings"
	"testing"
)

type dockerOptions struct {
	DockerPath      string
	PrivateRegistry []string
	StorageEngine   string
	UseIPTables     bool
	UseOverlay      bool
	UseTypeNotify   bool
	DockerExecArgs  []string
}

const (
	dockerFolder = "/var/lib/docker"

	dockerDaemonV1_10Arg = "daemon"
	dockerDaemonArg      = "-d"
	dockerIccEnabledArg  = "--icc=true"
)

var dockerCgroupDriverArgs = []string{"--exec-opt", "native.cgroupdriver=cgroupfs"}

func TestRender(t *testing.T) {
	// We picked the usr-bin template since it is rather simple and contains only a single value: "ExecStart"
	asset := "templates/usr-bin.mount.tmpl"
	data := map[string]string{
		"MountPoint": "THIS_IS_A_TEST",
	}
	content, err := Render(asset, data)
	if err != nil {
		t.Fatalf("rendering of %s failed: %v", asset, err.Error())
	}

	if !strings.Contains(content.String(), data["ExecStart"]) {
		t.Fatalf("expected rendered template to contain passed value.")
	}
}

func TestDockerServiceRender(t *testing.T) {
	// We picked the usr-bin template since it is rather simple and contains only a single value: "ExecStart"
	asset := "templates/docker.service.tmpl"
	opts := &dockerOptions{
		PrivateRegistry: []string{"private.registry.com"},
		StorageEngine:   "btrfs",
		UseIPTables:     true,
		UseOverlay:      true,
		UseTypeNotify:   true,
		DockerExecArgs:  []string{dockerDaemonV1_10Arg, dockerIccEnabledArg},
	}
	opts.DockerExecArgs = append(opts.DockerExecArgs, dockerCgroupDriverArgs...)

	content, err := Render(asset, opts)
	if err != nil {
		t.Fatalf("rendering of %s failed: %v", asset, err.Error())
	}

	if !strings.Contains(content.String(), "Type=Notify") {
		t.Fatalf("expected rendered template to contain passed value. %v", content.String())
	}
}
