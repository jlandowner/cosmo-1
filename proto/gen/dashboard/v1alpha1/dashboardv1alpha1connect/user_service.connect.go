//
//Cosmo Dashboard API
//Manipulate cosmo dashboard resource API

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: dashboard/v1alpha1/user_service.proto

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
	// UserServiceName is the fully-qualified name of the UserService service.
	UserServiceName = "dashboard.v1alpha1.UserService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// UserServiceDeleteUserProcedure is the fully-qualified name of the UserService's DeleteUser RPC.
	UserServiceDeleteUserProcedure = "/dashboard.v1alpha1.UserService/DeleteUser"
	// UserServiceGetUserProcedure is the fully-qualified name of the UserService's GetUser RPC.
	UserServiceGetUserProcedure = "/dashboard.v1alpha1.UserService/GetUser"
	// UserServiceGetUsersProcedure is the fully-qualified name of the UserService's GetUsers RPC.
	UserServiceGetUsersProcedure = "/dashboard.v1alpha1.UserService/GetUsers"
	// UserServiceGetEventsProcedure is the fully-qualified name of the UserService's GetEvents RPC.
	UserServiceGetEventsProcedure = "/dashboard.v1alpha1.UserService/GetEvents"
	// UserServiceCreateUserProcedure is the fully-qualified name of the UserService's CreateUser RPC.
	UserServiceCreateUserProcedure = "/dashboard.v1alpha1.UserService/CreateUser"
	// UserServiceUpdateUserDisplayNameProcedure is the fully-qualified name of the UserService's
	// UpdateUserDisplayName RPC.
	UserServiceUpdateUserDisplayNameProcedure = "/dashboard.v1alpha1.UserService/UpdateUserDisplayName"
	// UserServiceUpdateUserPasswordProcedure is the fully-qualified name of the UserService's
	// UpdateUserPassword RPC.
	UserServiceUpdateUserPasswordProcedure = "/dashboard.v1alpha1.UserService/UpdateUserPassword"
	// UserServiceUpdateUserRoleProcedure is the fully-qualified name of the UserService's
	// UpdateUserRole RPC.
	UserServiceUpdateUserRoleProcedure = "/dashboard.v1alpha1.UserService/UpdateUserRole"
	// UserServiceUpdateUserAddonsProcedure is the fully-qualified name of the UserService's
	// UpdateUserAddons RPC.
	UserServiceUpdateUserAddonsProcedure = "/dashboard.v1alpha1.UserService/UpdateUserAddons"
	// UserServiceUpdateUserDeletePolicyProcedure is the fully-qualified name of the UserService's
	// UpdateUserDeletePolicy RPC.
	UserServiceUpdateUserDeletePolicyProcedure = "/dashboard.v1alpha1.UserService/UpdateUserDeletePolicy"
)

// UserServiceClient is a client for the dashboard.v1alpha1.UserService service.
type UserServiceClient interface {
	// Delete user by ID
	DeleteUser(context.Context, *connect_go.Request[v1alpha1.DeleteUserRequest]) (*connect_go.Response[v1alpha1.DeleteUserResponse], error)
	// Returns a single User model
	GetUser(context.Context, *connect_go.Request[v1alpha1.GetUserRequest]) (*connect_go.Response[v1alpha1.GetUserResponse], error)
	// Returns an array of User model
	GetUsers(context.Context, *connect_go.Request[v1alpha1.GetUsersRequest]) (*connect_go.Response[v1alpha1.GetUsersResponse], error)
	// Returns events for User
	GetEvents(context.Context, *connect_go.Request[v1alpha1.GetEventsRequest]) (*connect_go.Response[v1alpha1.GetEventsResponse], error)
	// Create a new User
	CreateUser(context.Context, *connect_go.Request[v1alpha1.CreateUserRequest]) (*connect_go.Response[v1alpha1.CreateUserResponse], error)
	// Update user display name
	UpdateUserDisplayName(context.Context, *connect_go.Request[v1alpha1.UpdateUserDisplayNameRequest]) (*connect_go.Response[v1alpha1.UpdateUserDisplayNameResponse], error)
	// Update a single User password
	UpdateUserPassword(context.Context, *connect_go.Request[v1alpha1.UpdateUserPasswordRequest]) (*connect_go.Response[v1alpha1.UpdateUserPasswordResponse], error)
	// Update a single User role
	UpdateUserRole(context.Context, *connect_go.Request[v1alpha1.UpdateUserRoleRequest]) (*connect_go.Response[v1alpha1.UpdateUserRoleResponse], error)
	// Update a single User role
	UpdateUserAddons(context.Context, *connect_go.Request[v1alpha1.UpdateUserAddonsRequest]) (*connect_go.Response[v1alpha1.UpdateUserAddonsResponse], error)
	// Update user delete policy
	UpdateUserDeletePolicy(context.Context, *connect_go.Request[v1alpha1.UpdateUserDeletePolicyRequest]) (*connect_go.Response[v1alpha1.UpdateUserDeletePolicyResponse], error)
}

