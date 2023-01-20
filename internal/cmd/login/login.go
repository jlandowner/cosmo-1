package login

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/bufbuild/connect-go"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type LoginOption struct {
	*cmdutil.CliOptions

	Password      string
	PasswordStdin bool
	CAFile        string
}

func LoginCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &LoginOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVar(&o.LoginUser, "user", "", "login user name (required)")
	cmd.Flags().StringVar(&o.CAFile, "ca", "", "ca cert file path for server")
	cmd.Flags().StringVar(&o.Password, "password", "", "WARNING: this flag may be insecure. use --password-stdin")
	cmd.Flags().BoolVar(&o.PasswordStdin, "password-stdin", false, "input password by stdin")
	cmd.Flags().BoolVar(&o.Insecure, "insecure", false, "use http not https")
	return cmd
}

func (o *LoginOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *LoginOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("no args")
	}
	if o.LoginUser == "" {
		return errors.New("--user is required")
	}
	if (o.Password == "" && !o.PasswordStdin) || (o.Password != "" && o.PasswordStdin) {
		return errors.New("--password or --password-stdin is required")
	}
	if o.PasswordStdin {
		if isatty.IsTerminal(os.Stdin.Fd()) {
			return fmt.Errorf("no input via stdin")
		}
	}
	return nil
}

func (o *LoginOption) Complete(cmd *cobra.Command, args []string) error {
	o.ServerEndpoint = args[0]

	if o.CliConfigFilePath == "" || o.CliConfigFilePath == "$HOME/.cosmoctl" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		o.CliConfigFilePath = path.Join(home, ".cosmoctl")
	}

	// Complete logger
	o.CompleteLogger()

	if o.CAFile != "" {
		ca, err := ioutil.ReadFile(o.CAFile)
		if err != nil {
			return err
		}
		o.CliConfig.ServerCA = ca
	}

	// Complete Client
	if err := o.CompleteClient(); err != nil {
		return err
	}

	if o.Insecure {
		o.ServerEndpoint = "http://" + o.ServerEndpoint
	} else {
		o.ServerEndpoint = "https://" + o.ServerEndpoint
	}

	if o.PasswordStdin {
		// input data from stdin
		input, err := io.ReadAll(o.In)
		if err != nil {
			return fmt.Errorf("failed to read input file : %w", err)
		}
		if len(input) == 0 {
			return fmt.Errorf("no input")
		}
		o.Password = string(input)
	}
	return nil
}

func (o *LoginOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("login")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewAuthServiceClient(o.Client, o.ServerEndpoint)

	res, err := c.Login(ctx, connect.NewRequest(&dashv1alpha1.LoginRequest{
		UserName: o.LoginUser,
		Password: o.Password,
	}))
	if err != nil {
		return fmt.Errorf("failed to login server: %w", err)
	}
	log.Debug().Info(fmt.Sprintf("response: %v", res))

	o.Cookie = res.Header().Get("Cookie")

	log.Debug().Info(o.CliConfigFilePath)
	if err := o.CliOptions.CliConfig.Write(o.CliConfigFilePath); err != nil {
		return fmt.Errorf("failed to save config to %s: %w", o.CliConfigFilePath, err)
	}

	cmdutil.PrintfColorInfo(o.ErrOut, "Successfully Logined as User %s\n", o.LoginUser)

	if res.Msg.RequirePasswordUpdate {
		cmdutil.PrintfColorErr(o.ErrOut, "WARNING you should update password. Run `cosmoctl user update-password`\n")
	}

	return nil
}
