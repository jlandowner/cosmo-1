package get

import (
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/user"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/spf13/cobra"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get cosmo resources",
		Long: `
Get cosmo resources
`,
	}

	getCmd.AddCommand(user.GetCmd(&cobra.Command{
		Use:     "user USER_NAME",
		Short:   "Get users",
		Aliases: []string{"users"},
	}, o))
	cmd.AddCommand(getCmd)
}
