package cmdutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func GetKubeConfig(path string) (*api.Config, error) {
	if path == "" {
		rule := clientcmd.NewDefaultClientConfigLoadingRules()
		return rule.Load()
	} else {
		return clientcmd.LoadFromFile(path)
	}
}

var inclusterNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func GetDefaultNamespace(cfg *api.Config, kubecontext string) string {
	if cfg == nil || len(cfg.Contexts) == 0 {
		b, _ := ioutil.ReadFile(inclusterNamespaceFile)
		if len(b) != 0 {
			return string(b)
		}
		return ""
	}
	var ctxName string
	if kubecontext == "" {
		ctxName = cfg.CurrentContext
	} else {
		ctxName = kubecontext
	}
	ctx, ok := cfg.Contexts[ctxName]
	if !ok {
		return ""
	}
	return ctx.Namespace
}

func CreateFile(dir, fname string, data []byte) error {
	fullPath, err := filepath.Abs(dir + "/" + fname)
	if err != nil {
		return fmt.Errorf("invaid file path : %w", err)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create %s : %w", fname, err)
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return fmt.Errorf("failed to create %s : %w", fname, err)
	}
	return nil
}

func RemoveFile(dir, fname string) error {
	fullPath, err := filepath.Abs(dir + "/" + fname)
	if err != nil {
		return fmt.Errorf("invaid file path : %w", err)
	}
	return os.Remove(fullPath)
}

func PrintfColorErr(out io.Writer, msg string, a ...interface{}) {
	fmt.Fprintf(out, "\x1b[33m%s\x1b[0m", fmt.Sprintf(msg, a...))
}

func PrintfColorInfo(out io.Writer, msg string, a ...interface{}) {
	fmt.Fprintf(out, "\x1b[32m%s\x1b[0m", fmt.Sprintf(msg, a...))
}
