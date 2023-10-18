package get

import (
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/user"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/workspace"
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
		Use:     "user [USER_NAME]",
		Short:   "Get users. Alias of 'cosmoctl user get'",
		Aliases: []string{"users"},
	}, o))
	getCmd.AddCommand(workspace.GetCmd(&cobra.Command{
		Use:     "workspace [WORKSPACE_NAME]",
		Short:   "Get workspaces. Alias of 'cosmoctl workspace get'",
		Aliases: []string{"workspaces", "ws"},
	}, o))
	cmd.AddCommand(getCmd)
}
