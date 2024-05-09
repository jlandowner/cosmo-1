package workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type CreateOption struct {
	*cli.RootOptions

	WorkspaceName string
	UserName      string
	Template      string
	RawVars       string

	vars map[string]string
}

func CreateCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &CreateOption{RootOptions: cliOpt}
	cmd.RunE = cli.ConnectErrorHandler(o)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")
	cmd.Flags().StringVarP(&o.Template, "template", "t", "", "template name (Required)")
	cmd.MarkFlagRequired("template")
	cmd.Flags().StringVar(&o.RawVars, "vars", "", "template vars. the format is VarName:VarValue. also it can be set multiple vars by conma separated list. (example: VAR1:VAL1,VAR2:VAL2)")

	return cmd
}

func (o *CreateOption) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *CreateOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.WorkspaceName = args[0]

	if !o.UseKubeAPI && o.UserName == "" {
		o.UserName = o.CliConfig.User
	}

	if o.RawVars != "" {
		vars := make(map[string]string)
		varAndVals := strings.Split(o.RawVars, ",")
		for _, v := range varAndVals {
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
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*10)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	var err error
	if o.UseKubeAPI {
		err = o.CreateWorkspaceWithKubeClient(ctx)
		if err != nil {
			return err
		}
	} else {
		err = o.CreateWorkspaceWithDashClient(ctx)
		if err != nil {
			return err
		}
	}
	cmdutil.PrintfColorInfo(o.Out, "Successfully created workspace %s\n", o.WorkspaceName)

	return nil
}

func (o *CreateOption) CreateWorkspaceWithDashClient(ctx context.Context) error {
	req := &dashv1alpha1.CreateWorkspaceRequest{
		WsName:   o.WorkspaceName,
		UserName: o.UserName,
		Template: o.Template,
		Vars:     o.vars,
	}
	c := o.CosmoDashClient
	res, err := c.WorkspaceServiceClient.CreateWorkspace(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.CreateWorkspace", "res", res)

	return nil
}

func (o *CreateOption) CreateWorkspaceWithKubeClient(ctx context.Context) error {
	c := o.KosmoClient
	if _, err := c.CreateWorkspace(ctx, o.UserName, o.WorkspaceName, o.Template, o.vars); err != nil {
		return err
	}
	return nil
}
