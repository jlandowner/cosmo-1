package workspace

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type GetOption struct {
	*cli.RootOptions

	WorkspaceNames []string
	Filter         []string
	UserName       string
	AllUsers       bool

	tmplFilter []string
}

func GetCmd(cmd *cobra.Command, opt *cli.RootOptions) *cobra.Command {
	o := &GetOption{RootOptions: opt}
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVarP(&o.UserName, "user", "u", "", "user name (defualt: login user)")
	cmd.Flags().StringSliceVar(&o.Filter, "filter", nil, "filter option. 'template' is available for now. e.g. 'template=x'")
	return cmd
}

func (o *GetOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
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
		for _, f := range o.Filter {
			s := strings.Split(f, "=")
			if len(s) != 2 {
				return fmt.Errorf("invalid filter expression: %s", f)
			}
			switch s[0] {
			case "template", "tmpl":
				o.tmplFilter = append(o.tmplFilter, s[1])
			default:
				o.Logr.Info("invalid filter expression", "filter", f)
				return fmt.Errorf("invalid filter expression: %s", f)
			}
		}
	}
	o.Logr.Debug().Info("filter", "tmpl", o.tmplFilter)
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

	var workspaces []*dashv1alpha1.Workspace
	var err error
	if o.UseKubeAPI {
		workspaces, err = o.ListWorkspacesByKubeClient(ctx)
		if err != nil {
			return err
		}
	} else {
		workspaces, err = o.ListWorkspacesWithDashClient(ctx)
		if err != nil {
			return err
		}
	}
	o.Logr.Debug().Info("Workspaces", "workspaces", workspaces)

	workspaces = o.ApplyFilters(workspaces)

	o.Output(workspaces)

	return nil

}

func (o *GetOption) ListWorkspacesWithDashClient(ctx context.Context) ([]*dashv1alpha1.Workspace, error) {
	req := &dashv1alpha1.GetWorkspacesRequest{
		UserName: o.UserName,
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
	if len(o.tmplFilter) > 0 {
		// And loop
		for _, selected := range o.tmplFilter {
			ts := make([]*dashv1alpha1.Workspace, 0)
			for _, t := range workspaces {
				if matched, err := filepath.Match(selected, t.Spec.Template); err == nil && matched {
					ts = append(ts, t)
				}
			}
			workspaces = ts
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

func (o *GetOption) Output(workspaces []*dashv1alpha1.Workspace) {
	data := [][]string{}

	for _, v := range workspaces {
		data = append(data, []string{v.OwnerName, v.Name, v.Spec.Template, v.Status.Phase})
	}

	cli.OutputTable(o.Out,
		[]string{"USER", "NAME", "TEMPLATE", "PHASE"},
		data)
}

func (o *GetOption) ListWorkspacesByKubeClient(ctx context.Context) ([]*dashv1alpha1.Workspace, error) {
	c := o.KosmoClient
	workspaces, err := c.ListWorkspacesByUserName(ctx, o.UserName)
	if err != nil {
		return nil, err
	}
	return apiconv.C2D_Workspaces(workspaces), nil
}
