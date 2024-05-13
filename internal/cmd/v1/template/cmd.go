package template

import (
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	templateCmd := &cobra.Command{
		Use:     "template",
		Short:   "Template utitlity commands",
		Aliases: []string{"tmpl"},
	}
	templateCmd.AddCommand(generateCmd(&cobra.Command{
		Use:     "generate",
		Short:   "Generate Template",
		Aliases: []string{"gen"},
	}, o))
	cmd.AddCommand(templateCmd)
}
