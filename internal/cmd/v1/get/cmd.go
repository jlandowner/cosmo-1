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
		Use:     "template [TEMPLATE_NAME]",
		Short:   "Get workspace templates",
		Aliases: []string{"templates", "tmpl", "tmpls", "ws-tmpl", "ws-tmpls", "wstmpl", "wstmpls"},
	}, o))
	getCmd.AddCommand(user.GetAddonsCmd(&cobra.Command{
		Use:     "useraddon [ADDON_NAME]",
		Short:   "Get user addons. Alias of 'cosmoctl user get-addons'",
		Aliases: []string{"useraddon", "addon", "addons", "user-addon", "user-addons"},
	}, o))
	getCmd.AddCommand(workspace.GetNetworkCmd(&cobra.Command{
		Use:     "network WORKSPACE_NAME",
		Short:   "Get workspace networks",
		Aliases: []string{"workspace-network", "workspace-networks", "ws-net", "wsnet"},
	}, o))
	getCmd.AddCommand(workspace.GetNetworkCmd(&cobra.Command{
		Use:     "events [USER_NAME]",
		Short:   "Get events for user",
		Aliases: []string{"event"},
	}, o))
	cmd.AddCommand(getCmd)
}
