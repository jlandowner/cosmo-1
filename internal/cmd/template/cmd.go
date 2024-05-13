package template

import (
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	templateCmd := &cobra.Command{
		Use:     "template",
		Short:   "Manipulate Template resource",
		Aliases: []string{"tmpl"},
		Long: `
Manipulate COSMO Template resource.

"Template" is a set of Kubernetes resources
`,
	}
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate Template",
		Long: `Generate Template

For create generated template, just do "kubectl create -f cosmo-template.yaml"
`,
		Example: `
  * Pipe from kustomize build and apply to your cluster in a single line 
	
      kustomize build ./kubernetes/ | cosmoctl tmpl gen --name TEMPLATE_NAME | kubectl apply -f -

  * Input merged config file (kustomize build ... or helm template ... etc.) and save it to file

      cosmoctl tmpl gen --name TEMPLATE_NAME -o cosmo-template.yaml < merged.yaml
`,
		Aliases: []string{"gen"},
	}
	generateCmd.AddCommand(generateWorkspaceCmd(&cobra.Command{
		Use:     "workspace [< Input via Stdin or pipe]",
		Short:   "Generate WorkspaceTemplate",
		Aliases: []string{"ws"},
	}, o))
	generateCmd.AddCommand(generateUserAddonCmd(&cobra.Command{
		Use:     "useraddon [< Input via Stdin or pipe]",
		Short:   "Generate UserAddon",
		Aliases: []string{"addon"},
	}, o))

	templateCmd.AddCommand(validateCmd(&cobra.Command{
		Use:     "validate --file FILE",
		Aliases: []string{"valid", "check"},
		Short:   "Validate Template by dry-run",
		Example: `
  * Dry-run on server-side
	
      cosmoctl template validate -f cosmo-template.yaml

  * Dry-run on client-side using kubectl
	
      cosmoctl template validate -f cosmo-template.yaml --client

  * Input from stdin not file.

      cat cosmo-template.yaml | cosmoctl template validate -f -
`,
	}, o))

	templateCmd.AddCommand(generateCmd)
	cmd.AddCommand(templateCmd)
}
