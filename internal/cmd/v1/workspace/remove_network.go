package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type RemoveNetworkOption struct {
	*cli.RootOptions

	WorkspaceName    string
	UserName         string
	CustomHostPrefix string
	PortNumber       int32
	HTTPPath         string
	Public           bool

	rule cosmov1alpha1.NetworkRule
}

func RemoveNetworkCmd(cmd *cobra.Command, cliOpt *cli.RootOptions) *cobra.Command {
	o := &RemoveNetworkOption{RootOptions: cliOpt}

	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")
	cmd.Flags().Int32Var(&o.PortNumber, "port", 0, "serivce port number")
	cmd.Flags().StringVar(&o.CustomHostPrefix, "custom-host-prefix", "", "custom host prefix")
	cmd.Flags().StringVar(&o.HTTPPath, "path", "/", "path for Ingress path when using ingress")

	return cmd
}

func (o *RemoveNetworkOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	return nil
}

func (o *RemoveNetworkOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.WorkspaceName = args[0]
	} else {
		// o.WorkspaceName = GetCurrentWorkspaceName()
	}
	if o.UserName == "" {
		o.UserName = o.CliConfig.User
	}

	o.rule = cosmov1alpha1.NetworkRule{
		CustomHostPrefix: o.CustomHostPrefix,
		PortNumber:       o.PortNumber,
		HTTPPath:         o.HTTPPath,
		Public:           o.Public,
	}
	o.rule.Default()
	return nil
}

func (o *RemoveNetworkOption) RunE(cmd *cobra.Command, args []string) error {
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
		err := o.DeleteNetworkRuleByKubeClient(ctx)
		if err != nil {
			return err
		}
	} else {
		err := o.DeleteNetworkRuleWithDashClient(ctx)
		if err != nil {
			return err
		}
	}

	cmdutil.PrintfColorInfo(o.Out, "Successfully upsert network rule for workspace '%s'\n", o.WorkspaceName)
	return nil
}

func (o *RemoveNetworkOption) DeleteNetworkRuleWithDashClient(ctx context.Context) error {
	reqGet := &dashv1alpha1.GetWorkspaceRequest{
		WsName:   o.WorkspaceName,
		UserName: o.UserName,
	}
	c := o.CosmoDashClient
	resGet, err := c.WorkspaceServiceClient.GetWorkspace(ctx, cli.NewRequestWithToken(reqGet, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}

	rules := apiconv.D2C_NetworkRules(resGet.Msg.Workspace.Spec.Network)
	index := cosmov1alpha1.GetNetworkRuleIndex(rules, o.rule)

	if index < 0 || len(resGet.Msg.Workspace.Spec.Network) <= index {
		return fmt.Errorf("network rule is not found: %#v", o.rule)
	}

	req := &dashv1alpha1.DeleteNetworkRuleRequest{
		WsName:   o.WorkspaceName,
		UserName: o.UserName,
		Index:    int32(index),
	}
	res, err := c.WorkspaceServiceClient.DeleteNetworkRule(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.DeleteNetworkRule", "res", res)
	return nil
}

func (o *RemoveNetworkOption) DeleteNetworkRuleByKubeClient(ctx context.Context) error {
	c := o.KosmoClient

	ws, err := c.GetWorkspaceByUserName(ctx, o.WorkspaceName, o.UserName)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %v", err)
	}
	index := cosmov1alpha1.GetNetworkRuleIndex(ws.Spec.Network, o.rule)

	if index < 0 || len(ws.Spec.Network) <= index {
		return fmt.Errorf("network rule is not found: %#v", o.rule)
	}

	if _, err := c.DeleteNetworkRule(ctx, o.WorkspaceName, o.UserName, index); err != nil {
		return err
	}

	return nil
}
