package user

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

	UserName string
}

func DeleteCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &DeleteOption{RootOptions: cliOpt}
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
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
	o.UserName = args[0]
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

	if o.UseKubeAPI {
		if err := o.DeleteUserWithKubeClient(ctx); err != nil {
			return err
		}
	} else {
		if err := o.DeleteUserWithDashClient(ctx); err != nil {
			return err
		}
	}

	cmdutil.PrintfColorInfo(o.Out, "Successfully deleted user %s\n", o.UserName)
	return nil
}

func (o *DeleteOption) DeleteUserWithDashClient(ctx context.Context) error {
	req := &dashv1alpha1.DeleteUserRequest{
		UserName: o.UserName,
	}
	c := o.CosmoDashClient
	res, err := c.UserServiceClient.DeleteUser(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("UserServiceClient.DeleteUser", "res", res)

	return nil
}

func (o *DeleteOption) DeleteUserWithKubeClient(ctx context.Context) error {
	c := o.KosmoClient
	if _, err := c.DeleteUser(ctx, o.UserName); err != nil {
		return err
	}
	return nil
}
