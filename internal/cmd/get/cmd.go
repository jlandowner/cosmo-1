package get

import (
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/internal/cmd/template/useraddon"
	wstmpl "github.com/cosmo-workspace/cosmo/internal/cmd/template/workspace"
	"github.com/cosmo-workspace/cosmo/internal/cmd/user"
	"github.com/cosmo-workspace/cosmo/internal/cmd/workspace"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
)

func AddCommand(rootCmd *cobra.Command, o *cmdutil.CliOptions) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get cosmo resources",
	}

	cmd.AddCommand(workspace.GetCmd(&cobra.Command{
		Use:     "workspace WORKSPACE_NAME",
		Aliases: []string{"ws"},
		Short:   "Get workspace",
	}, o))
	cmd.AddCommand(user.GetCmd(&cobra.Command{
		Use:   "user USER_NAME",
		Short: "Get user",
	}, o))

	tmplCmd := &cobra.Command{
		Use:   "template",
		Short: "Get Templates",
	}
	tmplCmd.AddCommand(wstmpl.GetCmd(&cobra.Command{
		Use:     "workspace WORKSPACE_NAME",
		Aliases: []string{"wstmpl"},
		Short:   "Get WorkspaceTemplates",
	}, o))
	tmplCmd.AddCommand(useraddon.GetCmd(&cobra.Command{
		Use:     "useraddon ADDON_NAME",
		Aliases: []string{"ua"},
		Short:   "Get UserAddons",
	}, o))
	cmd.AddCommand(tmplCmd)

	rootCmd.AddCommand(cmd)
}
