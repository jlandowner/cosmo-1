package cmdutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/bufbuild/connect-go"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/kosmo"
)

type GetOutputFormat string

const (
	GetOutputFormatWide GetOutputFormat = "wide"
	GetOutputFormatJSON GetOutputFormat = "json"
)

var HelpGetOutputFormat string = fmt.Sprintf("available: [%s, %s]", GetOutputFormatWide, GetOutputFormatJSON)

func (o GetOutputFormat) String() string {
	return string(o)
}

func (o GetOutputFormat) Validate() error {
	switch o {
	case GetOutputFormatWide, GetOutputFormatJSON:
		return nil
	case "":
		return nil
	default:
		return fmt.Errorf("invalid output format '%s' %s", o, HelpGetOutputFormat)
	}
}

// CliConfig has client common infomation to communicate with the server
// which is loaded by incluster auto completion or from login config file
type CliConfig struct {
	LoginUser      string `json:"user,omitempty"`
	ServerCA       []byte `json:"ca,omitempty"`
	ServerEndpoint string `json:"endpoint,omitempty"`
	Token          string `json:"token,omitempty"`
	Cookie         string `json:"cookie,omitempty"`
}

func (c *CliConfig) CompleteInCluster() error {
	var err error
	c.LoginUser, err = loadInclusterUserName()
	if err != nil {
		return err
	}
	c.ServerEndpoint = os.Getenv(cosmov1alpha1.EnvServerEndpoint)
	if c.ServerEndpoint == "" {
		return fmt.Errorf("failed to get env %s", cosmov1alpha1.EnvServerEndpoint)
	}
	c.ServerCA, err = loadInclusterServerCA()
	if err != nil {
		return err
	}
	c.Token, err = loadInclusterServiceAccountToken()
	if err != nil {
		return err
	}
	return nil
}

func loadInclusterUserName() (string, error) {
	const inclusterNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	ns, err := ioutil.ReadFile(inclusterNamespaceFile)
	if len(ns) == 0 || err != nil {
		return "", fmt.Errorf("failed to load incluster namespace: %w", err)
	}
	userName := cosmov1alpha1.UserNameByNamespace(string(ns))
	if userName == "" {
		return "", fmt.Errorf("not cosmo user namespace: %s", ns)
	}
	return userName, nil
}

func loadInclusterServerCA() ([]byte, error) {
	v := os.Getenv(cosmov1alpha1.EnvServerCA)
	if v == "" {
		return nil, errors.New("env not found")
	}
	return base64.RawStdEncoding.DecodeString(v)
}

func loadInclusterServiceAccountToken() (string, error) {
	const inclusterTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	token, err := ioutil.ReadFile(inclusterTokenFile)
	if len(token) == 0 || err != nil {
		return "", fmt.Errorf("failed to load incluster token: %w", err)
	}
	return string(token), nil
}

func (c *CliConfig) CompleteByConfigFile(filepath string) error {
	f, err := ioutil.ReadFile(filepath)
	if len(f) == 0 || err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	buf := make([]byte, base64.RawStdEncoding.DecodedLen(len(f)))
	if n, err := base64.RawStdEncoding.Decode(buf, f); n == 0 || err != nil {
		return fmt.Errorf("failed to decode file: %w", err)
	}

	if err = json.Unmarshal(f, c); err != nil {
		return fmt.Errorf("failed to decode string: %w", err)
	}
	return nil
}

func (c *CliConfig) Write(filename string) error {
	r, err := json.Marshal(c)
	if err != nil {
		return err
	}
	buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(r)))
	base64.RawStdEncoding.Encode(buf, r)
	return os.WriteFile(filename, buf, 0600)
}

// CliOptions is a common infomations for cli
type CliOptions struct {
	*CliConfig

	LogLevel          int
	CliConfigFilePath string
	Insecure          bool

	Ctx    context.Context
	Client connect.HTTPClient

	UseKubeClient  bool
	KubeClient     *kosmo.Client
	KubeScheme     *runtime.Scheme
	KubeConfigPath string
	KubeContext    string

	Logr   *clog.Logger
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

func (o *CliOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&o.LogLevel, "verbose", "v", -1, "log level. -1:DISABLED, 0:INFO, 1:DEBUG, 2:ALL")
	cmd.PersistentFlags().StringVar(&o.CliConfigFilePath, "config", "$HOME/.cosmoctl", "config file path")
	cmd.PersistentFlags().BoolVar(&o.Insecure, "insecure", false, "use insecure client")
}

