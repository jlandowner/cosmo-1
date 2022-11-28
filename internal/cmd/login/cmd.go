package login

import (
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func AddCommand(cmd *cobra.Command, o *cmdutil.CliOptions) {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to cosmo server",
		Long: `
Login to cosmo server
`,
	}

	//TODO

	cmd.AddCommand(loginCmd)
}
