package workspace

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	dashboardv1alpha1connect "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type CreateOption struct {
	*cmdutil.CliOptions

	WorkspaceName string
	Template      string
	UserName      string
	Vars          []string
	// DryRun        bool    //TODO

	vars map[string]string
}

func CreateCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &CreateOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name")
	cmd.Flags().StringVarP(&o.Template, "template", "t", "", "template name")
	cmd.Flags().StringArrayVar(&o.Vars, "var", nil, "template vars. format is '--var=VarName1:VarValue1 --var=VarName2:VarValue2'")
	// cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "dry run")

	return cmd
}

func (o *CreateOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *CreateOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	if o.Template == "" {
		return errors.New("--template is required")
	}
	return nil
}

func (o *CreateOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.WorkspaceName = args[0]

	if len(o.Vars) > 0 {
		vars := make(map[string]string)
		for _, v := range o.Vars {
			varAndVal := strings.Split(v, ":")
			if len(varAndVal) != 2 {
				return fmt.Errorf("vars format error: vars %s must be 'VAR:VAL'", v)
			}
			vars[varAndVal[0]] = varAndVal[1]
		}
		o.vars = vars
	}

	return nil
}

func (o *CreateOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("create_workspace")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewWorkspaceServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.CreateWorkspace(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
		&dashv1alpha1.CreateWorkspaceRequest{
			UserName: o.UserName,
			WsName:   o.WorkspaceName,
			Template: o.Template,
			Vars:     o.vars,
		}))
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	cmdutil.PrintfColorInfo(o.ErrOut, "Successfully created workspace %s\n", o.WorkspaceName)

	return nil
}
