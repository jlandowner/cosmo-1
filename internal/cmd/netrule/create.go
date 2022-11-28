package netrule

import (
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
	"k8s.io/utils/pointer"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashboardv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type CreateOption struct {
	*cmdutil.CliOptions

	WorkspaceName string
	UserName      string
	NetRuleName   string
	PortNumber    int32
	Group         string
	HTTPPath      string
	Public        bool

	rule cosmov1alpha1.NetworkRule
}

func CreateCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &CreateOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVar(&o.WorkspaceName, "workspace", "", "workspace name (Required)")
	cmd.Flags().Int32Var(&o.PortNumber, "port", 0, "serivce port number (Required)")
	cmd.Flags().StringVar(&o.UserName, "user", "", "user name")
	cmd.Flags().StringVar(&o.Group, "group", "", "group of ports for URLVar. Ports in the same group are treated as the same domain. set 'name' value if empty")
	cmd.Flags().StringVar(&o.HTTPPath, "path", "/", "path for Ingress path when using ingress")
	cmd.Flags().BoolVar(&o.Public, "public", false, "disable authentication for this port")

	return cmd
}

func (o *CreateOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *CreateOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	if o.WorkspaceName == "" {
		return errors.New("--workspace is required")
	}
	if o.PortNumber == 0 {
		return errors.New("--port is required")
	}
	return nil
}

func (o *CreateOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.NetRuleName = args[0]

	if o.Group == "" {
		o.Group = o.NetRuleName
	}

	o.rule = cosmov1alpha1.NetworkRule{
		Name:       o.NetRuleName,
		PortNumber: o.PortNumber,
		HTTPPath:   o.HTTPPath,
		Group:      pointer.String(o.Group),
		Public:     o.Public,
	}
	return nil
}

func (o *CreateOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("create_networkrule")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewWorkspaceServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.UpsertNetworkRule(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
		&dashboardv1alpha1.UpsertNetworkRuleRequest{
			UserName: o.UserName,
			WsName:   o.WorkspaceName,
			NetworkRule: &dashboardv1alpha1.NetworkRule{
				Name:       o.NetRuleName,
				PortNumber: o.PortNumber,
				Group:      o.Group,
				HttpPath:   o.HTTPPath,
				Public:     o.Public,
			},
		}))
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	cmdutil.PrintfColorInfo(o.Out, "Successfully upserted network rule '%s' for workspace '%s'\n", o.NetRuleName, o.WorkspaceName)
	return nil
}
