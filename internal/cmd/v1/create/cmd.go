package create

import (
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/user"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/workspace"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/spf13/cobra"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create cosmo resources",
		Long: `
Create cosmo resources
`,
	}

	createCmd.AddCommand(user.CreateCmd(&cobra.Command{
		Use:     "user USER_NAME",
		Short:   "Create user. Alias of 'cosmoctl user create'",
		Aliases: []string{"us"},
	}, o))

	createCmd.AddCommand(workspace.CreateCmd(&cobra.Command{
		Use:     "workspace WORKSPACE_NAME --template TEMPLATE_NAME",
		Short:   "Create workspace. Alias of 'cosmoctl workspace create'",
		Aliases: []string{"ws"},
	}, o))

	createCmd.AddCommand(workspace.UpsertNetworkCmd(&cobra.Command{
		Use:     "workspace-network WORKSPACE_NAME",
		Short:   "Upsert workspace network. Alias of 'cosmoctl workspace upsert-network'",
		Aliases: []string{"workspace-net", "ws-net"},
	}, o))

	cmd.AddCommand(createCmd)
}
