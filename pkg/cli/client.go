package cli

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/bufbuild/connect-go"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
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
