package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/internal/dashboard"
	webhooks "github.com/cosmo-workspace/cosmo/internal/webhooks/dashboard"
	"github.com/cosmo-workspace/cosmo/pkg/auth"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/kosmo"
	"github.com/gorilla/securecookie"
)

var (
	o        = &options{}
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cosmov1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

type options struct {
	SessionAuthKey          string `printOption:"false"`
	StaticFileDir           string
	ResponseTimeoutSeconds  int64
	GracefulShutdownSeconds int64
	Insecure                bool
	ServerPort              int
	MaxAgeMinutes           int
	WebhookPort             int
	CertDir                 string
}

func main() {
	flag.StringVar(&o.SessionAuthKey, "auth-key", "", "Session authentication key. It must be 32 bytes secret string")
	flag.Int64Var(&o.ResponseTimeoutSeconds, "timeout-seconds", 3, "Timeout seconds for response")
	flag.Int64Var(&o.GracefulShutdownSeconds, "graceful-shutdown-seconds", 10, "Graceful shutdown seconds")
	flag.StringVar(&o.StaticFileDir, "serve-dir", "/app/public", "Static file dir to serve")
	flag.BoolVar(&o.Insecure, "insecure", false, "start http server not https server")
	flag.IntVar(&o.ServerPort, "port", 8443, "Port for dashboard server")
	flag.IntVar(&o.WebhookPort, "webhook-port", 9443, "Port for webhook server")
	flag.IntVar(&o.MaxAgeMinutes, "maxage-minutes", 720, "session maxage minutes")
	flag.StringVar(&o.CertDir, "cert-dir", "/tmp/k8s-webhook-server/serving-certs", "cert directory which has tls.key, tls.crt, ca.crt")

	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	printVersion()
	printOptions()

	ctx := ctrl.SetupSignalHandler()

	var authKey []byte
	if o.SessionAuthKey == "" {
		authKey = securecookie.GenerateRandomKey(32)
	} else {
		authKey = []byte(o.SessionAuthKey)
	}
	if len(authKey) != 32 {
		panic(fmt.Sprintf("auth key must be 32 bytes but %d", len(authKey)))
	}

	// Setup controller manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
		Port:               o.WebhookPort,
		LeaderElection:     false,
		CertDir:            o.CertDir,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if !o.Insecure {
		if err = (&webhooks.PodWebhook{
			RootCASecretKey: types.NamespacedName{Name: "cosmo-dashboard-cert", Namespace: "cosmo-system"},
		}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create pod webhook", "webhook", "Pod")
			os.Exit(1)
		}
	} else {
		setupLog.Info("WARNING: webhook is disabled when insecure")
	}

	// Setup server
	klient := kosmo.NewClient(mgr.GetClient())

	auths := make(map[cosmov1alpha1.UserAuthType]auth.Authorizer)
	auths[cosmov1alpha1.UserAuthTypePasswordSecert] = auth.NewPasswordSecretAuthorizer(klient)

	serv := (&dashboard.Server{
		Log:                 clog.NewLogger(ctrl.Log.WithName("dashboard")),
		Klient:              klient,
		GracefulShutdownDur: time.Second * time.Duration(o.GracefulShutdownSeconds),
		ResponseTimeout:     time.Second * time.Duration(o.ResponseTimeoutSeconds),
		StaticFileDir:       o.StaticFileDir,
		Port:                o.ServerPort,
		MaxAgeSeconds:       60 * o.MaxAgeMinutes,
		SessionAuthKey:      authKey,
		SessionName:         "cosmo-dashboard",
		TLSPrivateKeyPath:   filepath.Join(o.CertDir, "tls.key"),
		TLSCertPath:         filepath.Join(o.CertDir, "tls.crt"),
		Insecure:            o.Insecure,
		Authorizers:         auths,
	})
	if err := mgr.Add(serv); err != nil {
		setupLog.Error(err, "failed to add server to controller-manager")
		os.Exit(1)
	}

	// Start server
	setupLog.Info("Start controller manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func printOptions() {
	rv := reflect.ValueOf(*o)
	rt := rv.Type()
	options := make([]interface{}, rt.NumField()*2)

	for i := 0; i < rt.NumField(); i++ {
		options[i*2] = rt.Field(i).Name
		options[i*2+1] = rv.Field(i).Interface()

		if tag := rt.Field(i).Tag.Get("printOption"); tag != "" {
			if print, _ := strconv.ParseBool(tag); !print {
				options[i*2+1] = "*****"
			}
		}
	}

	setupLog.Info("options", options...)
}

func printVersion() {
	fmt.Println("cosmo-dashboard - cosmo v0.7.0 cosmo-workspace 2022")
}
