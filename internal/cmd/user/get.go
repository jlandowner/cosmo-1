package user

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"

	"k8s.io/cli-runtime/pkg/printers"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type GetOption struct {
	*cmdutil.CliOptions

	UserNames    []string
	OutputFormat string

	outputFormat cmdutil.GetOutputFormat
}

func GetCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &GetOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "", "output format. available: 'wide', 'json'")
	return cmd
}

func (o *GetOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *GetOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}
	o.outputFormat = cmdutil.GetOutputFormat(o.OutputFormat)
	if err := o.outputFormat.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *GetOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.UserNames = args
	}
	return nil
}

func (o *GetOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("get_users")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewUserServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.GetUsers(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	users := res.Msg.GetItems()

	if len(o.UserNames) > 0 {
		us := make([]*dashv1alpha1.User, 0, len(o.UserNames))
		for _, selected := range o.UserNames {
			for _, v := range users {
				if selected == v.GetName() {
					us = append(us, v)
				}
			}
		}
		users = us
	}

	if o.outputFormat == cmdutil.GetOutputFormatJSON {
		out, err := json.Marshal(users)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
		fmt.Fprintf(o.Out, "%s", out)
		return nil
	}

	w := printers.GetNewTabWriter(o.Out)
	defer w.Flush()

	columnNames := []string{"NAME", "ROLE", "NAMESPACE", "STATUS"}
	if o.outputFormat == cmdutil.GetOutputFormatWide {
		columnNames = append(columnNames, "DISPLAYNAME")
	}
	fmt.Fprintf(w, "%s\n", strings.Join(columnNames, "\t"))
	for _, v := range users {
		rowdata := []string{v.Name, v.Role, cosmov1alpha1.UserNamespace(v.Name), v.Status}
		if o.outputFormat == cmdutil.GetOutputFormatWide {
			rowdata = append(rowdata, v.DisplayName)
		}
		fmt.Fprintf(w, "%s\n", strings.Join(rowdata, "\t"))
	}

	return nil
}
