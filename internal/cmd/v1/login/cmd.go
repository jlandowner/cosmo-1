package login

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
)

func AddCommand(cmd *cobra.Command, o *cli.RootOptions) {
	loginCmd := &cobra.Command{
		Use:   "login USERID",
		Short: "Login to COSMO Dashboard Server",
		Long: `
Login to COSMO Dashboard Server.
`,
	}
	cmd.AddCommand(LoginCmd(loginCmd, o))
}

type LoginOption struct {
	*cli.RootOptions

	UserName      string
	Password      string
	PasswordStdin bool
}

func LoginCmd(cmd *cobra.Command, opt *cli.RootOptions) *cobra.Command {
	o := &LoginOption{RootOptions: opt}
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().BoolVar(&o.PasswordStdin, "password-stdin", false, "input new password from stdin pipe")
	return cmd
}

func (o *LoginOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	return nil
}

func (o *LoginOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.UserName = args[0]
	}

	if o.DashboardURL == "" {
		input, err := cli.AskInput(fmt.Sprintf("Dasboard URL(%s): ", o.DashboardURL), false)
		if err != nil {
			return err
		}
		o.DashboardURL = input
	}

	if o.UserName == "" {
		input, err := cli.AskInput("UserName: ", false)
		if err != nil {
			return err
		}
		o.UserName = input
	}

	if o.PasswordStdin {
		input, err := cli.ReadFromPipedStdin()
		if err != nil {
			return fmt.Errorf("failed to read from stdin pipe: %w", err)
		}
		o.Password = input
	} else {
		input, err := cli.AskInput("Password: ", true)
		if err != nil {
			return err
		}
		o.Password = input
	}
	return nil
}

func (o *LoginOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*10)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	c := o.CosmoDashClient
	ses, err := c.GetSession(ctx, o.UserName, o.Password)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	o.CliConfig.Token = ses
	o.CliConfig.User = o.UserName

	// save session
	err = o.CliConfig.Save()
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	cmdutil.PrintfColorInfo(o.Out, "Successfully logined to %s as %s\n", o.DashboardURL, o.UserName)

	return nil

}
