package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/internal/cmd/create"
	del "github.com/cosmo-workspace/cosmo/internal/cmd/delete"
	"github.com/cosmo-workspace/cosmo/internal/cmd/get"
	"github.com/cosmo-workspace/cosmo/internal/cmd/login"
	"github.com/cosmo-workspace/cosmo/internal/cmd/netrule"
	"github.com/cosmo-workspace/cosmo/internal/cmd/run"
	"github.com/cosmo-workspace/cosmo/internal/cmd/stop"
	"github.com/cosmo-workspace/cosmo/internal/cmd/template"
	"github.com/cosmo-workspace/cosmo/internal/cmd/user"
	"github.com/cosmo-workspace/cosmo/internal/cmd/workspace"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
)

func NewRootCmd(o *cmdutil.CliOptions) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "cosmoctl",
		Short: "Command line tool to manipulate cosmo resources",
		Long: `
Command line tool to manipulate cosmo resources
Complete documentation is available at http://github.com/cosmo-workspace/cosmo

MIT 2022 cosmo-workspace/cosmo
`,
	}

	rootCmd.SetIn(o.In)
	rootCmd.SetOut(o.Out)
	rootCmd.SetErr(o.ErrOut)
	o.AddFlags(rootCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(o.Out, "cosmoctl - cosmo v0.7.0 cosmo-workspace 2022")
		},
	}

	rootCmd.AddCommand(versionCmd)
	template.AddCommand(rootCmd, o)
	user.AddCommand(rootCmd, o)
	workspace.AddCommand(rootCmd, o)
	netrule.AddCommand(rootCmd, o)

	create.AddCommand(rootCmd, o)
	get.AddCommand(rootCmd, o)
	del.AddCommand(rootCmd, o)
	run.AddCommand(rootCmd, o)
	stop.AddCommand(rootCmd, o)

	login.AddCommand(rootCmd, o)

	return rootCmd
}

func Execute() {
	o := cmdutil.NewCliOptions()
	o.In = os.Stdin
	o.Out = os.Stdout
	o.ErrOut = os.Stderr
	rootCmd := NewRootCmd(o)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(o.Out, err)
		os.Exit(1)
	}

}
