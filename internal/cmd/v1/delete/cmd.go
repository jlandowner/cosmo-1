package delete

import (
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/user"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/workspace"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/spf13/cobra"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete cosmo resources",
		Aliases: []string{"rm", "remove"},
		Long: `
Delete cosmo resources
`,
	}

	deleteCmd.AddCommand(user.DeleteCmd(&cobra.Command{
		Use:     "user USER_NAME",
		Short:   "Delete user. Alias of 'cosmoctl user delete'",
		Aliases: []string{"us", "users"},
	}, o))
	deleteCmd.AddCommand(workspace.DeleteCmd(&cobra.Command{
		Use:     "workspace WORKSPACE_NAME",
		Short:   "Delete workspace. Alias of 'cosmoctl workspace delete'",
		Aliases: []string{"ws", "workspaces"},
	}, o))
	deleteCmd.AddCommand(workspace.RemoveNetworkCmd(&cobra.Command{
		Use:     "network WORKSPACE_NAME",
		Short:   "Remove workspace network. Alias of 'cosmoctl workspace remove-network'",
		Aliases: []string{"workspace-network", "workspace-networks", "ws-net", "wsnet"},
	}, o))
	cmd.AddCommand(deleteCmd)
}
