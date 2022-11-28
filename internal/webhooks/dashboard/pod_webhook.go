/*
MIT License

Copyright (c) 2022 cosmo-workspace
*/

package dashboard

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
)

type PodWebhook struct {
	// Logger
	Log clog.Logger
	// Kube client
	Client client.Client
	// RootCASecretKey is a key of certficate secret key. Secert type must be "kubernetes.io/tls"
	RootCASecretKey types.NamespacedName
	// RootCASecretDataKey is a key of root ca data key. if not set, uses "ca.crt"
	RootCASecretDataKey string
	// Hostname is a hostname for API Endpoint
	Hostname string
	// Port is a port for API Endpoint
	Port int

	endpoint string
	decoder  *admission.Decoder
}

func (r *PodWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	if r.RootCASecretDataKey == "" {
		r.RootCASecretDataKey = "ca.crt"
	}
	if _, err := r.GetCA(context.TODO()); err != nil {
		return fmt.Errorf("failed to get cert: %w", err)
	}

	r.endpoint = (&url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", r.Hostname, r.Port),
	}).String()

	mgr.GetWebhookServer().Register(
		"/mutate-core-v1-pod",
		&webhook.Admission{Handler: r},
	)
	return nil
}

//+kubebuilder:webhook:path=/mutate-core-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups=core,resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io,admissionReviewVersions=v1

func (r *PodWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := r.Log.WithValues("UID", req.UID, "GroupVersionKind", req.Kind.String(), "Name", req.Name, "Namespace", req.Namespace)
	ctx = clog.IntoContext(ctx, log)
	var pod corev1.Pod

	if err := r.decoder.Decode(req, &pod); err != nil {
		log.Error(err, "failed to decode request")
		return admission.Errored(http.StatusBadRequest, err)
	}
	log.DumpObject(r.Client.Scheme(), &pod, "request pod")

	if err := r.Default(ctx, &pod); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	return admission.Allowed("mutated")
}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PodWebhook) Default(ctx context.Context, pod *corev1.Pod) error {
	ca, err := r.GetCA(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ca cert: %w", err)
	}
	b64ca := base64.RawStdEncoding.EncodeToString(ca)

	for i, cont := range pod.Spec.Containers {
		if cont.Env == nil {
			cont.Env = make([]corev1.EnvVar, 0)
		}
		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env,
			corev1.EnvVar{Name: cosmov1alpha1.EnvServerCA, Value: b64ca})

		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env,
			corev1.EnvVar{Name: cosmov1alpha1.EnvServerEndpoint, Value: r.endpoint})
	}

	return nil
}

func (r *PodWebhook) GetCA(ctx context.Context) ([]byte, error) {
	var secret corev1.Secret
	if err := r.Client.Get(ctx, r.RootCASecretKey, &secret); err != nil {
		return nil, fmt.Errorf("failed to get cert secret: %w", err)
	}

	ca := secret.Data[r.RootCASecretDataKey]
	if ca == nil {
		return nil, fmt.Errorf("ca data key %s is not found", r.RootCASecretDataKey)
	}
	if ok := x509.NewCertPool().AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("invalid ca certificate format")
	}
	return ca, nil
}
