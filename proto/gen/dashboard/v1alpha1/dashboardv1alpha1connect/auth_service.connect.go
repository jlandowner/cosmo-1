//
//Cosmo Dashboard API
//Manipulate cosmo dashboard resource API

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: dashboard/v1alpha1/auth_service.proto

package dashboardv1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// AuthServiceName is the fully-qualified name of the AuthService service.
	AuthServiceName = "dashboard.v1alpha1.AuthService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AuthServiceLoginProcedure is the fully-qualified name of the AuthService's Login RPC.
	AuthServiceLoginProcedure = "/dashboard.v1alpha1.AuthService/Login"
	// AuthServiceLogoutProcedure is the fully-qualified name of the AuthService's Logout RPC.
	AuthServiceLogoutProcedure = "/dashboard.v1alpha1.AuthService/Logout"
	// AuthServiceVerifyProcedure is the fully-qualified name of the AuthService's Verify RPC.
	AuthServiceVerifyProcedure = "/dashboard.v1alpha1.AuthService/Verify"
	// AuthServiceServiceAccountLoginProcedure is the fully-qualified name of the AuthService's
	// ServiceAccountLogin RPC.
	AuthServiceServiceAccountLoginProcedure = "/dashboard.v1alpha1.AuthService/ServiceAccountLogin"
)

// AuthServiceClient is a client for the dashboard.v1alpha1.AuthService service.
type AuthServiceClient interface {
	// ID and password to login
	Login(context.Context, *connect_go.Request[v1alpha1.LoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error)
	// Delete session to logout
	Logout(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[emptypb.Empty], error)
	// Verify authorization
	Verify(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[v1alpha1.VerifyResponse], error)
	// Kubernetes ServiceAccount to login
	ServiceAccountLogin(context.Context, *connect_go.Request[v1alpha1.ServiceAccountLoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error)
}

// NewAuthServiceClient constructs a client for the dashboard.v1alpha1.AuthService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAuthServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AuthServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &authServiceClient{
		login: connect_go.NewClient[v1alpha1.LoginRequest, v1alpha1.LoginResponse](
			httpClient,
			baseURL+AuthServiceLoginProcedure,
			opts...,
		),
		logout: connect_go.NewClient[emptypb.Empty, emptypb.Empty](
			httpClient,
			baseURL+AuthServiceLogoutProcedure,
			opts...,
		),
		verify: connect_go.NewClient[emptypb.Empty, v1alpha1.VerifyResponse](
			httpClient,
			baseURL+AuthServiceVerifyProcedure,
			opts...,
		),
		serviceAccountLogin: connect_go.NewClient[v1alpha1.ServiceAccountLoginRequest, v1alpha1.LoginResponse](
			httpClient,
			baseURL+AuthServiceServiceAccountLoginProcedure,
			opts...,
		),
	}
}

// authServiceClient implements AuthServiceClient.
type authServiceClient struct {
	login               *connect_go.Client[v1alpha1.LoginRequest, v1alpha1.LoginResponse]
	logout              *connect_go.Client[emptypb.Empty, emptypb.Empty]
	verify              *connect_go.Client[emptypb.Empty, v1alpha1.VerifyResponse]
	serviceAccountLogin *connect_go.Client[v1alpha1.ServiceAccountLoginRequest, v1alpha1.LoginResponse]
}

// Login calls dashboard.v1alpha1.AuthService.Login.
func (c *authServiceClient) Login(ctx context.Context, req *connect_go.Request[v1alpha1.LoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error) {
	return c.login.CallUnary(ctx, req)
}

// Logout calls dashboard.v1alpha1.AuthService.Logout.
func (c *authServiceClient) Logout(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[emptypb.Empty], error) {
	return c.logout.CallUnary(ctx, req)
}

// Verify calls dashboard.v1alpha1.AuthService.Verify.
func (c *authServiceClient) Verify(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[v1alpha1.VerifyResponse], error) {
	return c.verify.CallUnary(ctx, req)
}

// ServiceAccountLogin calls dashboard.v1alpha1.AuthService.ServiceAccountLogin.
func (c *authServiceClient) ServiceAccountLogin(ctx context.Context, req *connect_go.Request[v1alpha1.ServiceAccountLoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error) {
	return c.serviceAccountLogin.CallUnary(ctx, req)
}

// AuthServiceHandler is an implementation of the dashboard.v1alpha1.AuthService service.
type AuthServiceHandler interface {
	// ID and password to login
	Login(context.Context, *connect_go.Request[v1alpha1.LoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error)
	// Delete session to logout
	Logout(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[emptypb.Empty], error)
	// Verify authorization
	Verify(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[v1alpha1.VerifyResponse], error)
	// Kubernetes ServiceAccount to login
	ServiceAccountLogin(context.Context, *connect_go.Request[v1alpha1.ServiceAccountLoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error)
}

// NewAuthServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAuthServiceHandler(svc AuthServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(AuthServiceLoginProcedure, connect_go.NewUnaryHandler(
		AuthServiceLoginProcedure,
		svc.Login,
		opts...,
	))
	mux.Handle(AuthServiceLogoutProcedure, connect_go.NewUnaryHandler(
		AuthServiceLogoutProcedure,
		svc.Logout,
		opts...,
	))
	mux.Handle(AuthServiceVerifyProcedure, connect_go.NewUnaryHandler(
		AuthServiceVerifyProcedure,
		svc.Verify,
		opts...,
	))
	mux.Handle(AuthServiceServiceAccountLoginProcedure, connect_go.NewUnaryHandler(
		AuthServiceServiceAccountLoginProcedure,
		svc.ServiceAccountLogin,
		opts...,
	))
	return "/dashboard.v1alpha1.AuthService/", mux
}

// UnimplementedAuthServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAuthServiceHandler struct{}

func (UnimplementedAuthServiceHandler) Login(context.Context, *connect_go.Request[v1alpha1.LoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.AuthService.Login is not implemented"))
}

func (UnimplementedAuthServiceHandler) Logout(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.AuthService.Logout is not implemented"))
}

func (UnimplementedAuthServiceHandler) Verify(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[v1alpha1.VerifyResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.AuthService.Verify is not implemented"))
}

func (UnimplementedAuthServiceHandler) ServiceAccountLogin(context.Context, *connect_go.Request[v1alpha1.ServiceAccountLoginRequest]) (*connect_go.Response[v1alpha1.LoginResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.AuthService.ServiceAccountLogin is not implemented"))
}
