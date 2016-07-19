package cli

import (
	"strings"

	"github.com/spf13/cobra"

	fsPkg "github.com/giantswarm/yochu/fs"
	httpClientPkg "github.com/giantswarm/yochu/httpclient"
	s3Pkg "github.com/giantswarm/yochu/s3"
	distribution "github.com/giantswarm/yochu/steps/distribution"
	dockerPkg "github.com/giantswarm/yochu/steps/docker"
	etcdPkg "github.com/giantswarm/yochu/steps/etcd"
	fleetPkg "github.com/giantswarm/yochu/steps/fleet"
	ip6tablesPkg "github.com/giantswarm/yochu/steps/ip6tables"
	iptablesPkg "github.com/giantswarm/yochu/steps/iptables"
	k8sPkg "github.com/giantswarm/yochu/steps/k8s"
	overlay "github.com/giantswarm/yochu/steps/overlay"
	rktPkg "github.com/giantswarm/yochu/steps/rkt"
	systemdPkg "github.com/giantswarm/yochu/systemd"
)

var (
	globalFlags = struct {
		debug   bool
		verbose bool

		version string
		steps   string
	}{}

	yochuCmd = &cobra.Command{
		Use:   "yochu",
		Short: "provision swarm cluster",
		Long:  "provision swarm cluster",
		Run:   yochuRun,
	}
)

func init() {
	fsPkg.Configure(Verbosef)
	systemdPkg.Configure(Verbosef)

	httpClientPkg.Configure(Verbosef)
	s3Pkg.Configure(Verbosef)

	iptablesPkg.Configure(Verbosef)
	ip6tablesPkg.Configure(Verbosef)
	dockerPkg.Configure(Verbosef)
	fleetPkg.Configure(Verbosef)
	etcdPkg.Configure(Verbosef)
	distribution.Configure(Verbosef)
	overlay.Configure(Verbosef)
	rktPkg.Configure(Verbosef)
	k8sPkg.Configure(Verbosef)

	yochuCmd.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "d", false, "print debug output")
	yochuCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "print verbose output")

	yochuCmd.PersistentFlags().StringVar(&globalFlags.steps, "steps", "all", "comma separated steps to execute")
}

// NewYochuCmd returns a Cobra command to run Yochu.
func NewYochuCmd(version string) *cobra.Command {
	globalFlags.version = version

	yochuCmd.AddCommand(versionCmd)
	yochuCmd.AddCommand(setupCmd)
	yochuCmd.AddCommand(teardownCmd)
	yochuCmd.AddCommand(iptablesRulesCmd)
	yochuCmd.AddCommand(ip6tablesRulesCmd)
	yochuCmd.Execute()

	return yochuCmd
}

func yochuRun(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}

// Check if the given step is listed in steps. steps is a comma separated list
// of steps. If there are no steps given by the cli, steps is all.
func execute(steps, step string) bool {
	if steps == "all" {
		return true
	}

	for _, s := range strings.Split(steps, ",") {
		if s == step {
			return true
		}
	}

	return false
}
