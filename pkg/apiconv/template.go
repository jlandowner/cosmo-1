package apiconv

import (
	"strconv"
	"strings"

	"github.com/cosmo-workspace/cosmo/pkg/kubeutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/utils/pointer"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

func C2D_Templates(tmpls []cosmov1alpha1.TemplateObject) []*dashv1alpha1.Template {
	dTmpls := make([]*dashv1alpha1.Template, len(tmpls))
	for i, v := range tmpls {
		dTmpls[i] = C2D_Template(v)
	}
	return dTmpls
}

func C2D_Template(tmpl cosmov1alpha1.TemplateObject) *dashv1alpha1.Template {
	requiredVars := make([]*dashv1alpha1.TemplateRequiredVars, len(tmpl.GetSpec().RequiredVars))
	for i, v := range tmpl.GetSpec().RequiredVars {
		requiredVars[i] = &dashv1alpha1.TemplateRequiredVars{
			VarName:      v.Var,
			DefaultValue: v.Default,
		}
	}

	return &dashv1alpha1.Template{
		Name:           tmpl.GetName(),
		Description:    tmpl.GetSpec().Description,
		RequiredVars:   requiredVars,
		IsClusterScope: tmpl.GetScope() == meta.RESTScopeRoot,
		IsDefaultUserAddon: func() *bool {
			if ann := tmpl.GetAnnotations(); ann != nil {
				if b, ok := ann[cosmov1alpha1.UserAddonTemplateAnnKeyDefaultUserAddon]; ok {
					if defaultAddon, err := strconv.ParseBool(b); err == nil && defaultAddon {
						return pointer.Bool(true)
					}
				}
			}
			return nil
		}(),
		RequiredUseraddons: func() []string {
			requiredAddons := kubeutil.GetAnnotation(tmpl, cosmov1alpha1.TemplateAnnKeyRequiredAddons)
			if requiredAddons != "" {
				return strings.Split(requiredAddons, ",")
			}
			return nil
		}(),
		Userroles: func() []string {
			requiredAddons := kubeutil.GetAnnotation(tmpl, cosmov1alpha1.TemplateAnnKeyUserRoles)
			if requiredAddons != "" {
				return strings.Split(requiredAddons, ",")
			}
			return nil
		}(),
	}
}
