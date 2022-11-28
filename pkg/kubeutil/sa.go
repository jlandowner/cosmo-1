package kubeutil

import (
	"context"
	"fmt"
	"net"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func AssumeServiceAccont(ctx context.Context, c client.Client, saKey types.NamespacedName) (client.Client, error) {
	var sa corev1.ServiceAccount
	err := c.Get(ctx, saKey, &sa)
	if err != nil {
		return nil, err
	}

	var secret *corev1.Secret
	for _, sec := range sa.Secrets {
		err = c.Get(ctx, types.NamespacedName{Name: sec.Name, Namespace: sa.Namespace}, secret)
		if err != nil {
			err = fmt.Errorf("failed to get secret %s/%s: %w", sec.Namespace, sec.Name, err)
			continue
		}
		break
	}
	if secret == nil && err != nil {
		return nil, err
	}

	token := secret.Data[corev1.ServiceAccountTokenKey]
	ca := secret.Data[corev1.ServiceAccountRootCAKey]
	// namespace := secret.Data[corev1.ServiceAccountNamespaceKey]

	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, fmt.Errorf("failed to get kubeapiserver host or port")
	}

	cfg := &rest.Config{
		Host: "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: ca,
		},
		BearerToken: string(token),
	}

	assumedClient, err := client.New(cfg, client.Options{
		Scheme: c.Scheme(),
		Mapper: c.RESTMapper(),
		Opts:   client.WarningHandlerOptions{},
	})

	return assumedClient, err
}
