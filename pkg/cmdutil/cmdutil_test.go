package cmdutil

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd/api"
)

func TestGetDefaultNamespace(t *testing.T) {
	inclusterNamespaceFile = "incluster-namespace-test"
	CreateFile(".", inclusterNamespaceFile, []byte("incluster-ns"))
	defer RemoveFile(".", inclusterNamespaceFile)

	type args struct {
		cfg         *api.Config
		kubecontext string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "incluster",
			args: args{
				cfg:         nil,
				kubecontext: "default",
			},
			want: "incluster-ns",
		},
		{
			name: "kubeconfig",
			args: args{
				cfg: &api.Config{
					Contexts: map[string]*api.Context{
						"foo-cluster": {
							Namespace: "cosmo-user-foo",
						},
						"bar-cluster": {
							Namespace: "bar",
						},
					},
					CurrentContext: "bar-cluster",
				},
				kubecontext: "foo-cluster",
			},
			want: "cosmo-user-foo",
		},
		{
			name: "kubecontext not found in config",
			args: args{
				cfg: &api.Config{
					Contexts: map[string]*api.Context{
						"foo-cluster": {
							Namespace: "cosmo-user-foo",
						},
						"bar-cluster": {
							Namespace: "bar",
						},
					},
					CurrentContext: "bar-cluster",
				},
				kubecontext: "notfound",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDefaultNamespace(tt.args.cfg, tt.args.kubecontext)
			if got != tt.want {
				t.Errorf("GetDefaultNamespace() got = %v, want %v", got, tt.want)
			}
		})
	}
}
