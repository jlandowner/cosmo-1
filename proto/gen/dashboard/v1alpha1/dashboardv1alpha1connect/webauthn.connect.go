//
//WebAuthn protobuf

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: dashboard/v1alpha1/webauthn.proto

package dashboardv1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
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
	// WebAuthnServiceName is the fully-qualified name of the WebAuthnService service.
	WebAuthnServiceName = "dashboard.v1alpha1.WebAuthnService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// WebAuthnServiceBeginRegistrationProcedure is the fully-qualified name of the WebAuthnService's
	// BeginRegistration RPC.
	WebAuthnServiceBeginRegistrationProcedure = "/dashboard.v1alpha1.WebAuthnService/BeginRegistration"
	// WebAuthnServiceFinishRegistrationProcedure is the fully-qualified name of the WebAuthnService's
	// FinishRegistration RPC.
	WebAuthnServiceFinishRegistrationProcedure = "/dashboard.v1alpha1.WebAuthnService/FinishRegistration"
	// WebAuthnServiceBeginLoginProcedure is the fully-qualified name of the WebAuthnService's
	// BeginLogin RPC.
	WebAuthnServiceBeginLoginProcedure = "/dashboard.v1alpha1.WebAuthnService/BeginLogin"
	// WebAuthnServiceFinishLoginProcedure is the fully-qualified name of the WebAuthnService's
	// FinishLogin RPC.
	WebAuthnServiceFinishLoginProcedure = "/dashboard.v1alpha1.WebAuthnService/FinishLogin"
	// WebAuthnServiceListCredentialsProcedure is the fully-qualified name of the WebAuthnService's
	// ListCredentials RPC.
	WebAuthnServiceListCredentialsProcedure = "/dashboard.v1alpha1.WebAuthnService/ListCredentials"
	// WebAuthnServiceUpdateCredentialProcedure is the fully-qualified name of the WebAuthnService's
	// UpdateCredential RPC.
	WebAuthnServiceUpdateCredentialProcedure = "/dashboard.v1alpha1.WebAuthnService/UpdateCredential"
	// WebAuthnServiceDeleteCredentialProcedure is the fully-qualified name of the WebAuthnService's
	// DeleteCredential RPC.
	WebAuthnServiceDeleteCredentialProcedure = "/dashboard.v1alpha1.WebAuthnService/DeleteCredential"
)

// WebAuthnServiceClient is a client for the dashboard.v1alpha1.WebAuthnService service.
type WebAuthnServiceClient interface {
	// BeginRegistration returns CredentialCreateOptions to window.navigator.create() which is serialized as JSON string
	// Also `publicKey.user.id“ and `publicKey.challenge` are base64url encoded
	BeginRegistration(context.Context, *connect_go.Request[v1alpha1.BeginRegistrationRequest]) (*connect_go.Response[v1alpha1.BeginRegistrationResponse], error)
	// FinishRegistration check the result of window.navigator.create()
	// `rawId`, `response.clientDataJSON` and `response.attestationObject` in the result must be base64url encoded
	// and all JSON must be serialized as string
	FinishRegistration(context.Context, *connect_go.Request[v1alpha1.FinishRegistrationRequest]) (*connect_go.Response[v1alpha1.FinishRegistrationResponse], error)
	// BeginLogin returns CredentialRequestOptions to window.navigator.get() which is serialized as JSON string
	// Also `publicKey.allowCredentials[*].id` and `publicKey.challenge` are base64url encoded
	BeginLogin(context.Context, *connect_go.Request[v1alpha1.BeginLoginRequest]) (*connect_go.Response[v1alpha1.BeginLoginResponse], error)
	// FinishLogin check the result of window.navigator.get()
	// `rawId`, `response.clientDataJSON`, `response.authenticatorData`, `response.signature`, `response.userHandle`
	// in the result must be base64url encoded and all JSON must be serialized as string
	FinishLogin(context.Context, *connect_go.Request[v1alpha1.FinishLoginRequest]) (*connect_go.Response[v1alpha1.FinishLoginResponse], error)
	// ListCredentials returns registered credential ID list
	ListCredentials(context.Context, *connect_go.Request[v1alpha1.ListCredentialsRequest]) (*connect_go.Response[v1alpha1.ListCredentialsResponse], error)
	// UpdateCredential updates registed credential's human readable infomations
	UpdateCredential(context.Context, *connect_go.Request[v1alpha1.UpdateCredentialRequest]) (*connect_go.Response[v1alpha1.UpdateCredentialResponse], error)
	// DeleteCredential remove registered credential
	DeleteCredential(context.Context, *connect_go.Request[v1alpha1.DeleteCredentialRequest]) (*connect_go.Response[v1alpha1.DeleteCredentialResponse], error)
}

