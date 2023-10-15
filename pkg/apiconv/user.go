package apiconv

import (
	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

func C2D_Users(users []cosmov1alpha1.User) []*dashv1alpha1.User {
	ts := make([]*dashv1alpha1.User, 0, len(users))
	for _, v := range users {
		ts = append(ts, C2D_User(v))
	}
	return ts
}

func C2D_User(user cosmov1alpha1.User) *dashv1alpha1.User {
	return &dashv1alpha1.User{
		Name:        user.Name,
		DisplayName: user.Spec.DisplayName,
		Roles:       C2S_UserRole(user.Spec.Roles),
		AuthType:    user.Spec.AuthType.String(),
		Addons:      C2D_UserAddons(user.Spec.Addons),
		Status:      string(user.Status.Phase),
	}
}

func C2S_UserRole(apiRoles []cosmov1alpha1.UserRole) []string {
	roles := make([]string, 0, len(apiRoles))
	for _, v := range apiRoles {
		roles = append(roles, v.Name)
	}
	return roles
}

func S2C_UserRoles(roles []string) []cosmov1alpha1.UserRole {
	apiRoles := make([]cosmov1alpha1.UserRole, 0, len(roles))
	for _, v := range roles {
		apiRoles = append(apiRoles, cosmov1alpha1.UserRole{Name: v})
	}
	return apiRoles
}

func D2C_UserAddons(addons []*dashv1alpha1.UserAddon) []cosmov1alpha1.UserAddon {
	a := make([]cosmov1alpha1.UserAddon, len(addons))
	for i, v := range addons {
		addon := cosmov1alpha1.UserAddon{
			Template: cosmov1alpha1.UserAddonTemplateRef{
				Name:          v.Template,
				ClusterScoped: v.ClusterScoped,
			},
			Vars: v.Vars,
		}
		a[i] = addon
	}
	return a
}

func C2D_UserAddons(addons []cosmov1alpha1.UserAddon) []*dashv1alpha1.UserAddon {
	da := make([]*dashv1alpha1.UserAddon, len(addons))
	for i, v := range addons {
		da[i] = &dashv1alpha1.UserAddon{
			Template:      v.Template.Name,
			ClusterScoped: v.Template.ClusterScoped,
			Vars:          v.Vars,
		}
	}
	return da
}
