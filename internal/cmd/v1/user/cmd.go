package user

import (
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Manipulate User resource",
		Long: `
Manipulate COSMO User resource.

"User" is a cluster-scoped Kubernetes CRD which represents a developer or user who use Workspace.

Once you create User, Kubernetes Namespace is created and bound to the User.
`,
	}

	// userCmd.AddCommand(resetPasswordCmd(&cobra.Command{
	// 	Use:   "reset-password USER_NAME",
	// 	Short: "Reset user password",
	// }, o))
	userCmd.AddCommand(CreateCmd(&cobra.Command{
		Use:   "create USER_NAME",
		Short: "Create user",
	}, o))
	userCmd.AddCommand(GetCmd(&cobra.Command{
		Use:   "get",
		Short: "Get users",
		Long: `
Get Users.
`,
	}, o))
	// userCmd.AddCommand(DeleteCmd(&cobra.Command{
	// 	Use:     "delete USER_NAME",
	// 	Aliases: []string{"del"},
	// 	Short:   "Delete user",
	// }, o))

	cmd.AddCommand(userCmd)
}
