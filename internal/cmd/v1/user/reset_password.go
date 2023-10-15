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
)

type resetPasswordOption struct {
	*changePasswordOption

	Force bool
}

func resetPasswordCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &resetPasswordOption{changePasswordOption: &changePasswordOption{RootOptions: cliOpt}}
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().BoolVar(&o.Force, "force", false, "not ask confirmation")
	return cmd
}

func (o *resetPasswordOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	if !o.UseKubeAPI {
		return errors.New("force reset is only available with -k")
	}
	return nil
}

func (o *resetPasswordOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.UserName = args[0]

	return nil
}

func (o *resetPasswordOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*10)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	if err := o.ValidateUser(ctx); err != nil {
		return err
	}

	newPassword, err := o.resetPasswordWithKubeClient(ctx)
	if err != nil {
		return err
	}

	cmdutil.PrintfColorInfo(o.Out, "Successfully reset password: user %s\n", o.UserName)

	fmt.Fprintln(o.Out, "New password:", *newPassword)

	return nil
}

func (o *resetPasswordOption) resetPasswordWithKubeClient(ctx context.Context) (*string, error) {
	c := o.KosmoClient
	if err := c.ResetPassword(ctx, o.UserName); err != nil {
		return nil, err
	}
	pass, err := c.GetDefaultPassword(ctx, o.UserName)
	if err != nil {
		return nil, err
	}
	if pass == nil {
		return nil, errors.New("password is nil")
	}
	return pass, nil
}
