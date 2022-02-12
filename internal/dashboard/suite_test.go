package dashboard

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/core/v1alpha1"
	wsv1alpha1 "github.com/cosmo-workspace/cosmo/api/workspace/v1alpha1"

	//+kubebuilder:scaffold:imports

	"github.com/cosmo-workspace/cosmo/pkg/auth"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/kosmo"
	"github.com/cosmo-workspace/cosmo/pkg/wsnet"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg       *rest.Config
	k8sClient kosmo.Client
	testEnv   *envtest.Environment

	userSession  []*http.Cookie
	adminSession []*http.Cookie
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Dashboard Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = cosmov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = wsv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	c, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())

	k8sClient = kosmo.NewClient(c)

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
	})
	Expect(err).NotTo(HaveOccurred())

	// Setup server
	By("bootstrapping server")
	klient := kosmo.NewClient(mgr.GetClient())

	auths := make(map[wsv1alpha1.UserAuthType]auth.Authorizer)
	auths[wsv1alpha1.UserAuthTypeKosmoSecert] = auth.NewKosmoSecretAuthorizer(klient)

	serv := (&Server{
		Log:                 clog.NewLogger(ctrl.Log.WithName("dashboard")),
		Klient:              klient,
		GracefulShutdownDur: time.Second * time.Duration(3),
		ResponseTimeout:     time.Second * time.Duration(3),
		StaticFileDir:       filepath.Join(".", "test"),
		Port:                8888,
		MaxAgeSeconds:       60,
		SessionName:         "test-server",
		Insecure:            true,
		Authorizers:         auths,
	})
	err = mgr.Add(serv)
	Expect(err).NotTo(HaveOccurred())

	ctx := ctrl.SetupSignalHandler()

	go func() {
		err = mgr.Start(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()

	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

type request struct {
	method string
	path   string
	query  map[string]string
	body   string
}
type response struct {
	statusCode int
	body       string
	bodyValues func() []interface{}
}

func test_HttpSend(session []*http.Cookie, reqParam request) (httpResponse *http.Response, responseBody []byte) {
	req, err := http.NewRequest(reqParam.method, "http://localhost:8888"+reqParam.path, bytes.NewBufferString(reqParam.body))
	Expect(err).NotTo(HaveOccurred())

	for _, v := range session {
		req.AddCookie(v)
	}

	got, err := http.DefaultClient.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer got.Body.Close()
	gotBody, err := ioutil.ReadAll(got.Body)
	Expect(err).NotTo(HaveOccurred())

	return got, gotBody
}

func test_HttpSendAndVerify(session []*http.Cookie, reqParam request, expect response) (responseBody []byte, statusCode int) {

	req, err := http.NewRequest(reqParam.method, "http://localhost:8888"+reqParam.path, bytes.NewBufferString(reqParam.body))
	Expect(err).NotTo(HaveOccurred())

	for _, v := range session {
		req.AddCookie(v)
	}

	got, err := http.DefaultClient.Do(req)
	Expect(err).NotTo(HaveOccurred())

	Expect(got).Should(HaveHTTPStatus(expect.statusCode))

	defer got.Body.Close()
	gotBody, err := ioutil.ReadAll(got.Body)
	Expect(err).NotTo(HaveOccurred())

	body := expect.body
	if expect.bodyValues != nil {
		body = fmt.Sprintf(body, expect.bodyValues()...)
	}

	if body == "@ignore" {
		// no assert
	} else if body != "" && body[0:1] == "{" {
		Expect(body).Should(MatchJSON(gotBody))
	} else {
		Expect(string(body)).Should(Equal(string(gotBody)))
	}
	return gotBody, got.StatusCode
}

func test_CreateTemplate(templateType string, templateName string) {
	ctx := context.Background()
	inst := cosmov1alpha1.Template{
		ObjectMeta: metav1.ObjectMeta{
			Name: templateName,
			Labels: map[string]string{
				cosmov1alpha1.TemplateLabelKeyType: templateType,
			},
		},
		Spec: cosmov1alpha1.TemplateSpec{
			RequiredVars: []cosmov1alpha1.RequiredVarSpec{
				{Var: "{{HOGE}}", Default: "FUGA"},
			},
		},
	}
	err := k8sClient.Create(ctx, &inst)
	Expect(err).ShouldNot(HaveOccurred())

	Eventually(func() error {
		err := k8sClient.Get(ctx, client.ObjectKey{Name: templateName}, &cosmov1alpha1.Template{})
		return err
	}, time.Second*5, time.Millisecond*100).Should(Succeed())
}

func test_DeleteTemplateAll() {
	ctx := context.Background()
	templates, err := k8sClient.ListTemplates(ctx)
	Expect(err).ShouldNot(HaveOccurred())
	for _, template := range templates {
		k8sClient.Delete(ctx, &template)
	}
	Eventually(func() ([]cosmov1alpha1.Template, error) {
		return k8sClient.ListTemplates(ctx)
	}, time.Second*5, time.Millisecond*100).Should(BeEmpty())
}

func test_CreateCosmoUser(id string, dispayName string, role wsv1alpha1.UserRole) {
	ctx := context.Background()
	user := wsv1alpha1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: wsv1alpha1.UserSpec{
			DisplayName: dispayName,
			Role:        role,
			AuthType:    wsv1alpha1.UserAuthTypeKosmoSecert,
		},
	}
	err := k8sClient.Create(ctx, &user)
	Expect(err).ShouldNot(HaveOccurred())

	Eventually(func() error {
		err := k8sClient.Get(ctx, client.ObjectKey{Name: id}, &wsv1alpha1.User{})
		return err
	}, time.Second*5, time.Millisecond*100).Should(Succeed())
}

func test_DeleteCosmoUserAll() {
	ctx := context.Background()
	users, err := k8sClient.ListUsers(ctx)
	Expect(err).ShouldNot(HaveOccurred())
	for _, user := range users {
		k8sClient.Delete(ctx, &user)
	}
	Eventually(func() ([]wsv1alpha1.User, error) {
		return k8sClient.ListUsers(ctx)
	}, time.Second*5, time.Millisecond*100).Should(BeEmpty())
}

func test_CreateUserNameSpaceandDefaultPasswordIfAbsent(id string) {
	ctx := context.Background()
	var ns v1.Namespace
	key := client.ObjectKey{Name: wsv1alpha1.UserNamespace(id)}
	err := k8sClient.Get(ctx, key, &ns)
	if err != nil {
		ns = v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: wsv1alpha1.UserNamespace(id),
			},
		}
		err = k8sClient.Create(ctx, &ns)
		Expect(err).ShouldNot(HaveOccurred())
	}
	// create default password
	err = k8sClient.ResetPassword(ctx, id)
	Expect(err).ShouldNot(HaveOccurred())
}

