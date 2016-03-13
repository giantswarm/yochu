package cli

import (
	"github.com/spf13/cobra"

	fsPkg "github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/steps/distribution"
	dockerPkg "github.com/giantswarm/yochu/steps/docker"
	etcdPkg "github.com/giantswarm/yochu/steps/etcd"
	fleetPkg "github.com/giantswarm/yochu/steps/fleet"
	ip6tablesPkg "github.com/giantswarm/yochu/steps/ip6tables"
	iptablesPkg "github.com/giantswarm/yochu/steps/iptables"
	k8sPkg "github.com/giantswarm/yochu/steps/k8s"
	"github.com/giantswarm/yochu/steps/overlay"
	rktPkg "github.com/giantswarm/yochu/steps/rkt"
	systemdPkg "github.com/giantswarm/yochu/systemd"
)

var (
	teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown swarm cluster",
		Long:  "teardown swarm cluster",
		Run:   teardownRun,
	}
)

func init() {
	teardownCmd.Flags().StringVarP(&distributionPath, "distribution-path", "", defaultDistributionPath, "Path to use for custom binary distribution")
	teardownCmd.Flags().BoolVarP(&stopDaemons, "stop-daemons", "", true, "Stop daemons before deploying")
}

func teardownRun(cmd *cobra.Command, args []string) {
	fs, err := fsPkg.NewFsClient()
	if err != nil {
		ExitStderr(err)
	}

	systemd, err := systemdPkg.NewSystemdClient()
	if err != nil {
		ExitStderr(err)
	}

	// iptables
	if execute(globalFlags.steps, "iptables") {
		if err := iptablesPkg.Teardown(fs, systemd); err != nil {
			ExitStderr(err)
		}
	}

	// ip6tables
	if execute(globalFlags.steps, "ip6tables") {
		if err := ip6tablesPkg.Teardown(fs, systemd); err != nil {
			ExitStderr(err)
		}
	}

	// k8s binaries
	if execute(globalFlags.steps, "k8s") {
		if err := k8sPkg.Teardown(fs, overlayMountPoint); err != nil {
			ExitStderr(err)
		}
	}

	// rkt binaries
	if execute(globalFlags.steps, "rkt") {
		if err := rktPkg.Teardown(fs, systemd, overlayMountPoint, stopDaemons); err != nil {
			ExitStderr(err)
		}
	}

	// docker service
	if execute(globalFlags.steps, "docker") {
		if err := dockerPkg.Teardown(fs, systemd, stopDaemons); err != nil {
			ExitStderr(err)
		}
	}

	// fleet service
	if execute(globalFlags.steps, "fleet") {
		if err := fleetPkg.Teardown(fs, systemd, overlayMountPoint, stopDaemons); err != nil {
			ExitStderr(err)
		}
	}

	// etcd service
	if execute(globalFlags.steps, "etcd") {
		if err := etcdPkg.Teardown(fs, systemd, overlayMountPoint, stopDaemons); err != nil {
			ExitStderr(err)
		}
	}

	// overlay service
	if execute(globalFlags.steps, "overlay") {
		if err := overlay.Teardown(fs, systemd, distributionPath, overlayWorkdir, overlayMountPoint); err != nil {
			ExitStderr(err)
		}
	}

	// distribution service
	if execute(globalFlags.steps, "distribution") {
		if err := distribution.Teardown(fs, distributionPath); err != nil {
			ExitStderr(err)
		}
	}
}
