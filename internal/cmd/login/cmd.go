package login

import (
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func AddCommand(rootCmd *cobra.Command, o *cmdutil.CliOptions) {
	cmd := LoginCmd(&cobra.Command{
		Use:   "login",
		Short: "Login to cosmo server",
		Long: `
Login to cosmo server
`,
	}, o)
	rootCmd.AddCommand(cmd)
}