// NewUserServiceClient constructs a client for the dashboard.v1alpha1.UserService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewUserServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) UserServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &userServiceClient{
		deleteUser: connect_go.NewClient[v1alpha1.DeleteUserRequest, v1alpha1.DeleteUserResponse](
			httpClient,
			baseURL+UserServiceDeleteUserProcedure,
			opts...,
		),
		getUser: connect_go.NewClient[v1alpha1.GetUserRequest, v1alpha1.GetUserResponse](
			httpClient,
			baseURL+UserServiceGetUserProcedure,
			opts...,
		),
		getUsers: connect_go.NewClient[v1alpha1.GetUsersRequest, v1alpha1.GetUsersResponse](
			httpClient,
			baseURL+UserServiceGetUsersProcedure,
			opts...,
		),
		getEvents: connect_go.NewClient[v1alpha1.GetEventsRequest, v1alpha1.GetEventsResponse](
			httpClient,
			baseURL+UserServiceGetEventsProcedure,
			opts...,
		),
		createUser: connect_go.NewClient[v1alpha1.CreateUserRequest, v1alpha1.CreateUserResponse](
			httpClient,
			baseURL+UserServiceCreateUserProcedure,
			opts...,
		),
		updateUserDisplayName: connect_go.NewClient[v1alpha1.UpdateUserDisplayNameRequest, v1alpha1.UpdateUserDisplayNameResponse](
			httpClient,
			baseURL+UserServiceUpdateUserDisplayNameProcedure,
			opts...,
		),
		updateUserPassword: connect_go.NewClient[v1alpha1.UpdateUserPasswordRequest, v1alpha1.UpdateUserPasswordResponse](
			httpClient,
			baseURL+UserServiceUpdateUserPasswordProcedure,
			opts...,
		),
		updateUserRole: connect_go.NewClient[v1alpha1.UpdateUserRoleRequest, v1alpha1.UpdateUserRoleResponse](
			httpClient,
			baseURL+UserServiceUpdateUserRoleProcedure,
			opts...,
		),
		updateUserAddons: connect_go.NewClient[v1alpha1.UpdateUserAddonsRequest, v1alpha1.UpdateUserAddonsResponse](
			httpClient,
			baseURL+UserServiceUpdateUserAddonsProcedure,
			opts...,
		),
		updateUserDeletePolicy: connect_go.NewClient[v1alpha1.UpdateUserDeletePolicyRequest, v1alpha1.UpdateUserDeletePolicyResponse](
			httpClient,
			baseURL+UserServiceUpdateUserDeletePolicyProcedure,
			opts...,
		),
	}
}

