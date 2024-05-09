package workspace

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type GetNetworkOption struct {
	*cli.RootOptions

	WorkspaceName string
	UserName      string
}

func GetNetworkCmd(cmd *cobra.Command, opt *cli.RootOptions) *cobra.Command {
	o := &GetNetworkOption{RootOptions: opt}
	cmd.RunE = cli.ConnectErrorHandler(o)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")
	return cmd
}

func (o *GetNetworkOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	return nil
}

func (o *GetNetworkOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.WorkspaceName = args[0]
	} else {
		o.WorkspaceName = cli.GetCurrentWorkspaceName()
		o.Logr.Info("Workspace name is auto detected from hostname", "name", o.WorkspaceName)
	}
	if !o.UseKubeAPI && o.UserName == "" {
		o.UserName = o.CliConfig.User
	}
	return nil
}

func (o *GetNetworkOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*30)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	var workspace *dashv1alpha1.Workspace
	var err error
	if o.UseKubeAPI {
		workspace, err = o.GetWorkspaceByKubeClient(ctx)
		if err != nil {
			return err
		}
	} else {
		workspace, err = o.GetWorkspaceWithDashClient(ctx)
		if err != nil {
			return err
		}
	}
	o.Logr.Debug().Info("Workspace", "workspace", workspace)

	o.Output(workspace)

	return nil

}

func (o *GetNetworkOption) GetWorkspaceWithDashClient(ctx context.Context) (*dashv1alpha1.Workspace, error) {
	req := &dashv1alpha1.GetWorkspaceRequest{
		WsName:   o.WorkspaceName,
		UserName: o.UserName,
	}
	c := o.CosmoDashClient
	res, err := c.WorkspaceServiceClient.GetWorkspace(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.GetWorkspace", "res", res)
	return res.Msg.Workspace, nil
}

func (o *GetNetworkOption) Output(workspace *dashv1alpha1.Workspace) {
	data := [][]string{}

	for _, v := range workspace.Spec.Network {
		data = append(data, []string{fmt.Sprintf("%d", v.PortNumber), v.CustomHostPrefix, v.HttpPath, strconv.FormatBool(v.Public), v.Url})
	}

	cli.OutputTable(o.Out,
		[]string{"PORT", "CUSTOM_HOST_PREFIX", "HTTP_PATH", "PUBLIC", "URL"},
		data)
}

func (o *GetNetworkOption) GetWorkspaceByKubeClient(ctx context.Context) (*dashv1alpha1.Workspace, error) {
	c := o.KosmoClient
	workspace, err := c.GetWorkspaceByUserName(ctx, o.WorkspaceName, o.UserName)
	if err != nil {
		return nil, err
	}
	return apiconv.C2D_Workspace(*workspace), nil
}
