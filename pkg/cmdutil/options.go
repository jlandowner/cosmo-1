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

	"github.com/bufbuild/connect-go"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
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

func (o *CliOptions) Complete(cmd *cobra.Command, args []string) error {
	// Complete Logger
	o.CompleteLogger()

	// Complete CliConfig
	// First, load config from file
	// If failed to load config file, then load incluster config
	// If both failed, return err
	if err := o.CliConfig.CompleteByConfigFile(o.CliConfigFilePath); err != nil {
		err = fmt.Errorf("failed to load config file %s: %w", o.CliConfigFilePath, err)

		if ierr := o.CliConfig.CompleteInCluster(); ierr != nil {
			return fmt.Errorf("failed to load config %v: %w", ierr, err)
		}
	}

	// Complete TLS Client
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(o.ServerCA)

	o.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caPool,
			},
		},
	}

	return nil
}

func (o *CliOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// NewConnectRequestWithAuth wraps a connect.NewRequest with authorization header
func NewConnectRequestWithAuth[T any](token string, message *T) *connect.Request[T] {
	req := connect.NewRequest(message)
	req.Header().Add("Authorization", token)
	return req
}
