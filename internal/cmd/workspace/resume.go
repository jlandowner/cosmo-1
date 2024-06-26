package workspace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/utils/ptr"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/kosmo"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type ResumeOption struct {
	*cli.RootOptions

	WorkspaceNames []string
	UserName       string
}

func ResumeCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &ResumeOption{RootOptions: cliOpt}

	cmd.RunE = cli.ConnectErrorHandler(o)

	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")

	return cmd
}

func (o *ResumeOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	if o.UseKubeAPI && o.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (o *ResumeOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.WorkspaceNames = args

	if !o.UseKubeAPI && o.UserName == "" {
		o.UserName = o.CliConfig.User
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return nil
}

func (o *ResumeOption) RunE(cmd *cobra.Command, args []string) error {
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
			if err := o.ResumeWorkspaceWithKubeClient(ctx, v); err != nil {
				return err
			}
		} else {
			if err := o.ResumeWorkspaceWithDashClient(ctx, v); err != nil {
				return err
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("Successfully resumed workspace %s", v))
	}

	return nil
}

func (o *ResumeOption) ResumeWorkspaceWithDashClient(ctx context.Context, workspaceName string) error {
	req := &dashv1alpha1.UpdateWorkspaceRequest{
		UserName: o.UserName,
		WsName:   workspaceName,
		Replicas: ptr.To(int64(1)),
	}
	c := o.CosmoDashClient
	o.Logr.DebugAll().Info("WorkspaceServiceClient.UpdateWorkspace", "req", req)
	res, err := c.WorkspaceServiceClient.UpdateWorkspace(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.UpdateWorkspace", "res", res)

	return nil
}

func (o *ResumeOption) ResumeWorkspaceWithKubeClient(ctx context.Context, workspaceName string) error {
	c := o.KosmoClient
	if _, err := c.UpdateWorkspace(ctx, workspaceName, o.UserName, kosmo.UpdateWorkspaceOpts{Replicas: ptr.To(int64(1))}); err != nil {
		return err
	}
	return nil
}
