package workspace

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	dashboardv1alpha1connect "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type GetOption struct {
	*cmdutil.CliOptions

	WorkspaceNames []string
	UserName       string
	OutputFormat   string

	outputFormat cmdutil.GetOutputFormat
}

func GetCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &GetOption{CliOptions: cliOpt}

	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name")
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "", cmdutil.HelpGetOutputFormat)
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
		o.WorkspaceNames = args
	}
	return nil
}

func (o *GetOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("get_workspaces")
	ctx := clog.IntoContext(o.Ctx, log)

	c := dashboardv1alpha1connect.NewWorkspaceServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

	res, err := c.GetWorkspaces(ctx, cmdutil.NewConnectRequestWithAuth(o.Token,
		&dashv1alpha1.GetWorkspacesRequest{
			UserName: o.UserName,
		}))
	if err != nil {
		return err
	}
	log.Debug().Info("response: %v", res)

	wss := res.Msg.GetItems()

	if len(o.WorkspaceNames) > 0 {
		ws := make([]*dashv1alpha1.Workspace, 0, len(o.WorkspaceNames))
		for _, selected := range o.WorkspaceNames {
			for _, v := range wss {
				if selected == v.GetName() {
					ws = append(ws, v)
				}
			}
		}
		wss = ws
	}

	if o.outputFormat == cmdutil.GetOutputFormatJSON {
		out, err := json.Marshal(wss)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
		fmt.Fprintf(o.Out, "%s", out)
		return nil
	}

	w := printers.GetNewTabWriter(o.Out)
	defer w.Flush()

	columnNames := []string{"NAME", "TEMPLATE", "PODPHASE"}
	if o.outputFormat == cmdutil.GetOutputFormatWide {
		columnNames = append(columnNames, "URLS")
	}
	fmt.Fprintf(w, "%s\n", strings.Join(columnNames, "\t"))

	for _, ws := range wss {
		rowdata := []string{ws.GetName(), ws.GetSpec().GetTemplate(), ws.Status.GetPhase()}
		if o.outputFormat == cmdutil.GetOutputFormatWide {
			rowdata = append(rowdata, fmt.Sprintf("%s", ws.GetStatus().GetMainUrl()))
		}
		fmt.Fprintf(w, "%s\n", strings.Join(rowdata, "\t"))
	}
	return nil
}
