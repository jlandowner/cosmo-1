package workspace

import (
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"

	"k8s.io/cli-runtime/pkg/printers"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type GetOption struct {
	*cmdutil.CliOptions

	args []string
}

func GetCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &GetOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
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
	return nil
}

func (o *GetOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.args = args
	}
	return nil
}

func (o *GetOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("get_workspacetemplates")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewTemplateServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.GetWorkspaceTemplates(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	items := res.Msg.GetItems()

	if len(o.args) > 0 {
		i := make([]*dashv1alpha1.Template, 0, len(o.args))
		for _, selected := range o.args {
			for _, v := range items {
				if selected == v.GetName() {
					i = append(i, v)
				}
			}
		}
		items = i
	}

	w := printers.GetNewTabWriter(o.Out)
	defer w.Flush()

	columnNames := []string{"NAME", "REQUIREDVARS"}
	fmt.Fprintf(w, "%s\n", strings.Join(columnNames, "\t"))
	for _, v := range items {
		reqVars := make([]string, len(v.RequiredVars))
		for i := 0; i < len(v.RequiredVars); i++ {
			reqVars[i] = v.RequiredVars[i].VarName
		}

		rowdata := []string{v.Name, fmt.Sprintf("%s", reqVars)}
		fmt.Fprintf(w, "%s\n", strings.Join(rowdata, "\t"))
	}

	return nil
}
