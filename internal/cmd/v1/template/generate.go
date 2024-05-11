package template

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/internal/cmd/version"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	cmdutil "github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	"github.com/cosmo-workspace/cosmo/pkg/template"
)

type generateOption struct {
	*cli.RootOptions
	wsConfig cosmov1alpha1.Config

	Name         string
	OutputFile   string
	RequiredVars []string
	Desc         string

	TypeWorkspace bool
	TypeUserAddon bool

	SetDefaultUserAddon bool
	DisableNamePrefix   bool
	ClusterScope        bool
	UserRoles           []string
	RequiredUserAddons  []string

	NoHeader bool
}

func generateCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &generateOption{RootOptions: cliOpt}
	cmd.RunE = cli.ConnectErrorHandler(o)

	cmd.Flags().StringVarP(&o.Name, "name", "n", "", "template name (use directory name if not specified)")
	cmd.Flags().StringVarP(&o.OutputFile, "output", "o", "", "write output into file (default: Stdout)")
	cmd.Flags().StringSliceVar(&o.RequiredVars, "required-vars", []string{}, "template custom vars to be replaced by instance. format --required-vars VAR1,VAR2:default-value")
	cmd.Flags().StringVar(&o.Desc, "desc", "", "template description")

	cmd.Flags().BoolVar(&o.TypeWorkspace, "workspace", false, "template as type workspace")
	cmd.Flags().StringVar(&o.wsConfig.DeploymentName, "workspace-deployment-name", "", "Deployment name for Workspace. use with --workspace (auto detected if not specified)")
	cmd.Flags().StringVar(&o.wsConfig.ServiceName, "workspace-service-name", "", "Service name for Workspace. use with --workspace (auto detected if not specified)")
	cmd.Flags().StringVar(&o.wsConfig.ServiceMainPortName, "workspace-main-service-port-name", "", "ServicePort name for Workspace main container port. use with --workspace (auto detected if not specified)")

	cmd.Flags().BoolVar(&o.TypeUserAddon, "useraddon", false, "template as type useraddon")
	cmd.Flags().BoolVar(&o.SetDefaultUserAddon, "useraddon-set-default", false, "set default user addon")

	cmd.Flags().BoolVar(&o.DisableNamePrefix, "disable-nameprefix", false, "disable adding instance name prefix on child resource name")
	cmd.Flags().BoolVar(&o.ClusterScope, "cluster-scope", false, "generate ClusterTemplate (default generate namespaced Template)")
	cmd.Flags().StringSliceVar(&o.UserRoles, "userroles", []string{}, "user roles to show this template (e.g. 'teama-*', 'teamb-admin', etc.)")
	cmd.Flags().StringSliceVar(&o.RequiredUserAddons, "required-useraddons", []string{}, "required user addons")

	cmd.Flags().BoolVar(&o.NoHeader, "no-header", false, "no output headers")

	return cmd
}

func (o *generateOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}

	if o.TypeWorkspace && o.TypeUserAddon {
		return errors.New("--workspace and --useraddon cannot be specified concurrently")
	}

	if o.TypeWorkspace && o.ClusterScope {
		return errors.New("workspace template cannot be cluster-scoped")
	}

	return nil
}

func (o *generateOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.CompleteWithoutClient(cmd, args); err != nil {
		return err
	}
	if o.Name == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		o.Name = filepath.Base(dir)
	}

	if o.OutputFile != "" {
		var err error
		o.OutputFile, err = filepath.Abs(o.OutputFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *generateOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	input, err := cli.ReadFromPipedStdin()
	if err != nil {
		return err
	}
	o.Logr.Debug().Info(input)
	if o.TypeWorkspace {
		unsts, err := preTemplateBuild(input)
		if err != nil {
			return fmt.Errorf("failed to pre-build template: %w", err)
		}
		if err := completeWorkspaceConfig(&o.wsConfig, unsts); err != nil {
			return fmt.Errorf("type workspace validation failed: %w", err)
		}
	}

	builder := NewTemplateObjectBuilder(o.ClusterScope).
		Name(o.Name).
		Description(o.Desc).
		RequiredVars(o.RequiredVars).
		SetRequiredAddons(o.RequiredUserAddons).
		SetUserRoles(o.UserRoles)

	if o.TypeWorkspace {
		builder = builder.TypeWorkspace(o.wsConfig)
	} else if o.TypeUserAddon {
		builder = builder.TypeUserAddon(o.SetDefaultUserAddon)
	}
	if o.DisableNamePrefix {
		builder = builder.DisableNamePrefix()
	}

	// exec kustomize to modify input
	rawYAML, err := ExecKustomizeToRawYAML(o.Ctx, builder.KustomizeBuilder(), input)
	if err != nil {
		return err
	}

	out, err := builder.RawYAML(rawYAML).Build()
	if err != nil {
		return err
	}

	return o.Output(out)
}

func (o *generateOption) Output(output []byte) error {
	if !o.NoHeader {
		output = append([]byte("# Generated by "+version.Footprint+"\n"), output...)
	}

	// output to Stdout or write the output to file given by option
	if o.OutputFile == "" {
		fmt.Fprintln(o.Out, string(output))
	} else {
		if err := cmdutil.CreateFile(filepath.Dir(o.OutputFile), filepath.Base(o.OutputFile), output); err != nil {
			return err
		}
	}
	return nil
}

func preTemplateBuild(rawTmpl string) ([]unstructured.Unstructured, error) {
	var inst cosmov1alpha1.Instance
	inst.SetName("dummy")
	inst.SetNamespace("dummy")

	builder := template.NewRawYAMLBuilder(rawTmpl, &inst)
	return builder.Build()
}