package workspace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/utils/pointer"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
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

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)

	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")

	return cmd
}

func (o *ResumeOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *ResumeOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *ResumeOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.WorkspaceNames = args

	if o.UserName == "" {
		o.UserName = o.CliConfig.User
	}
	return nil
}

func (o *ResumeOption) RunE(cmd *cobra.Command, args []string) error {
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
		cmdutil.PrintfColorInfo(o.Out, "Successfully resumed workspace %s\n", v)
	}

	return nil
}

func (o *ResumeOption) ResumeWorkspaceWithDashClient(ctx context.Context, workspaceName string) error {
	req := &dashv1alpha1.UpdateWorkspaceRequest{
		UserName: o.UserName,
		WsName:   workspaceName,
		Replicas: pointer.Int64(1),
	}
	c := o.CosmoDashClient
	res, err := c.WorkspaceServiceClient.UpdateWorkspace(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.UpdateWorkspace", "res", res)

	return nil
}

func (o *ResumeOption) ResumeWorkspaceWithKubeClient(ctx context.Context, workspaceName string) error {
	c := o.KosmoClient
	if _, err := c.UpdateWorkspace(ctx, workspaceName, o.UserName, kosmo.UpdateWorkspaceOpts{Replicas: pointer.Int64(1)}); err != nil {
		return err
	}
	return nil
}