package dashboard

import (
	"context"
	"errors"
	"net/http"
	"time"

	dashv1alpha1 "github.com/cosmo-workspace/cosmo/api/openapi/dashboard/v1alpha1"
	wsv1alpha1 "github.com/cosmo-workspace/cosmo/api/workspace/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/auth/session"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/gorilla/mux"
)

func (s *Server) useAuthMiddleWare(router *mux.Router, routes dashv1alpha1.Routes) {
	for _, rt := range routes {
		router.Get(rt.Name).Handler(s.contextMiddleware(router.Get(rt.Name).GetHandler()))
	}
	router.Get("Verify").Handler(s.authorizationMiddleware(router.Get("Verify").GetHandler()))
}

func (s *Server) Verify(ctx context.Context) (dashv1alpha1.ImplResponse, error) {

	caller := callerFromContext(ctx)
	if caller == nil {
		return dashv1alpha1.Response(http.StatusUnauthorized, nil), nil
	}
	deadline := deadlineFromContext(ctx)
	if deadline.Before(time.Now()) {
		return dashv1alpha1.Response(http.StatusUnauthorized, nil), nil
	}

	res := &dashv1alpha1.VerifyResponse{
		Id:       caller.ID,
		ExpireAt: deadline,
	}
	return dashv1alpha1.Response(http.StatusOK, res), nil
}

func (s *Server) Logout(ctx context.Context) (dashv1alpha1.ImplResponse, error) {

	w := responseWriterFromContext(ctx)
	r := requestFromContext(ctx)

	_, _, err := s.authorizeWithSession(r)
	if err != nil {
		if errors.Is(err, ErrNotAuthorized) {
			return dashv1alpha1.Response(http.StatusUnauthorized, nil), nil
		} else {
			return dashv1alpha1.Response(http.StatusInternalServerError, nil), nil
		}
	}

	cookie := s.sessionCookieKey()
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	return dashv1alpha1.Response(http.StatusOK, nil), nil
}

func (s *Server) Login(ctx context.Context, req dashv1alpha1.LoginRequest) (dashv1alpha1.ImplResponse, error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "req", req)

	w := responseWriterFromContext(ctx)
	r := requestFromContext(ctx)

	// Check ID
	user, err := s.Klient.GetUser(ctx, req.Id)
	if err != nil {
		log.Info(err.Error(), "userid", req.Id)
		return dashv1alpha1.Response(http.StatusForbidden, nil), nil
	}
	// Check password
	authrizer, ok := s.Authorizers[user.AuthType]
	if !ok {
		log.Info("authrizer not found", "userid", req.Id, "authType", user.AuthType)
		return dashv1alpha1.Response(http.StatusServiceUnavailable, nil), nil
	}
	verified, err := authrizer.Authorize(ctx, req)
	if err != nil {
		log.Error(err, "authorize failed", "userid", req.Id)
		return dashv1alpha1.Response(http.StatusForbidden, nil), nil
	}
	if !verified {
		log.Info("login failed: password invalid", "userid", req.Id)
		return dashv1alpha1.Response(http.StatusForbidden, nil), nil
	}
	var isDefault bool
	if wsv1alpha1.UserAuthType(user.AuthType) == wsv1alpha1.UserAuthTypeKosmoSecert {
		isDefault, err = s.Klient.IsDefaultPassword(ctx, req.Id)
		if err != nil {
			log.Error(err, "failed to check is default password", "userid", req.Id)
			return dashv1alpha1.Response(http.StatusInternalServerError, nil), nil
		}
	}

	// Create session
	now := time.Now()
	expireAt := now.Add(time.Duration(s.MaxAgeSeconds) * time.Second)

	ses, _ := s.sessionStore.New(r, s.SessionName)
	sesInfo := session.Info{
		UserID:   req.Id,
		Deadline: expireAt.Unix(),
	}
	log.DebugAll().Info("save session", "userID", sesInfo.UserID, "deadline", sesInfo.Deadline)
	ses = session.Set(ses, sesInfo)

	err = s.sessionStore.Save(r, w, ses)
	if err != nil {
		log.Error(err, "failed to save session")
		return dashv1alpha1.Response(http.StatusInternalServerError, nil), nil
	}

	res := &dashv1alpha1.LoginResponse{
		Id:                    req.Id,
		ExpireAt:              expireAt,
		RequirePasswordUpdate: isDefault,
	}

	return dashv1alpha1.Response(http.StatusOK, res), nil
}