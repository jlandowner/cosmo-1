package useraddon

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/cosmo-workspace/cosmo/pkg/cmdutil"
)

func AddCommand(rootCmd *cobra.Command, o *cmdutil.CliOptions) {
	cmd := &cobra.Command{
		Use:   "useraddon",
		Short: "Manipulate UserAddon",
		Long: `
Manipulate UserAddon.

UserAddon is a Template for User extentions
`,
		Aliases: []string{"ua"},
	}

	cmd.AddCommand(GenerateCmd(&cobra.Command{
		Use:     "generate [< Input via Stdin or pipe]",
		Aliases: []string{"gen"},
		Short:   "Generate UserAddon Template",
		Long: `Generate UserAddon Template

For create generated template, just do "kubectl create -f cosmo-template.yaml"

Example:
  * Pipe from kustomize build and apply to your cluster in a single line 
	
      kustomize build ./kubernetes/ | cosmoctl template useraddon generate --name ADDON_NAME | kubectl apply -f -

  * Pipe from helm template and generate UserAddon Template
	
  	  helm template code-server ci/helm-chart \
		| cosmoctl template useraddon generate --name ADDON_NAME

  * Input merged config file (kustomize build ... or helm template ... etc.) and save it to file

      cosmoctl template useraddon generate --name ADDON_NAME -o cosmo-template.yaml < merged.yaml
`,
	}, o))
	cmd.AddCommand(GetCmd(&cobra.Command{
		Use:   "get [ADDON_NAMES]",
		Short: "Get UserAddons",
		Long: `Get UserAddons
List or get UserAddons in cluster.
`,
	}, o))
	rootCmd.AddCommand(cmd)
}
