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
	getCmd.AddCommand(workspace.GetTemplatesCmd(&cobra.Command{
		Use:     "workspace-template [TEMPLATE_NAME]",
		Short:   "Get workspace templates. Alias of 'cosmoctl workspace get-template'",
		Aliases: []string{"workspace-templates", "workspace-template", "ws-templates", "ws-template", "ws-tmpl", "ws-tmpls"},
	}, o))
	getCmd.AddCommand(user.GetAddonsCmd(&cobra.Command{
		Use:     "useraddons [ADDON_NAME]",
		Short:   "Get user addons. Alias of 'cosmoctl user get-addons'",
		Aliases: []string{"useraddon", "addon", "addons", "user-addon", "user-addons"},
	}, o))
	getCmd.AddCommand(workspace.GetNetworkCmd(&cobra.Command{
		Use:     "workspace-network WORKSPACE_NAME",
		Short:   "Get workspace networks. Alias of 'cosmoctl workspace get-network'",
		Aliases: []string{"workspace-networks", "workspace-net", "ws-net"},
	}, o))
	cmd.AddCommand(getCmd)
}