func (o *CliOptions) AddKubeClientFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.UseKubeClient, "kube", "k", false, "use kube rbac authentication instead of cosmo user authentication")
	cmd.PersistentFlags().StringVar(&o.KubeContext, "kubecontext", "", "kube context")
	cmd.PersistentFlags().StringVar(&o.KubeConfigPath, "kubeconfig", "$HOME/.kube/config", "kubeconfig file path")
}

func NewCliOptions() *CliOptions {
	ctx := context.TODO()
	return &CliOptions{Ctx: ctx, CliConfig: &CliConfig{}, In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
}

func (o *CliOptions) CompleteLogger() {
	if o.LogLevel >= 0 {
		opt := zap.Options{
			Development: true,
			Level:       zapcore.Level(-o.LogLevel),
			TimeEncoder: zapcore.ISO8601TimeEncoder,
		}
		o.Logr = clog.NewLogger(zap.New(zap.UseFlagOptions(&opt)))
		o.Ctx = clog.IntoContext(o.Ctx, o.Logr)
	} else {
		o.Logr = clog.NewLogger(logr.Discard())
	}
}

func (o *CliOptions) CompleteCliConfig() error {
	// First, load config from file
	// If failed to load config file, then load incluster config
	// If both failed, return err
	if o.CliConfigFilePath == "" || o.CliConfigFilePath == "$HOME/.cosmoctl" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		o.CliConfigFilePath = path.Join(home, ".cosmoctl")
	}

	if err := o.CliConfig.CompleteByConfigFile(o.CliConfigFilePath); err != nil {
		err = fmt.Errorf("failed to load config file %s: %w", o.CliConfigFilePath, err)

		if ierr := o.CliConfig.CompleteInCluster(); ierr != nil {
			return fmt.Errorf("failed to load config %v: %w", ierr, err)
		}
	}
	return nil
}

func (o *CliOptions) CompleteKubeClient() error {
	cfgFlg := genericclioptions.NewConfigFlags(true)
	if o.KubeConfigPath != "" && o.KubeConfigPath != "$HOME/.kube/config" {
		cfgFlg.KubeConfig = &o.KubeConfigPath
	}
	if o.KubeContext != "" {
		cfgFlg.Context = &o.KubeContext
	}

	cfg, err := cfgFlg.ToRESTConfig()
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cosmov1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme

	baseclient, err := kosmo.NewClientByRestConfig(cfg, scheme)
	if err != nil {
		return err
	}
	o.KubeClient = &baseclient
	o.KubeScheme = scheme
	return nil
}

func (o *CliOptions) CompleteClient() error {
	if o.Insecure {
		o.Client = http.DefaultClient
	} else {
		caPool, err := x509.SystemCertPool()
		if err != nil {
			return err
		}

		if len(o.ServerCA) > 0 {
			caPool.AppendCertsFromPEM(o.ServerCA)
		}

		t := http.DefaultTransport.(*http.Transport).Clone()
		t.TLSClientConfig = &tls.Config{
			RootCAs:            caPool,
			InsecureSkipVerify: true,
		}
		o.Client = &http.Client{Transport: t}
	}
	return nil
}

func (o *CliOptions) Complete(cmd *cobra.Command, args []string) error {
	// Complete Logger
	o.CompleteLogger()

	if o.UseKubeClient {
		if err := o.CompleteKubeClient(); err != nil {
			return err
		}
	} else {
		// Complete CliConfig
		if err := o.CompleteCliConfig(); err != nil {
			return err
		}

		// Complete Client
		if err := o.CompleteClient(); err != nil {
			return err
		}
	}

	return nil
}

func (o *CliOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// NewConnectRequestWithAuth wraps a connect.NewRequest with authorization header
func NewConnectRequestWithAuth[T any](o *CliConfig, message *T) *connect.Request[T] {
	req := connect.NewRequest(message)
	if len(o.Cookie) > 0 {
		req.Header().Add("Cookie", o.Cookie)
	}
	if len(o.Token) > 0 {
		req.Header().Add("Authorization", o.Token)
	}
	return req
}
