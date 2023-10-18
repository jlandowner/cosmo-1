package workspace

import (
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	workspaceCmd := &cobra.Command{
		Use:     "workspace",
		Short:   "Manipulate Workspace resource",
		Aliases: []string{"ws"},
		Long: `
Manipulate COSMO Workspace resource.

"Workspace" is a namespaced Kubernetes CRD which represents a instance of workspace.
`,
	}

	workspaceCmd.AddCommand(CreateCmd(&cobra.Command{
		Use:   "create WORKSPACE_NAME",
		Short: "Create workspace",
	}, o))
	workspaceCmd.AddCommand(GetCmd(&cobra.Command{
		Use:     "get",
		Short:   "Get workspaces",
		Aliases: []string{"list"},
		Long: `
Get Workspaces.
`,
	}, o))
	workspaceCmd.AddCommand(GetTemplatesCmd(&cobra.Command{
		Use:     "get-templates",
		Short:   "Get workspace templates",
		Aliases: []string{"get-tmpls", "get-tmpl", "get-templates", "get-template"},
		Long: `
List workspaceaddon templates in cluster.
`,
	}, o))
	workspaceCmd.AddCommand(DeleteCmd(&cobra.Command{
		Use:     "delete WORKSPACE_NAME",
		Aliases: []string{"rm"},
		Short:   "Delete workspace",
	}, o))

	cmd.AddCommand(workspaceCmd)
}
