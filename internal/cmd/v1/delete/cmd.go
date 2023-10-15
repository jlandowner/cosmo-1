package delete

import (
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/user"
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
		Use:   "user USER_NAME",
		Short: "Delete user",
	}, o))
	cmd.AddCommand(deleteCmd)
}
