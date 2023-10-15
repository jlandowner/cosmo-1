package apiconv

import (
	"k8s.io/utils/ptr"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

func C2D_Workspaces(wss []cosmov1alpha1.Workspace) []*dashv1alpha1.Workspace {
	apiwss := make([]*dashv1alpha1.Workspace, len(wss))
	for i, v := range wss {
		apiwss[i] = C2D_Workspace(v)
	}
	return apiwss
}

func C2D_Workspace(ws cosmov1alpha1.Workspace) *dashv1alpha1.Workspace {
	replicas := ws.Spec.Replicas
	if replicas == nil {
		replicas = ptr.To(int64(1))
	}

	return &dashv1alpha1.Workspace{
		Name:      ws.Name,
		OwnerName: cosmov1alpha1.UserNameByNamespace(ws.Namespace),
		Spec: &dashv1alpha1.WorkspaceSpec{
			Template: ws.Spec.Template.Name,
			Replicas: *replicas,
			Vars:     ws.Spec.Vars,
			Network:  C2D_NetworkRules(ws.Spec.Network, ws.Status.URLs),
		},
		Status: &dashv1alpha1.WorkspaceStatus{
			Phase:   string(ws.Status.Phase),
			MainUrl: ws.Status.URLs[cosmov1alpha1.MainRuleKey(ws.Status.Config)],
		},
	}
}

func C2D_NetworkRules(netRules []cosmov1alpha1.NetworkRule, urlMap map[string]string) []*dashv1alpha1.NetworkRule {
	apirules := make([]*dashv1alpha1.NetworkRule, 0, len(netRules))
	for _, v := range netRules {
		r := C2D_NetworkRule(v)
		r.Url = urlMap[v.UniqueKey()]
		apirules = append(apirules, &r)
	}
	return apirules
}

func C2D_NetworkRule(v cosmov1alpha1.NetworkRule) dashv1alpha1.NetworkRule {
	return dashv1alpha1.NetworkRule{
		PortNumber:       int32(v.PortNumber),
		CustomHostPrefix: v.CustomHostPrefix,
		HttpPath:         v.HTTPPath,
		Public:           v.Public,
	}
}
