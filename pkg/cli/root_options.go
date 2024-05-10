package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

type RootOptions struct {
	UseKubeAPI     bool
	KubeConfigPath string
	KubeContext    string
	DashboardURL   string
	ConfigPath     string
	LogLevel       int

	In       io.Reader
	Out      io.Writer
	ErrOut   io.Writer
	Versions VersionInfo

	Ctx             context.Context
	Logr            *clog.Logger
	KosmoClient     *kosmo.Client
	CosmoDashClient *CosmoDashClient
	CliConfig       *Config
}

func NewRootOptions() *RootOptions {
	ctx := context.TODO()
	return &RootOptions{Ctx: ctx}
}

const (
	ENV_CONFIG        = "COSMOCTL_CONFIG"
	ENV_DASHBOARD_URL = "COSMOCTL_DASHBOARD_URL"
)

func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.SetIn(o.In)
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.PersistentFlags().BoolVarP(&o.UseKubeAPI,
		"kube", "k", false, "use kubernetes API client instead of cosmo dashboard API client")

	cmd.PersistentFlags().StringVar(&o.KubeConfigPath,
		"kubeconfig", "", "kubeconfig file path. env:KUBECONFIG (default: $HOME/.kube/config)")

	cmd.PersistentFlags().StringVar(&o.DashboardURL,
		"dashboard-url", "", "COSMO Dashboard server endpoint URL. env:COSMOCTL_DASHBOARD_URL")

	cmd.PersistentFlags().StringVar(&o.ConfigPath,
		"config", "", "cosmoctl config file path. env:COSMOCTL_CONFIG (default: $HOME/.config/cosmocfg)")

	cmd.PersistentFlags().StringVar(&o.KubeContext,
		"context", "", "kube-context (default: current context)")

	cmd.PersistentFlags().IntVarP(&o.LogLevel,
		"verbose", "v", 0, "log level. -1:DISABLED, 0:INFO, 1:DEBUG, 2:ALL")
}

func (o *RootOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *RootOptions) CompleteWithoutClient(cmd *cobra.Command, args []string) error {
	if err := o.buildLogger(); err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}
	return nil
}

func (o *RootOptions) Complete(cmd *cobra.Command, args []string) error {
	if err := o.buildLogger(); err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}
	if o.UseKubeAPI && o.KosmoClient == nil {
		if err := o.buildKosmoClient(); err != nil {
			return fmt.Errorf("failed to kubernetes client: %w", err)
		}
	} else {
		cfgPath, err := o.GetConfigFilePath()
		if err != nil {
			return fmt.Errorf("failed to get config file path: %w", err)
		}
		o.Logr.Debug().Info("config file path", "path", cfgPath, "dir", filepath.Dir(cfgPath))

		cfg, err := NewOrLoadConfigFile(cfgPath)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
		o.CliConfig = cfg

		if err := o.buildDashClient(); err != nil {
			return fmt.Errorf("failed to build COSMO Dashboard API client: %w", err)
		}
	}

	return nil
}

func (o *RootOptions) buildLogger() error {
	if o.LogLevel >= 0 {
		opt := zap.Options{
			Development: true,
			Level:       zapcore.Level(-o.LogLevel),
		}
		o.Logr = clog.NewLogger(zap.New(zap.UseFlagOptions(&opt)))
		o.Ctx = clog.IntoContext(o.Ctx, o.Logr)
	} else {
		o.Logr = clog.NewLogger(logr.Discard())
	}
	return nil
}

func (o *RootOptions) GetConfigFilePath() (string, error) {
	if o.ConfigPath != "" {
		return o.ConfigPath, nil
	} else if envCfg := os.Getenv(ENV_CONFIG); envCfg != "" {
		return envCfg, nil
	} else {
		d, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(d, ".config", "cosmocfg"), nil
	}
}

func (o *RootOptions) GetDashboardURL() string {
	if o.DashboardURL != "" {
		return o.DashboardURL
	} else if envURL := os.Getenv(ENV_DASHBOARD_URL); envURL != "" {
		return envURL
	} else if o.CliConfig.Endpoint != "" {
		return o.CliConfig.Endpoint
	} else {
		return ""
	}
}

func (o *RootOptions) buildDashClient() error {
	dashURL := o.GetDashboardURL()
	if dashURL == "" {
		return fmt.Errorf("failed to get dashboard URL. login first or run with --dashboard-url option")
	}
	o.Logr.Debug().Info("Dashboard URL", "url", dashURL)

	baseURL, err := url.Parse(dashURL)
	if err != nil {
		return err
	}

	o.CosmoDashClient = NewCosmoDashClient(http.DefaultClient, baseURL)
	return nil
}

func (o *RootOptions) buildKosmoClient() error {
	debug := o.Logr.WithCaller().DebugAll()

	cfgFlg := genericclioptions.NewConfigFlags(true)
	debug.Info("kubeconfigs", "kubeConfigPath", o.KubeConfigPath, "kubeContext", o.KubeContext)

	if o.KubeConfigPath != "" {
		cfgFlg.KubeConfig = &o.KubeConfigPath
	}
	if o.KubeContext != "" {
		cfgFlg.Context = &o.KubeContext
	}

	cfg, err := cfgFlg.ToRESTConfig()
	if err != nil {
		return err
	}
	debug.Info("RestConfig", "cfg", cfg)

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cosmov1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme

	baseclient, err := kosmo.NewClientByRestConfig(cfg, scheme)
	if err != nil {
		return err
	}
	o.KosmoClient = &baseclient

	return nil
}

func (o *RootOptions) Logger() *clog.Logger {
	return o.Logr
}

// GetCurrentWorkspaceName returns current workspace name.
// If running in Workspace pod, hostname is like `ws1-workspace-575db4c9cd-h558m`
// the first part is workspace name prefixed by cosmo.
func GetCurrentWorkspaceName() string {
	hostname := os.Getenv("HOSTNAME")
	h := strings.Split(hostname, "-")
	if len(h) > 3 && h[0] != "" {
		return h[0]
	}
	return ""
}
