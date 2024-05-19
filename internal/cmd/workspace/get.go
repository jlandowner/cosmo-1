package workspace

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type GetOption struct {
	*cli.RootOptions

	WorkspaceNames []string
	Filter         []string
	UserName       string
	AllUsers       bool
	OutputFormat   string

	filters []cli.Filter
}

func GetCmd(cmd *cobra.Command, opt *cli.RootOptions) *cobra.Command {
	o := &GetOption{RootOptions: opt}
	cmd.RunE = cli.ConnectErrorHandler(o)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")
	cmd.Flags().StringSliceVar(&o.Filter, "filter", nil, "filter option. available columns are ['NAME', 'TEMPLATE', 'PHASE']. available operators are ['==', '!=']. value format is filepath. e.g. '--filter TEMPLATE==dev-*'")
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", "output format. available values are ['table', 'yaml', 'wide']")
	return cmd
}

func (o *GetOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	if o.UseKubeAPI && o.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	switch o.OutputFormat {
	case "table", "yaml", "wide":
	default:
		return fmt.Errorf("invalid output format: %s", o.OutputFormat)
	}
	if o.UseKubeAPI && o.UserName == "" {
		return fmt.Errorf("user name is required")
	}
	return nil
}

func (o *GetOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.WorkspaceNames = args
	}
	if o.UserName == "" {
		o.UserName = o.CliConfig.User
	}
	if len(o.Filter) > 0 {
		f, err := cli.ParseFilters(o.Filter)
		if err != nil {
			return err
		}
		o.filters = f
	}
	for _, f := range o.filters {
		o.Logr.Debug().Info("filter", "key", f.Key, "value", f.Value, "op", f.Operator)
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return nil
}

func (o *GetOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*30)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	var (
		workspaces []*dashv1alpha1.Workspace
		err        error
	)
	if o.UseKubeAPI {
		workspaces, err = o.ListWorkspacesByKubeClient(ctx, o.OutputFormat == "yaml")
	} else {
		workspaces, err = o.ListWorkspacesWithDashClient(ctx, o.OutputFormat == "yaml")
	}
	if err != nil {
		return err
	}
	o.Logr.Debug().Info("Workspaces", "workspaces", workspaces)

	workspaces = o.ApplyFilters(workspaces)

	if o.OutputFormat == "yaml" {
		o.OutputYAML(workspaces)
		return nil
	} else if o.OutputFormat == "wide" {
		OutputWideTable(o.Out, workspaces)
		return nil
	} else {
		OutputTable(o.Out, workspaces)
		return nil
	}
}

func (o *GetOption) ListWorkspacesWithDashClient(ctx context.Context, withRaw bool) ([]*dashv1alpha1.Workspace, error) {
	req := &dashv1alpha1.GetWorkspacesRequest{
		UserName: o.UserName,
		WithRaw:  &withRaw,
	}
	c := o.CosmoDashClient
	res, err := c.WorkspaceServiceClient.GetWorkspaces(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("WorkspaceServiceClient.GetWorkspaces", "res", res)
	return res.Msg.Items, nil
}

func (o *GetOption) ApplyFilters(workspaces []*dashv1alpha1.Workspace) []*dashv1alpha1.Workspace {
	for _, f := range o.filters {
		o.Logr.Debug().Info("applying filter", "key", f.Key, "value", f.Value, "op", f.Operator)

		switch strings.ToUpper(f.Key) {
		case "NAME":
			workspaces = cli.DoFilter(workspaces, func(u *dashv1alpha1.Workspace) []string {
				return []string{u.Name}
			}, f)
		case "TEMPLATE":
			workspaces = cli.DoFilter(workspaces, func(u *dashv1alpha1.Workspace) []string {
				return []string{u.Spec.Template}
			}, f)
		case "PHASE":
			workspaces = cli.DoFilter(workspaces, func(u *dashv1alpha1.Workspace) []string {
				return []string{u.Status.Phase}
			}, f)
		default:
			o.Logr.Info("WARNING: unknown filter key", "key", f.Key)
		}
	}

	if len(o.WorkspaceNames) > 0 {
		ts := make([]*dashv1alpha1.Workspace, 0, len(o.WorkspaceNames))
	WorkspaceLoop:
		// Or loop
		for _, t := range workspaces {
			for _, selected := range o.WorkspaceNames {
				if selected == t.GetName() {
					ts = append(ts, t)
					continue WorkspaceLoop
				}
			}
		}
		workspaces = ts
	}
	return workspaces
}

func (o *GetOption) OutputYAML(objs []*dashv1alpha1.Workspace) {
	docs := make([]string, len(objs))
	for i, t := range objs {
		docs[i] = *t.Raw
	}
	fmt.Fprintln(o.Out, strings.Join(docs, "---\n"))
}

func OutputTable(out io.Writer, workspaces []*dashv1alpha1.Workspace) {
	data := [][]string{}

	for _, v := range workspaces {
		data = append(data, []string{v.OwnerName, v.Name, v.Spec.Template, v.Status.Phase, v.Status.MainUrl})
	}

	cli.OutputTable(out,
		[]string{"USER", "NAME", "TEMPLATE", "PHASE", "MAINURL"},
		data)
}

func OutputWideTable(out io.Writer, workspaces []*dashv1alpha1.Workspace) {
	data := [][]string{}

	for _, v := range workspaces {
		vars := make([]string, 0, len(v.Spec.Vars))
		for k, vv := range v.Spec.Vars {
			vars = append(vars, fmt.Sprintf("%s=%s", k, vv))
		}
		data = append(data, []string{v.OwnerName, v.Name, v.Spec.Template, strings.Join(vars, ","), v.Status.Phase, v.Status.MainUrl})
	}

	cli.OutputTable(out,
		[]string{"USER", "NAME", "TEMPLATE", "VARS", "PHASE", "MAINURL"},
		data)
}

func (o *GetOption) ListWorkspacesByKubeClient(ctx context.Context, withRaw bool) ([]*dashv1alpha1.Workspace, error) {
	c := o.KosmoClient
	workspaces, err := c.ListWorkspacesByUserName(ctx, o.UserName)
	if err != nil {
		return nil, err
	}
	return apiconv.C2D_Workspaces(workspaces, apiconv.WithWorkspaceRaw(&withRaw)), nil
}
