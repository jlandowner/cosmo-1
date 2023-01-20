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
	"io/ioutil"
	"net/http"
	"net/url"

	corev1 "k8s.io/api/core/v1"
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
	// CACertFile file path
	CACertFile string
	// Hostname is a hostname for API Endpoint
	Hostname string
	// Port is a port for API Endpoint
	Port int

	endpoint string
	b64ca    string
	decoder  *admission.Decoder
}

func (r *PodWebhook) SetupCA() error {
	if r.CACertFile == "" {
		r.CACertFile = "ca.crt"
	}
	ca, err := ioutil.ReadFile(r.CACertFile)
	if err != nil {
		return err
	}

	if ok := x509.NewCertPool().AppendCertsFromPEM(ca); !ok {
		return errors.New("invalid ca certificate format")
	}
	r.b64ca = base64.RawStdEncoding.EncodeToString(ca)

	return nil
}

func (r *PodWebhook) SetupEndpoint() {
	r.endpoint = (&url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", r.Hostname, r.Port),
	}).String()
}

func (r *PodWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	r.SetupEndpoint()
	r.SetupCA()
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
	for i, cont := range pod.Spec.Containers {
		if cont.Env == nil {
			cont.Env = make([]corev1.EnvVar, 0)
		}
		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env,
			corev1.EnvVar{Name: cosmov1alpha1.EnvServerCA, Value: r.b64ca})

		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env,
			corev1.EnvVar{Name: cosmov1alpha1.EnvServerEndpoint, Value: r.endpoint})
	}

	return nil
}