// NewWebAuthnServiceClient constructs a client for the dashboard.v1alpha1.WebAuthnService service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewWebAuthnServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) WebAuthnServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &webAuthnServiceClient{
		beginRegistration: connect_go.NewClient[v1alpha1.BeginRegistrationRequest, v1alpha1.BeginRegistrationResponse](
			httpClient,
			baseURL+WebAuthnServiceBeginRegistrationProcedure,
			opts...,
		),
		finishRegistration: connect_go.NewClient[v1alpha1.FinishRegistrationRequest, v1alpha1.FinishRegistrationResponse](
			httpClient,
			baseURL+WebAuthnServiceFinishRegistrationProcedure,
			opts...,
		),
		beginLogin: connect_go.NewClient[v1alpha1.BeginLoginRequest, v1alpha1.BeginLoginResponse](
			httpClient,
			baseURL+WebAuthnServiceBeginLoginProcedure,
			opts...,
		),
		finishLogin: connect_go.NewClient[v1alpha1.FinishLoginRequest, v1alpha1.FinishLoginResponse](
			httpClient,
			baseURL+WebAuthnServiceFinishLoginProcedure,
			opts...,
		),
		listCredentials: connect_go.NewClient[v1alpha1.ListCredentialsRequest, v1alpha1.ListCredentialsResponse](
			httpClient,
			baseURL+WebAuthnServiceListCredentialsProcedure,
			opts...,
		),
		updateCredential: connect_go.NewClient[v1alpha1.UpdateCredentialRequest, v1alpha1.UpdateCredentialResponse](
			httpClient,
			baseURL+WebAuthnServiceUpdateCredentialProcedure,
			opts...,
		),
		deleteCredential: connect_go.NewClient[v1alpha1.DeleteCredentialRequest, v1alpha1.DeleteCredentialResponse](
			httpClient,
			baseURL+WebAuthnServiceDeleteCredentialProcedure,
			opts...,
		),
	}
}

// webAuthnServiceClient implements WebAuthnServiceClient.
type webAuthnServiceClient struct {
	beginRegistration  *connect_go.Client[v1alpha1.BeginRegistrationRequest, v1alpha1.BeginRegistrationResponse]
	finishRegistration *connect_go.Client[v1alpha1.FinishRegistrationRequest, v1alpha1.FinishRegistrationResponse]
	beginLogin         *connect_go.Client[v1alpha1.BeginLoginRequest, v1alpha1.BeginLoginResponse]
	finishLogin        *connect_go.Client[v1alpha1.FinishLoginRequest, v1alpha1.FinishLoginResponse]
	listCredentials    *connect_go.Client[v1alpha1.ListCredentialsRequest, v1alpha1.ListCredentialsResponse]
	updateCredential   *connect_go.Client[v1alpha1.UpdateCredentialRequest, v1alpha1.UpdateCredentialResponse]
	deleteCredential   *connect_go.Client[v1alpha1.DeleteCredentialRequest, v1alpha1.DeleteCredentialResponse]
}

// BeginRegistration calls dashboard.v1alpha1.WebAuthnService.BeginRegistration.
func (c *webAuthnServiceClient) BeginRegistration(ctx context.Context, req *connect_go.Request[v1alpha1.BeginRegistrationRequest]) (*connect_go.Response[v1alpha1.BeginRegistrationResponse], error) {
	return c.beginRegistration.CallUnary(ctx, req)
}

