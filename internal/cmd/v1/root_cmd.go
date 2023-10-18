/*
Copyright © 2023 cosmo-workspace
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/create"
	del "github.com/cosmo-workspace/cosmo/internal/cmd/v1/delete"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/get"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/login"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/user"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/version"
	"github.com/cosmo-workspace/cosmo/internal/cmd/v1/workspace"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
)

func NewRootCmd(o *cli.RootOptions) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "cosmoctl",
		Short: "Command line tool for cosmo API",
		Long: `
Command line tool for cosmo API
Complete documentation is available at http://github.com/cosmo-workspace/cosmo

MIT 2023 cosmo-workspace/cosmo
`,
	}
	o.AddFlags(rootCmd)

	version.AddCommand(rootCmd, o)
	login.AddCommand(rootCmd, o)
	user.AddCommand(rootCmd, o)
	workspace.AddCommand(rootCmd, o)
	// template.AddCommand(rootCmd, o)
	// netrule.AddCommand(rootCmd, o)

	create.AddCommand(rootCmd, o)
	get.AddCommand(rootCmd, o)
	del.AddCommand(rootCmd, o)
	// run.AddCommand(rootCmd, o)
	// stop.AddCommand(rootCmd, o)

	return rootCmd
}

func Execute(v cli.VersionInfo) {
	o := cli.NewRootOptions()
	o.Versions = v
	o.In = os.Stdin
	o.Out = os.Stdout
	o.ErrOut = os.Stderr
	rootCmd := NewRootCmd(o)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(o.Out, err)
		os.Exit(1)
	}

}
