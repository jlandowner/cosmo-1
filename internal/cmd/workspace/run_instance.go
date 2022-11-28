package workspace

import (
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
	"k8s.io/utils/pointer"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	dashboardv1alpha1connect "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type RunInstanceOption struct {
	*cmdutil.CliOptions

	InstanceName string
	UserName     string
}

func RunInstanceCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &RunInstanceOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name")
	return cmd
}

func (o *RunInstanceOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *RunInstanceOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *RunInstanceOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.InstanceName = args[0]
	return nil
}

func (o *RunInstanceOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("run_instance")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewWorkspaceServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.UpdateWorkspace(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
		&dashv1alpha1.UpdateWorkspaceRequest{
			UserName: o.UserName,
			WsName:   o.InstanceName,
			Replicas: pointer.Int64(1),
		}))
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	cmdutil.PrintfColorInfo(o.Out, "Successfully run workspace %s\n", o.InstanceName)
	return nil
}
