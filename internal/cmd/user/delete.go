package user

import (
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type DeleteOption struct {
	*cmdutil.CliOptions

	UserName string
}

func DeleteCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &DeleteOption{CliOptions: cliOpt}
	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	return cmd
}

func (o *DeleteOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *DeleteOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *DeleteOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.UserName = args[0]
	return nil
}

func (o *DeleteOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("delete_user")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewUserServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.DeleteUser(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
		&dashv1alpha1.DeleteUserRequest{
			UserName: o.UserName,
		}))
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	cmdutil.PrintfColorInfo(o.Out, "Successfully deleted user %s\n", o.UserName)
	return nil
}