func test_CreateLoginUserSession(id, displayName string, role wsv1alpha1.UserRole, password string) (session []*http.Cookie) {
	ctx := context.Background()

	test_CreateCosmoUser(id, displayName, role)
	test_CreateUserNameSpaceandDefaultPasswordIfAbsent(id)
	err := k8sClient.RegisterPassword(ctx, id, []byte(password))
	Expect(err).ShouldNot(HaveOccurred())

	Eventually(func() error {
		_, _, err := k8sClient.VerifyPassword(ctx, id, []byte(password))
		return err
	}, time.Second*5, time.Millisecond*100).Should(Succeed())

	return test_Login(id, password)
}

func test_Login(userId string, password string) []*http.Cookie {
	got, _ := test_HttpSend(nil, request{
		method: http.MethodPost, path: "/api/v1alpha1/auth/login", body: `{"id": "` + userId + `", "password": "` + password + `"}`,
	})
	Expect(http.StatusOK).Should(Equal(got.StatusCode))

	return got.Cookies()
}

func test_CreateWorkspace(userId string, name string, template string, vars map[string]string) {
	ctx := context.Background()

	ws := &wsv1alpha1.Workspace{}
	ws.SetName(name)
	ws.SetNamespace(wsv1alpha1.UserNamespace(userId))
	ws.Spec = wsv1alpha1.WorkspaceSpec{
		Template: cosmov1alpha1.TemplateRef{
			Name: template,
		},
		Replicas: pointer.Int64(0),
		Vars:     vars,
	}

	err := k8sClient.Create(ctx, ws)
	Expect(err).ShouldNot(HaveOccurred())

	Eventually(func() (*wsv1alpha1.Workspace, error) {
		return k8sClient.GetWorkspaceByUserID(ctx, name, userId)
	}, time.Second*5, time.Millisecond*100).ShouldNot(BeNil())
}

func test_DeleteWorkspaceAllByUserId(userId string) {
	ctx := context.Background()
	workspaces, err := k8sClient.ListWorkspaces(ctx, wsv1alpha1.UserNamespace(userId))
	Expect(err).ShouldNot(HaveOccurred())
	for _, workspace := range workspaces {
		k8sClient.Delete(ctx, &workspace)
	}
	Eventually(func() ([]wsv1alpha1.Workspace, error) {
		return k8sClient.ListWorkspaces(ctx, wsv1alpha1.UserNamespace(userId))
	}, time.Second*5, time.Millisecond*100).Should(BeEmpty())
}

func test_DeleteWorkspaceAll() {
	ctx := context.Background()
	users, err := k8sClient.ListUsers(ctx)
	Expect(err).ShouldNot(HaveOccurred())
	for _, user := range users {
		test_DeleteWorkspaceAllByUserId(user.Name)
	}
}

func test_createNetworkRule(userId, workspaceName, networkRuleName string, portNumber int, group, httpPath string) {
	ctx := context.Background()

	ws, err := k8sClient.GetWorkspaceByUserID(ctx, workspaceName, userId)
	Expect(err).ShouldNot(HaveOccurred())

	netRule := wsv1alpha1.NetworkRule{
		PortName:   networkRuleName,
		PortNumber: portNumber,
		Group:      pointer.String(group),
		HTTPPath:   httpPath,
	}

	ws.Spec.Network, err = wsnet.UpsertNetRule(ws.Spec.Network, netRule)
	Expect(err).ShouldNot(HaveOccurred())

	err = k8sClient.Update(ctx, ws)
	Expect(err).ShouldNot(HaveOccurred())

	Eventually(func() []wsv1alpha1.NetworkRule {
		ws, _ := k8sClient.GetWorkspaceByUserID(ctx, workspaceName, userId)
		return ws.Spec.Network
	}, time.Second*5, time.Millisecond*100).Should(ContainElement(netRule))
}
