package workspace

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/cosmo-workspace/cosmo/pkg/cmdutil"
)

func AddCommand(rootCmd *cobra.Command, o *cmdutil.CliOptions) {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manipulate WorkspaceTemplate",
		Long: `
Manipulate WorkspaceTemplate.

WorkspaceTemplate is a Template for Workspace
`,
		Aliases: []string{"ws"},
	}

	cmd.AddCommand(GenerateCmd(&cobra.Command{
		Use:     "generate --name TEMPLATE_NAME [< Input via Stdin or pipe]",
		Aliases: []string{"gen"},
		Short:   "Generate Template",
		Long: `Generate Template

For create generated template, just do "kubectl create -f cosmo-template.yaml"

Example:
  * Pipe from kustomize build and apply to your cluster in a single line 
	
      kustomize build ./kubernetes/ | cosmoctl template workspace generate --name TEMPLATE_NAME | kubectl apply -f -

  * Pipe from helm template and generate Workspace Template with cosmo-auth-proxy injection
	
  	  helm template code-server ci/helm-chart \
		| cosmoctl template workspace generate --name TEMPLATE_NAME

  * Input merged config file (kustomize build ... or helm template ... etc.) and save it to file

      cosmoctl template workspace generate --name TEMPLATE_NAME -o cosmo-template.yaml < merged.yaml
`,
	}, o))
	cmd.AddCommand(GetCmd(&cobra.Command{
		Use:   "get [TEMPLATE]",
		Short: "Get WorkspaceTemplates",
		Long: `Get Templates
List or get WorkspaceTemplates in cluster.
`,
	}, o))

	rootCmd.AddCommand(cmd)
}