// userServiceClient implements UserServiceClient.
type userServiceClient struct {
	deleteUser             *connect_go.Client[v1alpha1.DeleteUserRequest, v1alpha1.DeleteUserResponse]
	getUser                *connect_go.Client[v1alpha1.GetUserRequest, v1alpha1.GetUserResponse]
	getUsers               *connect_go.Client[v1alpha1.GetUsersRequest, v1alpha1.GetUsersResponse]
	getEvents              *connect_go.Client[v1alpha1.GetEventsRequest, v1alpha1.GetEventsResponse]
	createUser             *connect_go.Client[v1alpha1.CreateUserRequest, v1alpha1.CreateUserResponse]
	updateUserDisplayName  *connect_go.Client[v1alpha1.UpdateUserDisplayNameRequest, v1alpha1.UpdateUserDisplayNameResponse]
	updateUserPassword     *connect_go.Client[v1alpha1.UpdateUserPasswordRequest, v1alpha1.UpdateUserPasswordResponse]
	updateUserRole         *connect_go.Client[v1alpha1.UpdateUserRoleRequest, v1alpha1.UpdateUserRoleResponse]
	updateUserAddons       *connect_go.Client[v1alpha1.UpdateUserAddonsRequest, v1alpha1.UpdateUserAddonsResponse]
	updateUserDeletePolicy *connect_go.Client[v1alpha1.UpdateUserDeletePolicyRequest, v1alpha1.UpdateUserDeletePolicyResponse]
}

// DeleteUser calls dashboard.v1alpha1.UserService.DeleteUser.
func (c *userServiceClient) DeleteUser(ctx context.Context, req *connect_go.Request[v1alpha1.DeleteUserRequest]) (*connect_go.Response[v1alpha1.DeleteUserResponse], error) {
	return c.deleteUser.CallUnary(ctx, req)
}

// GetUser calls dashboard.v1alpha1.UserService.GetUser.
func (c *userServiceClient) GetUser(ctx context.Context, req *connect_go.Request[v1alpha1.GetUserRequest]) (*connect_go.Response[v1alpha1.GetUserResponse], error) {
	return c.getUser.CallUnary(ctx, req)
}

// GetUsers calls dashboard.v1alpha1.UserService.GetUsers.
func (c *userServiceClient) GetUsers(ctx context.Context, req *connect_go.Request[v1alpha1.GetUsersRequest]) (*connect_go.Response[v1alpha1.GetUsersResponse], error) {
	return c.getUsers.CallUnary(ctx, req)
}

// GetEvents calls dashboard.v1alpha1.UserService.GetEvents.
func (c *userServiceClient) GetEvents(ctx context.Context, req *connect_go.Request[v1alpha1.GetEventsRequest]) (*connect_go.Response[v1alpha1.GetEventsResponse], error) {
	return c.getEvents.CallUnary(ctx, req)
}

// CreateUser calls dashboard.v1alpha1.UserService.CreateUser.
func (c *userServiceClient) CreateUser(ctx context.Context, req *connect_go.Request[v1alpha1.CreateUserRequest]) (*connect_go.Response[v1alpha1.CreateUserResponse], error) {
	return c.createUser.CallUnary(ctx, req)
}

// UpdateUserDisplayName calls dashboard.v1alpha1.UserService.UpdateUserDisplayName.
func (c *userServiceClient) UpdateUserDisplayName(ctx context.Context, req *connect_go.Request[v1alpha1.UpdateUserDisplayNameRequest]) (*connect_go.Response[v1alpha1.UpdateUserDisplayNameResponse], error) {
	return c.updateUserDisplayName.CallUnary(ctx, req)
}

// UpdateUserPassword calls dashboard.v1alpha1.UserService.UpdateUserPassword.
func (c *userServiceClient) UpdateUserPassword(ctx context.Context, req *connect_go.Request[v1alpha1.UpdateUserPasswordRequest]) (*connect_go.Response[v1alpha1.UpdateUserPasswordResponse], error) {
	return c.updateUserPassword.CallUnary(ctx, req)
}

// UpdateUserRole calls dashboard.v1alpha1.UserService.UpdateUserRole.
func (c *userServiceClient) UpdateUserRole(ctx context.Context, req *connect_go.Request[v1alpha1.UpdateUserRoleRequest]) (*connect_go.Response[v1alpha1.UpdateUserRoleResponse], error) {
	return c.updateUserRole.CallUnary(ctx, req)
}

