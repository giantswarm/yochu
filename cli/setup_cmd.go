package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/giantswarm/yochu/fetchclient"
	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/httpclient"
	"github.com/giantswarm/yochu/s3"
	"github.com/giantswarm/yochu/steps/distribution"
	"github.com/giantswarm/yochu/steps/docker"
	"github.com/giantswarm/yochu/steps/etcd"
	"github.com/giantswarm/yochu/steps/fleet"
	"github.com/giantswarm/yochu/steps/ip6tables"
	"github.com/giantswarm/yochu/steps/iptables"
	"github.com/giantswarm/yochu/steps/overlay"
	"github.com/giantswarm/yochu/systemd"
)

const (
	defaultSubnet           = "172.31.0.0/16"
	defaultDockerSubnet     = "172.17.0.0/16"
	defaultGateway          = "172.31.0.2/32"
	defaultDistributionPath = "/opt/giantswarm/bin"
	defaultOverlayWorkdir   = "/opt/giantswarm/overlay_workdir"
	defaultFleetVersion     = "v0.11.3-gs-2"
	defaultEtcdVersion      = "v2.1.0-gs-1"
	defaultDockerVersion    = "1.6.2"
	defaultS3Bucket         = "downloads.giantswarm.io"
	defaultHTTPEndpoint     = "https://downloads.giantswarm.io"
	defaultPrivateRegistry  = ""

	overlayMountPoint = "/usr/bin"
)

var (
	subnet       string
	dockerSubnet string
	gateway      string

	privateRegistry []string

	distributionPath string
	overlayWorkdir   string
	fleetVersion     string
	etcdVersion      string
	dockerVersion    string
	awsAccessKey     string
	awsSecretKey     string
	s3Endpoint       string
	s3Bucket         string
	httpEndpoint     string

	startDaemons bool
	stopDaemons  bool

	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "setup swarm cluster",
		Long:  "setup swarm cluster",
		Run:   setupRun,
	}
)

func init() {
	setupCmd.Flags().StringVarP(&subnet, "subnet", "", defaultSubnet, "subnet for the iptables rules")
	setupCmd.Flags().StringVarP(&dockerSubnet, "docker-subnet", "", defaultDockerSubnet, "docker subnet for the iptables rules")
	setupCmd.Flags().StringVarP(&gateway, "gateway", "", defaultGateway, "gateway for the host")

	setupCmd.Flags().StringSliceVarP(&privateRegistry, "private-registry", "", strings.Split(defaultPrivateRegistry, ","), "private registry without ssl (for multiple private registry use comma separation)")

	setupCmd.Flags().StringVarP(&distributionPath, "distribution-path", "", defaultDistributionPath, "path to use for custom binary distribution")
	setupCmd.Flags().StringVarP(&overlayWorkdir, "overlay-workdir", "", defaultOverlayWorkdir, "workdir to use for distribution overlay mount")
	setupCmd.Flags().StringVarP(&fleetVersion, "fleet-version", "", defaultFleetVersion, "version to use when provisioning fleetd/fleetctl binaries")
	setupCmd.Flags().StringVarP(&etcdVersion, "etcd-version", "", defaultEtcdVersion, "version to use when provisioning etcd binaries")
	setupCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", defaultDockerVersion, "version to use when provisioning docker binaries")
	setupCmd.Flags().StringVarP(&awsAccessKey, "aws-access-key", "", "", "AWS access key to use (can also be set via the AWS_ACCESS_KEY_ID env variable)")
	setupCmd.Flags().StringVarP(&awsSecretKey, "aws-secret-key", "", "", "AWS secret key to use (can also be set via the AWS_SECRET_ACCESS_KEY env variable)")
	setupCmd.Flags().StringVarP(&s3Endpoint, "s3-endpoint", "", "", "S3 endpoint to use (can also be set via the S3_ENDPOINT env variable)")
	setupCmd.Flags().StringVarP(&s3Bucket, "s3-bucket", "", defaultS3Bucket, "S3 bucket to use")
	setupCmd.Flags().StringVarP(&httpEndpoint, "http-endpoint", "", defaultHTTPEndpoint, "HTTP endpoint to use")
	setupCmd.Flags().BoolVarP(&startDaemons, "start-daemons", "", true, "start daemons after deploying")
	setupCmd.Flags().BoolVarP(&stopDaemons, "stop-daemons", "", true, "stop daemons before deploying")
}

