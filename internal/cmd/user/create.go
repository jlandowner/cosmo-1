package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

type CreateOption struct {
	*cmdutil.CliOptions

	UserName      string
	DisplayName   string
	Role          string
	Admin         bool
	AuthType      string
	Addons        []string
	ClusterAddons []string
	Quiet         bool

	userAddons []*dashv1alpha1.UserAddons
}

func CreateCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &CreateOption{CliOptions: cliOpt}
	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)
	cmd.Flags().StringVar(&o.DisplayName, "name", "", "user display name (default: same as USER_NAME)")
	cmd.Flags().StringVar(&o.Role, "role", "", "user role")
	cmd.Flags().BoolVar(&o.Admin, "admin", false, "user admin role")
	cmd.Flags().StringVar(&o.AuthType, "auth", cosmov1alpha1.UserAuthTypePasswordSecert.String(), "user auth type")
	cmd.Flags().StringArrayVar(&o.Addons, "addon", nil, "user addons by Template, which created in UserNamespace\nformat is '--addon TEMPLATE_NAME1,KEY:VAL,KEY:VAL --addon TEMPLATE_NAME2,KEY:VAL ...' ")
	cmd.Flags().StringArrayVar(&o.ClusterAddons, "cluster-addon", nil, "user addons by ClusterTemplate\nformat is '--cluster-addon TEMPLATE_NAME1,KEY:VAL,KEY:VAL --cluster-addon TEMPLATE_NAME2,KEY:VAL ...' ")
	cmd.Flags().BoolVarP(&o.Quiet, "quiet", "q", false, "only output default password")

	// support using kube client
	o.AddKubeClientFlags(cmd)

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
	if o.Role != "" {
		if o.Admin {
			return errors.New("--role and --admin is not used at the same time")
		}
		if !cosmov1alpha1.UserRole(o.Role).IsValid() {
			return fmt.Errorf("role %s is invalid", o.Role)
		}
	}
	if !cosmov1alpha1.UserAuthType(o.AuthType).IsValid() {
		return fmt.Errorf("auth %s is invalid. acceptable values: %v", o.AuthType, cosmov1alpha1.UserAuthTypeList)
	}
	return nil
}

func (o *CreateOption) Complete(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Complete(cmd, args); err != nil {
		return err
	}

	o.UserName = args[0]

	if o.Admin {
		o.Role = cosmov1alpha1.UserAdminRole.String()
	}

	o.userAddons = make([]*dashv1alpha1.UserAddons, 0, len(o.Addons)+len(o.ClusterAddons))
	if len(o.Addons) > 0 {
		userAddons, err := parseUserAddonOptions(o.Addons, false)
		if err != nil {
			return err
		}
		o.userAddons = append(o.userAddons, userAddons...)
	}
	if len(o.ClusterAddons) > 0 {
		userAddons, err := parseUserAddonOptions(o.ClusterAddons, true)
		if err != nil {
			return err
		}
		o.userAddons = append(o.userAddons, userAddons...)
	}

	return nil
}

func parseUserAddonOptions(rawAddonOptionArray []string, isClusterScope bool) ([]*dashv1alpha1.UserAddons, error) {
	// format
	//   TEMPLATE_NAME
	//   TEMPLATE_NAME,KEY1:XXX,KEY2:YYY ZZZ,KEY3:
	r1 := regexp.MustCompile(`^[^: ,]+(,([^: ,]+):([^,]*))*$`)
	r2 := regexp.MustCompile(`^([^: ,]+):([^,]*)$`)

	userAddons := make([]*dashv1alpha1.UserAddons, 0, len(rawAddonOptionArray))

	for _, addonParm := range rawAddonOptionArray {
		if !r1.MatchString(addonParm) {
			return nil, fmt.Errorf("invalid addon vars format: %s", addonParm)
		}

		addonSplits := strings.Split(addonParm, ",")

		userAddon := &dashv1alpha1.UserAddons{
			Template:      addonSplits[0],
			ClusterScoped: isClusterScope,
			Vars:          make(map[string]string, len(addonSplits)-1),
		}

		for _, k_v := range addonSplits[1:] {
			kv := r2.FindStringSubmatch(k_v)
			userAddon.Vars[kv[1]] = kv[2]
		}
		userAddons = append(userAddons, userAddon)
	}
	return userAddons, nil
}

func (o *CreateOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("create_user")
	ctx := clog.IntoContext(o.Ctx, log)

	var defaultPassword string
	if o.UseKubeClient {
		addons := apiconv.ConvertDashv1alpha1UserAddonToUserAddon(o.userAddons)
		if _, err := o.KubeClient.CreateUser(ctx, o.UserName, o.DisplayName, o.Role, o.AuthType, addons); err != nil {
			return err
		}

		dp, err := o.KubeClient.GetDefaultPasswordAwait(ctx, o.UserName)
		if err != nil {
			return err
		}
		defaultPassword = *dp

	} else {
		c := dashboardv1alpha1connect.NewUserServiceClient(o.Client, o.ServerEndpoint, connect.WithGRPC())

		res, err := c.CreateUser(ctx, cmdutil.NewConnectRequestWithAuth(o.CliConfig,
			&dashv1alpha1.CreateUserRequest{
				UserName:    o.UserName,
				DisplayName: o.DisplayName,
				Role:        o.Role,
				AuthType:    o.AuthType,
				Addons:      o.userAddons,
			}))
		if err != nil {
			return err
		}
		log.Debug().Info("response: %v", res)

		defaultPassword = res.Msg.User.DefaultPassword
	}

	if o.Quiet {
		fmt.Fprintln(o.Out, defaultPassword)
	} else {
		cmdutil.PrintfColorInfo(o.Out, "Successfully created user %s\n", o.UserName)
		fmt.Fprintln(o.Out, "Default password:", defaultPassword)
	}

	return nil
}
