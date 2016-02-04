package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/giantswarm/yochu/steps/iptables"
)

var (
	iptablesRulesCmd = &cobra.Command{
		Use:   "iptables-rules",
		Short: "show iptables-rules",
		Long:  "show iptables-rules",
		Run:   iptablesRulesRun,
	}
)

func init() {
	iptablesRulesCmd.Flags().StringVarP(&subnet, "subnet", "", defaultSubnet, "subnet for the iptables rules")
	iptablesRulesCmd.Flags().StringVarP(&dockerSubnet, "docker-subnet", "", defaultDockerSubnet, "docker subnet for the iptables rules")
	iptablesRulesCmd.Flags().StringVarP(&gateway, "gateway", "", defaultGateway, "gateway for the host")
	iptablesRulesCmd.Flags().StringSliceVarP(&privateRegistry, "private-registry", "", strings.Split(defaultPrivateRegistry, ","), "private registry without SSL (for multiple private registry use comma separation)")
}

func iptablesRulesRun(cmd *cobra.Command, args []string) {
	b, err := iptables.RenderRulesFromTemplate(subnet, dockerSubnet, gateway)
	if err != nil {
		ExitStderr(Mask(err))
	}

	Stdoutf(string(b))
}
