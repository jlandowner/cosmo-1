package apiconv

import (
	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

func ConvertDashv1alpha1UserAddonToUserAddon(addons []*dashv1alpha1.UserAddons) []cosmov1alpha1.UserAddon {
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