// UpdateUserAddons calls dashboard.v1alpha1.UserService.UpdateUserAddons.
func (c *userServiceClient) UpdateUserAddons(ctx context.Context, req *connect_go.Request[v1alpha1.UpdateUserAddonsRequest]) (*connect_go.Response[v1alpha1.UpdateUserAddonsResponse], error) {
	return c.updateUserAddons.CallUnary(ctx, req)
}

// UpdateUserDeletePolicy calls dashboard.v1alpha1.UserService.UpdateUserDeletePolicy.
func (c *userServiceClient) UpdateUserDeletePolicy(ctx context.Context, req *connect_go.Request[v1alpha1.UpdateUserDeletePolicyRequest]) (*connect_go.Response[v1alpha1.UpdateUserDeletePolicyResponse], error) {
	return c.updateUserDeletePolicy.CallUnary(ctx, req)
}

// UserServiceHandler is an implementation of the dashboard.v1alpha1.UserService service.
type UserServiceHandler interface {
	// Delete user by ID
	DeleteUser(context.Context, *connect_go.Request[v1alpha1.DeleteUserRequest]) (*connect_go.Response[v1alpha1.DeleteUserResponse], error)
	// Returns a single User model
	GetUser(context.Context, *connect_go.Request[v1alpha1.GetUserRequest]) (*connect_go.Response[v1alpha1.GetUserResponse], error)
	// Returns an array of User model
	GetUsers(context.Context, *connect_go.Request[v1alpha1.GetUsersRequest]) (*connect_go.Response[v1alpha1.GetUsersResponse], error)
	// Returns events for User
	GetEvents(context.Context, *connect_go.Request[v1alpha1.GetEventsRequest]) (*connect_go.Response[v1alpha1.GetEventsResponse], error)
	// Create a new User
	CreateUser(context.Context, *connect_go.Request[v1alpha1.CreateUserRequest]) (*connect_go.Response[v1alpha1.CreateUserResponse], error)
	// Update user display name
	UpdateUserDisplayName(context.Context, *connect_go.Request[v1alpha1.UpdateUserDisplayNameRequest]) (*connect_go.Response[v1alpha1.UpdateUserDisplayNameResponse], error)
	// Update a single User password
	UpdateUserPassword(context.Context, *connect_go.Request[v1alpha1.UpdateUserPasswordRequest]) (*connect_go.Response[v1alpha1.UpdateUserPasswordResponse], error)
	// Update a single User role
	UpdateUserRole(context.Context, *connect_go.Request[v1alpha1.UpdateUserRoleRequest]) (*connect_go.Response[v1alpha1.UpdateUserRoleResponse], error)
	// Update a single User role
	UpdateUserAddons(context.Context, *connect_go.Request[v1alpha1.UpdateUserAddonsRequest]) (*connect_go.Response[v1alpha1.UpdateUserAddonsResponse], error)
	// Update user delete policy
	UpdateUserDeletePolicy(context.Context, *connect_go.Request[v1alpha1.UpdateUserDeletePolicyRequest]) (*connect_go.Response[v1alpha1.UpdateUserDeletePolicyResponse], error)
}

// NewUserServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewUserServiceHandler(svc UserServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(UserServiceDeleteUserProcedure, connect_go.NewUnaryHandler(
		UserServiceDeleteUserProcedure,
		svc.DeleteUser,
		opts...,
	))
	mux.Handle(UserServiceGetUserProcedure, connect_go.NewUnaryHandler(
		UserServiceGetUserProcedure,
		svc.GetUser,
		opts...,
	))
	mux.Handle(UserServiceGetUsersProcedure, connect_go.NewUnaryHandler(
		UserServiceGetUsersProcedure,
		svc.GetUsers,
		opts...,
	))
	mux.Handle(UserServiceGetEventsProcedure, connect_go.NewUnaryHandler(
		UserServiceGetEventsProcedure,
		svc.GetEvents,
		opts...,
	))
	mux.Handle(UserServiceCreateUserProcedure, connect_go.NewUnaryHandler(
		UserServiceCreateUserProcedure,
		svc.CreateUser,
		opts...,
	))
	mux.Handle(UserServiceUpdateUserDisplayNameProcedure, connect_go.NewUnaryHandler(
		UserServiceUpdateUserDisplayNameProcedure,
		svc.UpdateUserDisplayName,
		opts...,
	))
	mux.Handle(UserServiceUpdateUserPasswordProcedure, connect_go.NewUnaryHandler(
		UserServiceUpdateUserPasswordProcedure,
		svc.UpdateUserPassword,
		opts...,
	))
	mux.Handle(UserServiceUpdateUserRoleProcedure, connect_go.NewUnaryHandler(
		UserServiceUpdateUserRoleProcedure,
		svc.UpdateUserRole,
		opts...,
	))
	mux.Handle(UserServiceUpdateUserAddonsProcedure, connect_go.NewUnaryHandler(
		UserServiceUpdateUserAddonsProcedure,
		svc.UpdateUserAddons,
		opts...,
	))
	mux.Handle(UserServiceUpdateUserDeletePolicyProcedure, connect_go.NewUnaryHandler(
		UserServiceUpdateUserDeletePolicyProcedure,
		svc.UpdateUserDeletePolicy,
		opts...,
	))
	return "/dashboard.v1alpha1.UserService/", mux
}

// UnimplementedUserServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedUserServiceHandler struct{}

func (UnimplementedUserServiceHandler) DeleteUser(context.Context, *connect_go.Request[v1alpha1.DeleteUserRequest]) (*connect_go.Response[v1alpha1.DeleteUserResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.DeleteUser is not implemented"))
}

func (UnimplementedUserServiceHandler) GetUser(context.Context, *connect_go.Request[v1alpha1.GetUserRequest]) (*connect_go.Response[v1alpha1.GetUserResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.GetUser is not implemented"))
}

func (UnimplementedUserServiceHandler) GetUsers(context.Context, *connect_go.Request[v1alpha1.GetUsersRequest]) (*connect_go.Response[v1alpha1.GetUsersResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.GetUsers is not implemented"))
}

func (UnimplementedUserServiceHandler) GetEvents(context.Context, *connect_go.Request[v1alpha1.GetEventsRequest]) (*connect_go.Response[v1alpha1.GetEventsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.GetEvents is not implemented"))
}

func (UnimplementedUserServiceHandler) CreateUser(context.Context, *connect_go.Request[v1alpha1.CreateUserRequest]) (*connect_go.Response[v1alpha1.CreateUserResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.CreateUser is not implemented"))
}

func (UnimplementedUserServiceHandler) UpdateUserDisplayName(context.Context, *connect_go.Request[v1alpha1.UpdateUserDisplayNameRequest]) (*connect_go.Response[v1alpha1.UpdateUserDisplayNameResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.UpdateUserDisplayName is not implemented"))
}

func (UnimplementedUserServiceHandler) UpdateUserPassword(context.Context, *connect_go.Request[v1alpha1.UpdateUserPasswordRequest]) (*connect_go.Response[v1alpha1.UpdateUserPasswordResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.UpdateUserPassword is not implemented"))
}

func (UnimplementedUserServiceHandler) UpdateUserRole(context.Context, *connect_go.Request[v1alpha1.UpdateUserRoleRequest]) (*connect_go.Response[v1alpha1.UpdateUserRoleResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.UpdateUserRole is not implemented"))
}

func (UnimplementedUserServiceHandler) UpdateUserAddons(context.Context, *connect_go.Request[v1alpha1.UpdateUserAddonsRequest]) (*connect_go.Response[v1alpha1.UpdateUserAddonsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.UpdateUserAddons is not implemented"))
}

func (UnimplementedUserServiceHandler) UpdateUserDeletePolicy(context.Context, *connect_go.Request[v1alpha1.UpdateUserDeletePolicyRequest]) (*connect_go.Response[v1alpha1.UpdateUserDeletePolicyResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("dashboard.v1alpha1.UserService.UpdateUserDeletePolicy is not implemented"))
}
