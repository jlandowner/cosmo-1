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
	workspaceCmd.AddCommand(ResumeCmd(&cobra.Command{
		Use:     "resume",
		Short:   "Resume stopped workspaces",
		Aliases: []string{"start", "run"},
	}, o))
	workspaceCmd.AddCommand(SuspendCmd(&cobra.Command{
		Use:     "suspend",
		Short:   "Suspend workspaces",
		Aliases: []string{"stop"},
	}, o))
	workspaceCmd.AddCommand(GetNetworkCmd(&cobra.Command{
		Use:     "get-network WORKSPACE_NAME",
		Short:   "Get workspace network",
		Aliases: []string{"get-ns", "ns", "network", "get-networks"},
	}, o))
	workspaceCmd.AddCommand(UpsertNetworkCmd(&cobra.Command{
		Use:     "upsert-network WORKSPACE_NAME --port 8080",
		Short:   "Upsert workspace network",
		Aliases: []string{"upsert-ns", "add-ns", "add-network"},
	}, o))
	workspaceCmd.AddCommand(RemoveNetworkCmd(&cobra.Command{
		Use:     "remove-network WORKSPACE_NAME --port 8080",
		Short:   "Remove workspace network",
		Aliases: []string{"remove-ns", "delete-ns", "delete-network"},
	}, o))

	cmd.AddCommand(workspaceCmd)
}
