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

	// This cli argument is used within newer docker versions to avoid setting our rules.
	iptablesRulesCmd.Flags().BoolVarP(&useDockerIptableRules, "use-docker-rules", "", true, "set our docker iptables rules")

	iptablesRulesCmd.Flags().StringSliceVarP(&privateRegistry, "private-registry", "", strings.Split(defaultPrivateRegistry, ","), "private registry without SSL (for multiple private registry use comma separation)")
}

func iptablesRulesRun(cmd *cobra.Command, args []string) {
	b, err := iptables.RenderRulesFromTemplate(subnet, dockerSubnet, gateway, useDockerIptableRules)
	if err != nil {
		ExitStderr(mask(err))
	}

	Stdoutf(string(b))
}
