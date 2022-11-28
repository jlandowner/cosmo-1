package user

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
)

type resetPasswordOption struct {
	*cmdutil.CliOptions

	UserName string
}

func resetPasswordCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &resetPasswordOption{CliOptions: cliOpt}
	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	return cmd
}

func (o *resetPasswordOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *resetPasswordOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *resetPasswordOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.UserName = args[0]
	return nil
}

func (o *resetPasswordOption) RunE(cmd *cobra.Command, args []string) error {
	// log := o.Logr.WithName("reset_password")
	// ctx := clog.IntoContext(o.Ctx, log)

	// c := dashboardv1alpha1connect.NewUserServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	// // res, err := c.UpdateUserPassword(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
	// &dashv1alpha1.UpdateUserPasswordRequest{
	// // 	UserName:    o.UserName,
	// // 	CurrentPassword: ,
	// // }))
	// if err != nil {
	// 	return err
	// }
	// log.Debug().Info("response: %v", res)

	// fmt.Fprintln(o.Out, "New password:", *pass)

	return nil
}
