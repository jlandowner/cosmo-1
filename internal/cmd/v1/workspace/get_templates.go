package workspace

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type GetTemplatesOption struct {
	*cli.RootOptions
	TemplateNames []string

	Filter []string

	roleFilter []string
	showDetail bool
}

func GetTemplatesCmd(cmd *cobra.Command, opt *cli.RootOptions) *cobra.Command {
	o := &GetTemplatesOption{RootOptions: opt}
	cmd.RunE = cli.ConnectErrorHandler(o)
	cmd.Flags().StringSliceVar(&o.Filter, "filter", nil, "filter option. 'userrole' is available for now. e.g. 'userrole=x'")
	return cmd
}

func (o *GetTemplatesOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Validate(cmd, args); err != nil {
		return err
	}
	return nil
}

func (o *GetTemplatesOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.RootOptions.Complete(cmd, args); err != nil {
		return err
	}
	if len(args) > 0 {
		o.TemplateNames = args
	}
	if len(args) == 1 {
		o.showDetail = true
	}
	if len(o.Filter) > 0 {
		for _, f := range o.Filter {
			s := strings.Split(f, "=")
			if len(s) != 2 {
				return fmt.Errorf("invalid filter expression: %s", f)
			}
			switch s[0] {
			case "userrole":
				o.roleFilter = append(o.roleFilter, s[1])
			default:
				o.Logr.Info("invalid filter expression", "filter", f)
				return fmt.Errorf("invalid filter expression: %s", f)
			}
		}
	}
	o.Logr.Debug().Info("filter", "role", o.roleFilter)
	return nil
}

func (o *GetTemplatesOption) RunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	ctx, cancel := context.WithTimeout(o.Ctx, time.Second*30)
	defer cancel()
	ctx = clog.IntoContext(ctx, o.Logr)

	var tmpls []*dashv1alpha1.Template
	var err error
	if o.UseKubeAPI {
		tmpls, err = o.ListWorkspaceTemplatesByKubeClient(ctx)
		if err != nil {
			return err
		}
	} else {
		tmpls, err = o.ListWorkspaceTemplatesWithDashClient(ctx)
		if err != nil {
			return err
		}
	}
	o.Logr.Debug().Info("WorkspaceTemplate templates", "templates", tmpls)

	tmpls = o.ApplyFilters(tmpls)

	if o.showDetail {
		if len(tmpls) == 0 {
			return fmt.Errorf("template not found")
		} else {
			o.OutputDetail(tmpls[0])
			return nil
		}
	} else {
		o.Output(tmpls)
	}

	return nil

}

func (o *GetTemplatesOption) ListWorkspaceTemplatesWithDashClient(ctx context.Context) ([]*dashv1alpha1.Template, error) {
	req := &dashv1alpha1.GetWorkspaceTemplatesRequest{
		UseRoleFilter: pointer.Bool(false),
	}
	c := o.CosmoDashClient
	res, err := c.TemplateServiceClient.GetWorkspaceTemplates(ctx, cli.NewRequestWithToken(req, o.CliConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("TemplateServiceClient.GetWorkspaceTemplates", "res", res)
	return res.Msg.Items, nil
}

func (o *GetTemplatesOption) ApplyFilters(tmpls []*dashv1alpha1.Template) []*dashv1alpha1.Template {
	// filter userroles
	if len(o.roleFilter) > 0 {
		// And loop
		for _, selected := range o.roleFilter {
			ts := make([]*dashv1alpha1.Template, 0)
			for _, t := range tmpls {
			RoleFilterLoop:
				for _, v := range t.Userroles {
					if matched, err := filepath.Match(selected, v); err == nil && matched {
						ts = append(ts, t)
						break RoleFilterLoop
					}
				}
			}
			tmpls = ts
		}
	}

	if len(o.TemplateNames) > 0 {
		ts := make([]*dashv1alpha1.Template, 0, len(o.TemplateNames))
	WorkspaceLoop:
		// Or loop
		for _, t := range tmpls {
			for _, selected := range o.TemplateNames {
				if selected == t.GetName() {
					ts = append(ts, t)
					continue WorkspaceLoop
				}
			}
		}
		tmpls = ts
	}
	return tmpls
}

func (o *GetTemplatesOption) Output(tmpls []*dashv1alpha1.Template) {
	data := [][]string{}

	for _, v := range tmpls {
		rawRequiredUseraddons := strings.Join(v.RequiredUseraddons, ",")

		rawUserroles := strings.Join(v.Userroles, ",")

		var isDefaultUserAddon bool
		if v.IsDefaultUserAddon != nil {
			isDefaultUserAddon = *v.IsDefaultUserAddon
		}

		data = append(data, []string{v.GetName(), strconv.FormatBool(isDefaultUserAddon), rawUserroles, rawRequiredUseraddons})

	}

	cli.OutputTable(o.Out,
		[]string{"NAME", "DEFAULT", "REQUIRED_USERROLES", "REQUIRED_USERADDONS"},
		data)
}

func (o *GetTemplatesOption) OutputDetail(tmpl *dashv1alpha1.Template) {
	b, err := yaml.Marshal(tmpl)
	if err != nil {
		fmt.Printf("failed to marshal template: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

func (o *GetTemplatesOption) ListWorkspaceTemplatesByKubeClient(ctx context.Context) ([]*dashv1alpha1.Template, error) {
	c := o.KosmoClient
	tmpls, err := c.ListWorkspaceTemplates(ctx)
	if err != nil {
		return nil, err
	}
	return apiconv.C2D_Templates(tmpls), nil
}
