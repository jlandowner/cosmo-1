package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/cli"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type GetOption struct {
	*cli.RootOptions

	UserNames []string
	Filter    []string

	filters []cli.Filter
}

func GetCmd(cmd *cobra.Command, opt *cli.RootOptions) *cobra.Command {
	o := &GetOption{RootOptions: opt}
	cmd.RunE = cli.ConnectErrorHandler(o)
	cmd.Flags().StringSliceVar(&o.Filter, "filter", nil, "filter option. available columns are ['NAME', 'ROLE', 'ADDON', 'AUTHTYPE', 'PHASE']. available operators are ['==', '!=']. value format is filepath. e.g. '--filter ROLE==*-dev --filter ROLE!=team-a'")
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
		o.UserNames = args
	}
	if len(o.Filter) > 0 {
		o.filters = cli.ParseFilters(o.Filter)
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

	var users []*dashv1alpha1.User
	var err error
	if o.UseKubeAPI {
		users, err = o.ListUsersByKubeClient(ctx)
		if err != nil {
			return err
		}
	} else {
		users, err = o.ListUsersWithDashClient(ctx)
		if err != nil {
			if connect_go.CodeOf(err) == connect_go.CodePermissionDenied {

				if len(o.UserNames) == 0 {
					cmdutil.PrintfColorErr(o.ErrOut, "WARNING: Without Admin roles, you can get only login user\n")
				} else {
					for _, v := range o.UserNames {
						if v != o.CliConfig.User {
							return fmt.Errorf("permission denied: failed to get user: %s", v)
						}
					}
				}
				me, err := o.GetUserWithDashClient(ctx, o.CliConfig.User)
				if err != nil {
					return err
				}
				users = []*dashv1alpha1.User{me}
			} else {
				return err
			}
		}
	}
	o.Logr.Debug().Info("Users", "users", users)

	users = o.ApplyFilters(users)

	o.Output(users)

	return nil

}

func (o *GetOption) ListUsersWithDashClient(ctx context.Context) ([]*dashv1alpha1.User, error) {
	c := o.CosmoDashClient
	res, err := c.UserServiceClient.GetUsers(ctx, cli.NewRequestWithToken(&emptypb.Empty{}, o.CliConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("UserServiceClient.GetUsers", "res", res)
	return res.Msg.Items, nil
}

func (o *GetOption) GetUserWithDashClient(ctx context.Context, userName string) (*dashv1alpha1.User, error) {
	c := o.CosmoDashClient
	res, err := c.UserServiceClient.GetUser(ctx, cli.NewRequestWithToken(&dashv1alpha1.GetUserRequest{UserName: userName}, o.CliConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to connect dashboard server: %w", err)
	}
	o.Logr.DebugAll().Info("UserServiceClient.GetUser", "res", res)
	return res.Msg.User, nil
}

func (o *GetOption) ApplyFilters(users []*dashv1alpha1.User) []*dashv1alpha1.User {
	for _, f := range o.filters {
		o.Logr.Debug().Info("applying filter", "key", f.Key, "value", f.Value, "op", f.Operator)

		switch strings.ToUpper(f.Key) {
		case "NAME":
			users = cli.DoFilter(users, func(u *dashv1alpha1.User) []string {
				return []string{u.Name}
			}, f)
		case "ROLE", "ROLES":
			users = cli.DoFilter(users, func(u *dashv1alpha1.User) []string {
				arr := make([]string, 0, len(u.Roles))
				arr = append(arr, u.Roles...)
				return arr
			}, f)
		case "ADDON", "ADDONS":
			users = cli.DoFilter(users, func(u *dashv1alpha1.User) []string {
				arr := make([]string, 0, len(u.Addons))
				for _, a := range u.Addons {
					arr = append(arr, a.Template)
				}
				return arr
			}, f)
		case "AUTHTYPE":
			users = cli.DoFilter(users, func(u *dashv1alpha1.User) []string {
				return []string{u.AuthType}
			}, f)
		case "PHASE":
			users = cli.DoFilter(users, func(u *dashv1alpha1.User) []string {
				return []string{u.Status}
			}, f)
		default:
			o.Logr.Info("WARNING: unknown filter key", "key", f.Key)
		}
	}

	if len(o.UserNames) > 0 {
		ts := make([]*dashv1alpha1.User, 0, len(o.UserNames))
	UserLoop:
		// Or loop
		for _, t := range users {
			for _, selected := range o.UserNames {
				if selected == t.GetName() {
					ts = append(ts, t)
					continue UserLoop
				}
			}
		}
		users = ts
	}
	return users
}

func (o *GetOption) Output(users []*dashv1alpha1.User) {
	data := [][]string{}

	for _, v := range users {
		role := make([]string, 0, len(v.Roles))
		role = append(role, v.Roles...)

		addons := make([]string, 0, len(v.Addons))
		for _, v := range v.Addons {
			addons = append(addons, v.Template)
		}
		data = append(data, []string{v.Name, strings.Join(role, ","), v.AuthType, cosmov1alpha1.UserNamespace(v.Name), v.Status, strings.Join(addons, ",")})
	}

	cli.OutputTable(o.Out,
		[]string{"NAME", "ROLES", "AUTHTYPE", "NAMESPACE", "PHASE", "ADDONS"},
		data)
}

func (o *GetOption) ListUsersByKubeClient(ctx context.Context) ([]*dashv1alpha1.User, error) {
	c := o.KosmoClient
	users, err := c.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return apiconv.C2D_Users(users), nil
}