// FinishRegistration calls dashboard.v1alpha1.WebAuthnService.FinishRegistration.
func (c *webAuthnServiceClient) FinishRegistration(ctx context.Context, req *connect_go.Request[v1alpha1.FinishRegistrationRequest]) (*connect_go.Response[v1alpha1.FinishRegistrationResponse], error) {
	return c.finishRegistration.CallUnary(ctx, req)
}

// BeginLogin calls dashboard.v1alpha1.WebAuthnService.BeginLogin.
func (c *webAuthnServiceClient) BeginLogin(ctx context.Context, req *connect_go.Request[v1alpha1.BeginLoginRequest]) (*connect_go.Response[v1alpha1.BeginLoginResponse], error) {
	return c.beginLogin.CallUnary(ctx, req)
}

// FinishLogin calls dashboard.v1alpha1.WebAuthnService.FinishLogin.
func (c *webAuthnServiceClient) FinishLogin(ctx context.Context, req *connect_go.Request[v1alpha1.FinishLoginRequest]) (*connect_go.Response[v1alpha1.FinishLoginResponse], error) {
	return c.finishLogin.CallUnary(ctx, req)
}

// ListCredentials calls dashboard.v1alpha1.WebAuthnService.ListCredentials.
func (c *webAuthnServiceClient) ListCredentials(ctx context.Context, req *connect_go.Request[v1alpha1.ListCredentialsRequest]) (*connect_go.Response[v1alpha1.ListCredentialsResponse], error) {
	return c.listCredentials.CallUnary(ctx, req)
}

// UpdateCredential calls dashboard.v1alpha1.WebAuthnService.UpdateCredential.
func (c *webAuthnServiceClient) UpdateCredential(ctx context.Context, req *connect_go.Request[v1alpha1.UpdateCredentialRequest]) (*connect_go.Response[v1alpha1.UpdateCredentialResponse], error) {
	return c.updateCredential.CallUnary(ctx, req)
}

// DeleteCredential calls dashboard.v1alpha1.WebAuthnService.DeleteCredential.
func (c *webAuthnServiceClient) DeleteCredential(ctx context.Context, req *connect_go.Request[v1alpha1.DeleteCredentialRequest]) (*connect_go.Response[v1alpha1.DeleteCredentialResponse], error) {
	return c.deleteCredential.CallUnary(ctx, req)
}

// WebAuthnServiceHandler is an implementation of the dashboard.v1alpha1.WebAuthnService service.
type WebAuthnServiceHandler interface {
	// BeginRegistration returns CredentialCreateOptions to window.navigator.create() which is serialized as JSON string
	// Also `publicKey.user.id“ and `publicKey.challenge` are base64url encoded
	BeginRegistration(context.Context, *connect_go.Request[v1alpha1.BeginRegistrationRequest]) (*connect_go.Response[v1alpha1.BeginRegistrationResponse], error)
	// FinishRegistration check the result of window.navigator.create()
	// `rawId`, `response.clientDataJSON` and `response.attestationObject` in the result must be base64url encoded
	// and all JSON must be serialized as string
	FinishRegistration(context.Context, *connect_go.Request[v1alpha1.FinishRegistrationRequest]) (*connect_go.Response[v1alpha1.FinishRegistrationResponse], error)
	// BeginLogin returns CredentialRequestOptions to window.navigator.get() which is serialized as JSON string
	// Also `publicKey.allowCredentials[*].id` and `publicKey.challenge` are base64url encoded
	BeginLogin(context.Context, *connect_go.Request[v1alpha1.BeginLoginRequest]) (*connect_go.Response[v1alpha1.BeginLoginResponse], error)
	// FinishLogin check the result of window.navigator.get()
	// `rawId`, `response.clientDataJSON`, `response.authenticatorData`, `response.signature`, `response.userHandle`
	// in the result must be base64url encoded and all JSON must be serialized as string
	FinishLogin(context.Context, *connect_go.Request[v1alpha1.FinishLoginRequest]) (*connect_go.Response[v1alpha1.FinishLoginResponse], error)
	// ListCredentials returns registered credential ID list
	ListCredentials(context.Context, *connect_go.Request[v1alpha1.ListCredentialsRequest]) (*connect_go.Response[v1alpha1.ListCredentialsResponse], error)
	// UpdateCredential updates registed credential's human readable infomations
	UpdateCredential(context.Context, *connect_go.Request[v1alpha1.UpdateCredentialRequest]) (*connect_go.Response[v1alpha1.UpdateCredentialResponse], error)
	// DeleteCredential remove registered credential
	DeleteCredential(context.Context, *connect_go.Request[v1alpha1.DeleteCredentialRequest]) (*connect_go.Response[v1alpha1.DeleteCredentialResponse], error)
}

// NewWebAuthnServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewWebAuthnServiceHandler(svc WebAuthnServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(WebAuthnServiceBeginRegistrationProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceBeginRegistrationProcedure,
		svc.BeginRegistration,
		opts...,
	))
	mux.Handle(WebAuthnServiceFinishRegistrationProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceFinishRegistrationProcedure,
		svc.FinishRegistration,
		opts...,
	))
	mux.Handle(WebAuthnServiceBeginLoginProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceBeginLoginProcedure,
		svc.BeginLogin,
		opts...,
	))
	mux.Handle(WebAuthnServiceFinishLoginProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceFinishLoginProcedure,
		svc.FinishLogin,
		opts...,
	))
	mux.Handle(WebAuthnServiceListCredentialsProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceListCredentialsProcedure,
		svc.ListCredentials,
		opts...,
	))
	mux.Handle(WebAuthnServiceUpdateCredentialProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceUpdateCredentialProcedure,
		svc.UpdateCredential,
		opts...,
	))
	mux.Handle(WebAuthnServiceDeleteCredentialProcedure, connect_go.NewUnaryHandler(
		WebAuthnServiceDeleteCredentialProcedure,
		svc.DeleteCredential,
		opts...,
	))
	return "/dashboard.v1alpha1.WebAuthnService/", mux
}

// UnimplementedWebAuthnServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedWebAuthnServiceHandler struct{}

func (UnimplementedWebAuthnServiceHandler) BeginRegistration(context.Context, *connect_go.Request[v1alpha1.BeginRegistrationRequest]) (*connect_go.Response[v1alpha1.BeginRegistrationResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.BeginRegistration is not implemented"))
}

func (UnimplementedWebAuthnServiceHandler) FinishRegistration(context.Context, *connect_go.Request[v1alpha1.FinishRegistrationRequest]) (*connect_go.Response[v1alpha1.FinishRegistrationResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.FinishRegistration is not implemented"))
}

func (UnimplementedWebAuthnServiceHandler) BeginLogin(context.Context, *connect_go.Request[v1alpha1.BeginLoginRequest]) (*connect_go.Response[v1alpha1.BeginLoginResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.BeginLogin is not implemented"))
}

func (UnimplementedWebAuthnServiceHandler) FinishLogin(context.Context, *connect_go.Request[v1alpha1.FinishLoginRequest]) (*connect_go.Response[v1alpha1.FinishLoginResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.FinishLogin is not implemented"))
}

func (UnimplementedWebAuthnServiceHandler) ListCredentials(context.Context, *connect_go.Request[v1alpha1.ListCredentialsRequest]) (*connect_go.Response[v1alpha1.ListCredentialsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.ListCredentials is not implemented"))
}

func (UnimplementedWebAuthnServiceHandler) UpdateCredential(context.Context, *connect_go.Request[v1alpha1.UpdateCredentialRequest]) (*connect_go.Response[v1alpha1.UpdateCredentialResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.UpdateCredential is not implemented"))
}

func (UnimplementedWebAuthnServiceHandler) DeleteCredential(context.Context, *connect_go.Request[v1alpha1.DeleteCredentialRequest]) (*connect_go.Response[v1alpha1.DeleteCredentialResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.WebAuthnService.DeleteCredential is not implemented"))
}
