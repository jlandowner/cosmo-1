package workspace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type DeleteOption struct {
	*cli.RootOptions

	WorkspaceNames []string
	UserName       string
}

func DeleteCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &DeleteOption{RootOptions: cliOpt}
	cmd.RunE = cli.ConnectErrorHandler(o)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")
	return cmd
}

func (o *DeleteOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *DeleteOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.WorkspaceNames = args

	if !o.UseKubeAPI && o.UserName == "" {
		o.UserName = o.CliConfig.User
	}
	return nil
}

func (o *DeleteOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*10)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	for _, v := range o.WorkspaceNames {
		if o.UseKubeAPI {
			if err := o.DeleteWorkspaceWithKubeClient(ctx, v); err != nil {
				return err
			}
		} else {
			if err := o.DeleteWorkspaceWithDashClient(ctx, v); err != nil {
				return err
			}
		}
		cmdutil.PrintfColorInfo(o.Out, "Successfully deleted workspace %s\n", v)
	}

	return nil
}

func (o *DeleteOption) DeleteWorkspaceWithDashClient(ctx context.Context, workspaceName string) error {
	req := &dashv1alpha1.DeleteWorkspaceRequest{
		UserName: o.UserName,
		WsName:   workspaceName,
	}
	c := o.CosmoDashClient
	res, err := c.WorkspaceServiceClient.DeleteWorkspace(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.DeleteWorkspace", "res", res)

	return nil
}

func (o *DeleteOption) DeleteWorkspaceWithKubeClient(ctx context.Context, workspaceName string) error {
	c := o.KosmoClient
	if _, err := c.DeleteWorkspace(ctx, workspaceName, o.UserName); err != nil {
		return err
	}
	return nil
}
