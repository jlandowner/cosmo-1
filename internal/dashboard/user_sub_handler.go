package dashboard

import (
	"context"

	connect_go "github.com/bufbuild/connect-go"

	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/kosmo"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

func (s *Server) UpdateUserAddons(ctx context.Context, req *connect_go.Request[dashv1alpha1.UpdateUserAddonsRequest]) (*connect_go.Response[dashv1alpha1.UpdateUserAddonsResponse], error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "req", req)

	// targetUser, err := s.Klient.GetUser(ctx, req.Msg.UserName)
	// if err != nil {
	// 	return nil, ErrResponse(log, err)
	// }

	caller := callerFromContext(ctx)
	if caller == nil {
		return nil, kosmo.NewInternalServerError("unable get caller", nil)
	}

	// changingAddons := diffUserAddons(targetUser.Spec.Addons, convertDashv1alpha1UserAddonToUserAddon(req.Msg.Addons))
	// for _, addon := range changingAddons {
	// 	if _, ok := find(availableTemplates, addon,
	// 		func(e cosmov1alpha1.TemplateObject, v cosmov1alpha1.UserAddon) bool {
	// 			return e.GetName() == v.Template.Name
	// 		}); !ok {
	// 		return nil, ErrResponse(log, kosmo.NewForbiddenError(fmt.Sprintf("denied to attach '%s' group", addon.Template.Name), nil))
	// 	}
	// }

	addons := convertDashv1alpha1UserAddonToUserAddon(req.Msg.Addons)
	user, err := s.Klient.UpdateUser(ctx, req.Msg.UserName, kosmo.UpdateUserOpts{
		UserAddons: addons,
	})
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	res := &dashv1alpha1.UpdateUserAddonsResponse{
		Message: "Successfully updated",
		User:    convertUserToDashv1alpha1User(*user),
	}
	log.Info(res.Message, "username", req.Msg.UserName)
	return connect_go.NewResponse(res), nil
}

func (s *Server) UpdateUserDisplayName(ctx context.Context, req *connect_go.Request[dashv1alpha1.UpdateUserDisplayNameRequest]) (*connect_go.Response[dashv1alpha1.UpdateUserDisplayNameResponse], error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "req", req)

	if err := userAuthentication(ctx, req.Msg.UserName); err != nil {
		return nil, ErrResponse(log, err)
	}

	user, err := s.Klient.UpdateUser(ctx, req.Msg.UserName, kosmo.UpdateUserOpts{DisplayName: &req.Msg.DisplayName})
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	res := &dashv1alpha1.UpdateUserDisplayNameResponse{
		Message: "Successfully updated",
		User:    convertUserToDashv1alpha1User(*user),
	}
	log.Info(res.Message, "username", req.Msg.UserName)
	return connect_go.NewResponse(res), nil
}

func diff(slice1 []string, slice2 []string) []string {
	var diff []string
	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices
		slice1, slice2 = slice2, slice1
	}
	return diff
}

func (s *Server) UpdateUserRole(ctx context.Context, req *connect_go.Request[dashv1alpha1.UpdateUserRoleRequest]) (*connect_go.Response[dashv1alpha1.UpdateUserRoleResponse], error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "req", req)

	currentUser, err := s.Klient.GetUser(ctx, req.Msg.UserName)
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	// group-admin can attach or detach only group-roles
	changingRoles := diff(convertUserRolesToStringSlice(currentUser.Spec.Roles), req.Msg.Roles)
	if err := adminAuthentication(ctx, validateCallerHasAdminForAllRoles(changingRoles)); err != nil {
		return nil, ErrResponse(log, err)
	}

	roles := convertStringSliceToUserRoles(req.Msg.Roles)
	user, err := s.Klient.UpdateUser(ctx, req.Msg.UserName, kosmo.UpdateUserOpts{UserRoles: roles})
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	res := &dashv1alpha1.UpdateUserRoleResponse{
		Message: "Successfully updated",
		User:    convertUserToDashv1alpha1User(*user),
	}
	log.Info(res.Message, "username", req.Msg.UserName)
	return connect_go.NewResponse(res), nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, req *connect_go.Request[dashv1alpha1.UpdateUserPasswordRequest]) (*connect_go.Response[dashv1alpha1.UpdateUserPasswordResponse], error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "username", req.Msg.UserName)

	if err := userAuthentication(ctx, req.Msg.UserName); err != nil {
		return nil, ErrResponse(log, err)
	}

	// check current password is valid
	verified, _, err := s.Klient.VerifyPassword(ctx, req.Msg.UserName, []byte(req.Msg.CurrentPassword))
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	if !verified {
		return nil, ErrResponse(log, kosmo.NewForbiddenError("incorrect user or password", nil))
	}

	// Upsert password
	if err := s.Klient.RegisterPassword(ctx, req.Msg.UserName, []byte(req.Msg.NewPassword)); err != nil {
		return nil, ErrResponse(log, err)
	}

	res := &dashv1alpha1.UpdateUserPasswordResponse{
		Message: "Successfully updated",
	}
	log.Info(res.Message, "username", req.Msg.UserName)
	return connect_go.NewResponse(res), nil
}
