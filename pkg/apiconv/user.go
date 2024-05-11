package apiconv

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	eventsv1 "k8s.io/api/events/v1"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

type UserConvertOptions func(c *cosmov1alpha1.User, d *dashv1alpha1.User)

func WithUserRaw(withRaw *bool) func(c *cosmov1alpha1.User, d *dashv1alpha1.User) {
	return func(c *cosmov1alpha1.User, d *dashv1alpha1.User) {
		if withRaw != nil && *withRaw {
			d.Raw = ToYAML(c)
		}
	}
}

func C2D_Users(users []cosmov1alpha1.User, opts ...UserConvertOptions) []*dashv1alpha1.User {
	ts := make([]*dashv1alpha1.User, len(users))
	for i, v := range users {
		ts[i] = C2D_User(v, opts...)
	}
	return ts
}

func C2D_User(user cosmov1alpha1.User, opts ...UserConvertOptions) *dashv1alpha1.User {
	d := &dashv1alpha1.User{
		Name:        user.Name,
		DisplayName: user.Spec.DisplayName,
		Roles:       C2S_UserRole(user.Spec.Roles),
		AuthType:    user.Spec.AuthType.String(),
		Addons:      C2D_UserAddons(user.Spec.Addons),
		Status:      string(user.Status.Phase),
	}
	for _, opt := range opts {
		opt(&user, d)
	}
	return d
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

func K2D_Events(events []eventsv1.Event) []*dashv1alpha1.Event {
	es := make([]*dashv1alpha1.Event, len(events))
	for i, v := range events {
		var eventTime *timestamppb.Timestamp
		if v.EventTime.Year() != 1 {
			eventTime = timestamppb.New(v.EventTime.Time)
		} else {
			eventTime = timestamppb.New(v.DeprecatedLastTimestamp.Time)
		}

		var count int32
		if v.Series != nil {
			count = v.Series.Count
		} else {
			count = v.DeprecatedCount
		}

		var lastTime *timestamppb.Timestamp
		if v.Series != nil {
			lastTime = timestamppb.New(v.Series.LastObservedTime.Time)
		} else {
			lastTime = timestamppb.New(v.DeprecatedLastTimestamp.Time)
		}

		e := &dashv1alpha1.Event{
			EventTime: eventTime,
			Reason:    v.Reason,
			Note:      v.Note,
			Type:      v.Type,
			Regarding: &dashv1alpha1.ObjectReference{
				ApiVersion: v.Regarding.APIVersion,
				Kind:       v.Regarding.Kind,
				Name:       v.Regarding.Name,
				Namespace:  v.Regarding.Namespace,
			},
			ReportingController: v.ReportingController,
		}
		if count > 1 {
			e.Series = &dashv1alpha1.EventSeries{
				Count:            count,
				LastObservedTime: lastTime,
			}
		}
		es[i] = e
	}
	return es
}
