package cli

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/yochu/steps/ip6tables"
)

var (
	ip6tablesRulesCmd = &cobra.Command{
		Use:   "ip6tables-rules",
		Short: "show ip6tables-rules",
		Long:  "show ip6tables-rules",
		Run:   ip6tablesRulesRun,
	}
)

func ip6tablesRulesRun(cmd *cobra.Command, args []string) {
	b, err := ip6tables.RenderRulesFromTemplate()
	if err != nil {
		ExitStderr(mask(err))
	}

	Stdoutf("%s", string(b))
}
