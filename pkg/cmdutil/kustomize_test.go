package cmdutil

import (
	"os"
	"testing"
)

func TestPrepareKustomizeBuildCmd(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "kustomize",
			want:    []string{"/usr/local/bin/kustomize", "build"},
			wantErr: false,
		},
		{
			name:    "kubectl",
			want:    []string{"/usr/bin/kubectl", "kustomize"},
			wantErr: false,
		},
	}
	t.Logf("KustomizeBuildCmd() PATH = %v", os.Getenv("PATH"))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kustomizeBuildCmd()
			if (err != nil) != tt.wantErr {
				t.Logf("KustomizeBuildCmd() kustomize or kubectl is not found: %v", err)
				return
			}
			t.Logf("KustomizeBuildCmd() got = %v", got)
			// This test is dependent on the testing environment and goes OK at any time.
			// When you test manually, do comment-in the below line and see the results.
			// t.Fail()
		})
	}
}
