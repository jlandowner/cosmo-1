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

type StopInstanceOption struct {
	*cmdutil.CliOptions

	InstanceName string
	UserName     string
}

func StopInstanceCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &StopInstanceOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	return cmd
}

func (o *StopInstanceOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *StopInstanceOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("invalid args")
	}
	return nil
}

func (o *StopInstanceOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	o.InstanceName = args[0]
	return nil
}

func (o *StopInstanceOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("stop_instance")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewWorkspaceServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.UpdateWorkspace(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
		&dashv1alpha1.UpdateWorkspaceRequest{
			UserName: o.UserName,
			WsName:   o.InstanceName,
			Replicas: pointer.Int64(0),
		}))
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	cmdutil.PrintfColorInfo(o.Out, "Successfully stopped workspace %s\n", o.InstanceName)
	return nil
}