func setupRun(cmd *cobra.Command, args []string) {
	fsClient, err := fs.NewFsClient()
	if err != nil {
		ExitStderr(Mask(err))
	}

	systemdClient, err := systemd.NewSystemdClient()
	if err != nil {
		ExitStderr(Mask(err))
	}

	var fetchClient fetchclient.FetchClient
	if awsAccessKey == "" || awsSecretKey == "" {
		fetchClient, err = httpclient.NewHTTPClient(httpEndpoint)
		if err != nil {
			ExitStderr(Mask(err))
		}
	} else {
		fetchClient, err = s3.NewS3Client(awsAccessKey, awsSecretKey, s3Endpoint, s3Bucket)
		if err != nil {
			ExitStderr(Mask(err))
		}
	}

	// distribution service
	if execute(globalFlags.steps, "distribution") {
		if err := distribution.Teardown(fsClient, distributionPath); err != nil {
			ExitStderr(Mask(err))
		}

		if err := distribution.Setup(fsClient, distributionPath); err != nil {
			ExitStderr(Mask(err))
		}
	}

	// overlay mount service
	if execute(globalFlags.steps, "overlay") {
		if err := overlay.Teardown(fsClient, systemdClient, distributionPath, overlayWorkdir, overlayMountPoint); err != nil {
			ExitStderr(Mask(err))
		}

		if err := overlay.Setup(fsClient, systemdClient, distributionPath, overlayWorkdir, overlayMountPoint); err != nil {
			ExitStderr(Mask(err))
		}
	}

	// iptables
	useIPTables := execute(globalFlags.steps, "iptables")
	if useIPTables {
		if err := iptables.Teardown(fsClient, systemdClient); err != nil {
			ExitStderr(Mask(err))
		}

		if err := iptables.Setup(fsClient, systemdClient, subnet, dockerSubnet, gateway); err != nil {
			ExitStderr(Mask(err))
		}
	}

	// ip6tables
	useIP6Tables := execute(globalFlags.steps, "ip6tables")
	if useIP6Tables {
		if err := ip6tables.Teardown(fsClient, systemdClient); err != nil {
			ExitStderr(Mask(err))
		}

		if err := ip6tables.Setup(fsClient, systemdClient); err != nil {
			ExitStderr(Mask(err))
		}
	}

	// docker service
	if execute(globalFlags.steps, "docker") {
		if err := docker.Teardown(fsClient, systemdClient, stopDaemons); err != nil {
			ExitStderr(Mask(err))
		}

		// !useIPTables is used, because when the iptables step is enabled, we want
		// --iptables=false for the docker daemon.
		if err := docker.Setup(fsClient, systemdClient, fetchClient, overlayMountPoint, dockerVersion, privateRegistry, !useIPTables, startDaemons); err != nil {
			ExitStderr(Mask(err))
		}
	}

	// fleet service
	if execute(globalFlags.steps, "fleet") {
		if err := fleet.Teardown(fsClient, systemdClient, overlayMountPoint, stopDaemons); err != nil {
			ExitStderr(Mask(err))
		}

		if err := fleet.Setup(fsClient, systemdClient, fetchClient, overlayMountPoint, fleetVersion, startDaemons); err != nil {
			ExitStderr(Mask(err))
		}
	}

	// etcd service
	if execute(globalFlags.steps, "etcd") {
		if err := etcd.Teardown(fsClient, systemdClient, overlayMountPoint, stopDaemons); err != nil {
			ExitStderr(Mask(err))
		}

		if err := etcd.Setup(fsClient, systemdClient, fetchClient, overlayMountPoint, etcdVersion, startDaemons); err != nil {
			ExitStderr(Mask(err))
		}
	}
}
