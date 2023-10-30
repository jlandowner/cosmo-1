package cli

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
	"github.com/spf13/cobra"
)

type CosmoDashClient struct {
	BaseURL *url.URL
	Cookie  string

	AuthServiceClient      dashboardv1alpha1connect.AuthServiceClient
	UserServiceClient      dashboardv1alpha1connect.UserServiceClient
	WorkspaceServiceClient dashboardv1alpha1connect.WorkspaceServiceClient
	TemplateServiceClient  dashboardv1alpha1connect.TemplateServiceClient
	WebAuthnServiceClient  dashboardv1alpha1connect.WebAuthnServiceClient
}

func NewCosmoDashClient(httpClient connect.HTTPClient, baseURL *url.URL) *CosmoDashClient {
	clientOptions := connect.WithClientOptions(
		connect.WithGRPCWeb(),
		connect.WithSendGzip(),
	)
	return &CosmoDashClient{
		BaseURL:                baseURL,
		AuthServiceClient:      dashboardv1alpha1connect.NewAuthServiceClient(httpClient, baseURL.String(), clientOptions),
		UserServiceClient:      dashboardv1alpha1connect.NewUserServiceClient(httpClient, baseURL.String(), clientOptions),
		WorkspaceServiceClient: dashboardv1alpha1connect.NewWorkspaceServiceClient(httpClient, baseURL.String(), clientOptions),
		TemplateServiceClient:  dashboardv1alpha1connect.NewTemplateServiceClient(httpClient, baseURL.String(), clientOptions),
		WebAuthnServiceClient:  dashboardv1alpha1connect.NewWebAuthnServiceClient(httpClient, baseURL.String(), clientOptions),
	}
}

func (c *CosmoDashClient) GetSession(ctx context.Context, userName string, password string) (string, error) {
	res, err := c.AuthServiceClient.Login(ctx, connect.NewRequest(&dashv1alpha1.LoginRequest{UserName: userName, Password: password}))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(res.Header().Get("Set-Cookie"))), nil
}

func NewRequestWithToken[T any](message *T, cfg *Config) *connect.Request[T] {
	req := connect.NewRequest(message)
	if cfg != nil {
		s, err := base64.StdEncoding.DecodeString(cfg.Token)
		if err != nil {
			panic(err)
		}
		req.Header().Add("Cookie", string(s))
	}
	return req
}

type RunCommand interface {
	RunE(*cobra.Command, []string) error
	Logger() *clog.Logger
}

func ConnectErrorHandler(rcmd RunCommand) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := rcmd.RunE(cmd, args)
		var connectErr *connect.Error
		if errors.As(err, &connectErr) {
			rcmd.Logger().Debug().Info("connectErr", "code", connectErr.Code(), "message", connectErr.Message())
			if connectErr.Code() == connect.CodeUnknown {
				if strings.Index(connectErr.Message(), fmt.Sprintf("%d", http.StatusFound)) > 0 {
					return fmt.Errorf("session has been expired: please login again")
				}
				return fmt.Errorf("%w: session might have been expired", err)
			}
		}
		return err
	}
}
