package dashboard

import (
	"context"
	"net/http"
	"time"

	. "github.com/cosmo-workspace/cosmo/pkg/snap"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"google.golang.org/protobuf/types/known/emptypb"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
	"github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1/dashboardv1alpha1connect"
)

var _ = Describe("Dashboard server [User]", func() {

	type SessionName string

	const (
		user       SessionName = "user"
		admin      SessionName = "admin"
		privileged SessionName = "priv"
	)

	var (
		userSession       string
		adminSession      string
		privilegedSession string
		client            dashboardv1alpha1connect.UserServiceClient
	)

	BeforeEach(func() {
		userSession = test_CreateLoginUserSession("normal-user", "お名前", []cosmov1alpha1.UserRole{{Name: "team-developer"}}, "password")
		adminSession = test_CreateLoginUserSession("admin-user", "アドミン", []cosmov1alpha1.UserRole{{Name: "team-admin"}}, "password")
		privilegedSession = test_CreateLoginUserSession("priv-user", "特権", []cosmov1alpha1.UserRole{cosmov1alpha1.PrivilegedRole}, "password")
		client = dashboardv1alpha1connect.NewUserServiceClient(http.DefaultClient, "http://localhost:8888")
	})

	AfterEach(func() {
		clientMock.Clear()
		testUtil.DeleteCosmoUserAll()
		testUtil.DeleteTemplateAll()
	})

	//==================================================================================
	userSnap := func(us *cosmov1alpha1.User) struct{ Name, Namespace, Spec, Status interface{} } {
		return struct{ Name, Namespace, Spec, Status interface{} }{
			Name:      us.Name,
			Namespace: us.Namespace,
			Spec:      us.Spec,
			Status:    us.Status,
		}
	}

	getSession := func(loginUser SessionName) string {
		if loginUser == admin {
			return adminSession
		} else if loginUser == privileged {
			return privilegedSession
		} else {
			return userSession
		}
	}

	//==================================================================================
	Describe("[CreateUser]", func() {

		withUser := func(name string) func(req *dashv1alpha1.CreateUserRequest) *time.Timer {
			return func(req *dashv1alpha1.CreateUserRequest) *time.Timer {
				testUtil.CreateUserNameSpaceandDefaultPasswordIfAbsent(name)
				_, err := client.CreateUser(ctx, NewRequestWithSession(req, privilegedSession))
				Expect(err).NotTo(HaveOccurred())
				return nil
			}
		}

		withNamespace := func(delay time.Duration) func(req *dashv1alpha1.CreateUserRequest) *time.Timer {
			return func(req *dashv1alpha1.CreateUserRequest) *time.Timer {
				createNamespace := func() {
					testUtil.CreateUserNameSpaceandDefaultPasswordIfAbsent(req.UserName)
				}
				if delay > 0 {
					return time.AfterFunc(delay, createNamespace)
				} else {
					createNamespace()
					return nil
				}
			}
		}
		withUserAddon := func(name string) func(req *dashv1alpha1.CreateUserRequest) *time.Timer {
			return func(req *dashv1alpha1.CreateUserRequest) *time.Timer {
				testUtil.CreateTemplate(cosmov1alpha1.TemplateLabelEnumTypeUserAddon, name)
				return nil
			}
		}

		run_test := func(loginUser string, req *dashv1alpha1.CreateUserRequest, beforeFuncs ...func(req *dashv1alpha1.CreateUserRequest) *time.Timer) {
			for _, f := range beforeFuncs {
				if f != nil {
					if t := f(req); t != nil {
						defer t.Stop()
					}
				}
			}

			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.CreateUser(ctx, NewRequestWithSession(req, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg.User.DefaultPassword).ShouldNot(BeEmpty())
				res.Msg.User.DefaultPassword = "xxxxxxxx"
				Ω(res.Msg).To(MatchSnapShot())
				wsv1User, err := k8sClient.GetUser(context.Background(), req.UserName)
				Expect(err).NotTo(HaveOccurred())
				Ω(userSnap(wsv1User)).To(MatchSnapShot())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
			}
			By("---------------test end---------------")
		}

		DescribeTable("✅ success create user by privileged role:",
			run_test,
			Entry("create normal user", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName:    "create-user-by-priv",
				DisplayName: "create 1",
				Roles:       []string{"team-a", "team-b"},
				AuthType:    "kosmo-secret",
				Addons: []*dashv1alpha1.UserAddons{{
					Template: "user-tmpl1",
					Vars:     map[string]string{"HOGE": "FUGA"},
				}},
			}, withNamespace(0), withUserAddon("user-tmpl1")),
			Entry("create user with only name", "priv-user", &dashv1alpha1.CreateUserRequest{UserName: "create-user-only-name"}, withNamespace(0)),
			Entry("create privileged user", "priv-user", &dashv1alpha1.CreateUserRequest{UserName: "create-user-priv-by-priv", Roles: []string{"cosmo-admin"}}, withNamespace(0)),
			Entry("create user with namespace creation short delay before timeout", "priv-user", &dashv1alpha1.CreateUserRequest{UserName: "create-user-with-delay"}, withNamespace(100*time.Millisecond)),
		)

		DescribeTable("✅ success create user by admin:",
			run_test,
			Entry("create group-developer user", "admin-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-team-by-admin",
				Roles:    []string{"team-developer"},
			}, withNamespace(0)),
			Entry("create group-admin user", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-admin-by-admin",
				Roles:    []string{"team-admin"},
			}, withNamespace(0)),
			Entry("create group-admin user", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-admin-by-admin",
				Roles:    []string{"team-developer", "team-etc"},
			}, withNamespace(0)),
		)

		DescribeTable("❌ fail with invalid request:",
			run_test,
			Entry("invalid auth type", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-invalid-auth",
				AuthType: "INVALID"}),
			Entry("no name", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName: ""}),
			Entry("user already exist", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-existing",
			}, withUser("create-user-existing")),
			Entry("including invalid charactor in username", "priv-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-INVALID"}),
		)

		DescribeTable("❌ fail with authorization by role:",
			run_test,
			Entry("normal user cannot create user", "normal-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-by-normal"}),
			Entry("admin user cannot create user including other roles", "admin-user", &dashv1alpha1.CreateUserRequest{
				UserName: "create-user-other-role-by-admin",
				Roles:    []string{"team-developer", "cosmo-admin"},
			}),
		)

		DescribeTable("❌ fail to create password timeout",
			run_test,
			Entry(nil, "priv-user", &dashv1alpha1.CreateUserRequest{UserName: "create-user-no-namespace"}),
		)
	})

	//==================================================================================
	Describe("[GetUsers]", func() {

		run_test := func(loginUser string) {
			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.GetUsers(ctx, NewRequestWithSession(&emptypb.Empty{}, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg).To(MatchSnapShot())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
			}
			By("---------------test end---------------")
		}

		DescribeTable("✅ success in normal context:",
			run_test,
			Entry(nil, "priv-user"),
		)

		// DescribeTable("✅ success with empty user:",
		// 	func(loginUser string) {
		// 		clientMock.SetListError((*Server).GetUsers, nil)
		// 		run_test(loginUser)
		// 	},
		// 	Entry(nil, "priv-user"),
		// )

		// DescribeTable("❌ fail with authorization by role:",
		// 	run_test,
		// 	Entry(nil, "normal-user"),
		// )

		// DescribeTable("❌ fail with an unexpected error at list:",
		// 	func(loginUser string) {
		// 		clientMock.SetListError((*Server).GetUsers, errors.New("mock user list error"))
		// 		run_test(loginUser)
		// 	},
		// 	Entry(nil, "priv-user"),
		// )
	})

	//==================================================================================
	Describe("[GetUser]", func() {

		_ = func(loginUser string, username string) {
			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.GetUser(ctx, NewRequestWithSession(&dashv1alpha1.GetUserRequest{UserName: username}, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg).To(MatchSnapShot())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
			}
			By("---------------test end---------------")
		}

		// DescribeTable("✅ success in normal context:",
		// 	run_test,
		// 	Entry(nil, "priv-user", "normal-user"),
		// )

		// DescribeTable("❌ fail with invalid request:",
		// 	run_test,
		// 	Entry(nil, "priv-user", "XXXXX"),
		// )

		// DescribeTable("❌ fail with authorization by role:",
		// 	run_test,
		// 	Entry(nil, "normal-user", "priv-user"),
		// )

		// DescribeTable("❌ fail with an unexpected error to get:",
		// 	func(loginUser string, username string) {
		// 		clientMock.SetGetError((*Server).GetUser, errors.New("get user error"))
		// 		run_test(loginUser, username)
		// 	},
		// 	Entry(nil, "priv-user", "normal-user"),
		// )
	})

	//==================================================================================
	Describe("[DeleteUser]", func() {

		_ = func(loginUser string, req *dashv1alpha1.DeleteUserRequest) {
			testUtil.CreateCosmoUser("user-delete1", "delete", nil)
			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.DeleteUser(ctx, NewRequestWithSession(req, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg).To(MatchSnapShot())
				_, err = k8sClient.GetUser(context.Background(), req.UserName)
				Expect(err).To(HaveOccurred())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
				_, err = k8sClient.GetUser(context.Background(), "user-delete1")
				Expect(err).NotTo(HaveOccurred())
			}
			By("---------------test end---------------")
		}

		// DescribeTable("✅ success in normal context:",
		// 	run_test,
		// 	Entry(nil, "priv-user", &dashv1alpha1.DeleteUserRequest{UserName: "user-delete1"}),
		// )

		// DescribeTable("❌ fail with invalid request:",
		// 	run_test,
		// 	Entry(nil, "priv-user", &dashv1alpha1.DeleteUserRequest{UserName: "xxxxxx"}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.DeleteUserRequest{UserName: "priv-user"}),
		// )

		// DescribeTable("❌ fail with authorization by role:",
		// 	run_test,
		// 	Entry(nil, "normal-user", &dashv1alpha1.DeleteUserRequest{UserName: "user-delete1"}),
		// )

		// DescribeTable("❌ fail with an unexpected error to get:",
		// 	func(loginUser string, req *dashv1alpha1.DeleteUserRequest) {
		// 		clientMock.SetGetError(`\.preFetchUserMiddleware\.|\.DeleteUser$`, errors.New("mock get user error")) ///
		// 		//clientMock.SetGetError((*Server).DeleteUser, errors.New("mock get user error"))
		// 		run_test(loginUser, req)
		// 	},
		// 	Entry(nil, "priv-user", &dashv1alpha1.DeleteUserRequest{UserName: "user-delete1"}),
		// )

		// DescribeTable("❌ fail with an unexpected error to delete:",
		// 	func(loginUser string, req *dashv1alpha1.DeleteUserRequest) {
		// 		clientMock.SetDeleteError((*Server).DeleteUser, errors.New("mock delete user error"))
		// 		run_test(loginUser, req)
		// 	},
		// 	Entry(nil, "priv-user", &dashv1alpha1.DeleteUserRequest{UserName: "user-delete1"}),
		// )
	})

	//==================================================================================
	Describe("[UpdateUserDisplayName]", func() {

		_ = func(loginUser string, req *dashv1alpha1.UpdateUserDisplayNameRequest) {
			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.UpdateUserDisplayName(ctx, NewRequestWithSession(req, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg).To(MatchSnapShot())
				wsv1User, err := k8sClient.GetUser(context.Background(), req.UserName)
				Expect(err).NotTo(HaveOccurred())
				Ω(userSnap(wsv1User)).To(MatchSnapShot())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
			}
			By("---------------test end---------------")
		}

		// DescribeTable("✅ success in normal context:",
		// 	run_test,
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "normal-user", DisplayName: "namechanged"}),
		// 	Entry(nil, "normal-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "normal-user", DisplayName: "namechanged"}),
		// )

		// DescribeTable("❌ fail with invalid request:",
		// 	run_test,
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "XXXXXX", DisplayName: "namechanged"}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "normal-user", DisplayName: ""}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "", DisplayName: ""}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "normal-user", DisplayName: "お名前"}),
		// )

		// DescribeTable("❌ fail with authorization by role:",
		// 	run_test,
		// 	Entry(nil, "normal-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "priv-user", DisplayName: "namechanged"}),
		// )

		// DescribeTable("❌ fail with an unexpected error to update:",
		// 	func(loginUser string, req *dashv1alpha1.UpdateUserDisplayNameRequest) {
		// 		clientMock.SetUpdateError((*Server).UpdateUserDisplayName, errors.New("mock update user error"))
		// 		run_test(loginUser, req)
		// 	},
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserDisplayNameRequest{UserName: "normal-user", DisplayName: "namechanged"}),
		// )
	})

	//==================================================================================
	Describe("[UpdateUserRole]", func() {

		_ = func(loginUser string, req *dashv1alpha1.UpdateUserRoleRequest) {
			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.UpdateUserRole(ctx, NewRequestWithSession(req, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg).To(MatchSnapShot())
				wsv1User, err := k8sClient.GetUser(context.Background(), req.UserName)
				Expect(err).NotTo(HaveOccurred())
				Ω(userSnap(wsv1User)).To(MatchSnapShot())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
				wsv1User, err := k8sClient.GetUser(context.Background(), req.UserName)
				if err != nil {
					Ω(err.Error()).To(MatchSnapShot())
				}
				if wsv1User != nil {
					Ω(userSnap(wsv1User)).To(MatchSnapShot())
				}
			}
			By("---------------test end---------------")
		}

		// DescribeTable("✅ success in normal context:",
		// 	run_test,
		// 	Entry("attach cosmo-admin to normal-user", "priv-user", &dashv1alpha1.UpdateUserRoleRequest{UserName: "normal-user", Roles: []string{"cosmo-admin"}}),
		// 	Entry("attach custom-role to normal-user", "priv-user", &dashv1alpha1.UpdateUserRoleRequest{UserName: "normal-user", Roles: []string{"xxxxx"}}),
		// 	Entry("detach role from priv-user", "priv-user", &dashv1alpha1.UpdateUserRoleRequest{UserName: "priv-user", Roles: []string{""}}),
		// )

		// DescribeTable("❌ fail with invalid request:",
		// 	run_test,
		// 	Entry("user not found", "priv-user", &dashv1alpha1.UpdateUserRoleRequest{UserName: "XXXXXX", Roles: []string{"cosmo-admin"}}),
		// 	Entry("no change", "priv-user", &dashv1alpha1.UpdateUserRoleRequest{UserName: "priv-user", Roles: []string{"cosmo-admin"}}),
		// )

		// DescribeTable("❌ fail with authorization by role:",
		// 	run_test,
		// 	Entry("permission denied", "normal-user", &dashv1alpha1.UpdateUserRoleRequest{UserName: "normal-user", Roles: []string{"cosmo-admin"}}),
		// )
	})

	//==================================================================================
	Describe("[UpdateUserPassword]", func() {

		_ = func(loginUser string, req *dashv1alpha1.UpdateUserPasswordRequest) {
			By("---------------test start----------------")
			ctx := context.Background()
			res, err := client.UpdateUserPassword(ctx, NewRequestWithSession(req, getSession(loginUser)))
			if err == nil {
				Ω(res.Msg).To(MatchSnapShot())
				verified, _, _ := k8sClient.VerifyPassword(context.Background(), req.UserName, []byte("newPassword"))
				Expect(verified).Should(BeTrue())
			} else {
				Ω(err.Error()).To(MatchSnapShot())
				Expect(res).Should(BeNil())
			}
			By("---------------test end---------------")
		}

		// DescribeTable("✅ success with invalid request:",
		// 	run_test,
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "password", NewPassword: "newPassword"}),
		// 	Entry(nil, "normal-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "normal-user", CurrentPassword: "password", NewPassword: "newPassword"}),
		// )

		// DescribeTable("❌ fail with authorization by role:",
		// 	run_test,
		// 	Entry(nil, "normal-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "password", NewPassword: "newPassword"}),
		// )

		// DescribeTable("❌ fail with invalid request:",
		// 	run_test,
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "XXXXXX", CurrentPassword: "password", NewPassword: "newPassword"}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "", NewPassword: "newPassword"}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "xxxxxx", NewPassword: "newPassword"}),
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "password", NewPassword: ""}),
		// )

		// DescribeTable("❌ fail to verify password:",
		// 	func(loginUser string, req *dashv1alpha1.UpdateUserPasswordRequest) {
		// 		clientMock.GetMock = func(ctx context.Context, key ctrl_client.ObjectKey, obj ctrl_client.Object, opts ...ctrl_client.GetOption) (mocked bool, err error) {
		// 			if key.Name == cosmov1alpha1.UserPasswordSecretName {
		// 				return true, apierrs.NewNotFound(schema.GroupResource{}, "secret")
		// 			}
		// 			return false, nil
		// 		}
		// 		//clientMock.SetGetError((*Server).PutUserPassword, apierrs.NewNotFound(schema.GroupResource{}, "secret"))
		// 		run_test(loginUser, req)
		// 	},
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "password", NewPassword: "newPassword"}),
		// )

		// DescribeTable("❌ fail with an unexpected error :",
		// 	func(loginUser string, req *dashv1alpha1.UpdateUserPasswordRequest) {
		// 		clientMock.SetUpdateError((*Server).UpdateUserPassword, errors.New("mock update error"))
		// 		run_test(loginUser, req)
		// 	},
		// 	Entry(nil, "priv-user", &dashv1alpha1.UpdateUserPasswordRequest{UserName: "priv-user", CurrentPassword: "password", NewPassword: "newPassword"}),
		// )
	})
})
