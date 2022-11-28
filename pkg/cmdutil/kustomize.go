package cmdutil

import (
	"context"
	"fmt"
	"os/exec"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/types"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/template"
)

const (
	KustomizationFile   = "kustomization.yaml"
	defaultPackagedFile = "packaged.yaml"
)

type Kustomize struct {
	kust *types.Kustomization

	dir          string
	packagedFile string
}

func NewKustomize(buildDir string) *Kustomize {
	k := &Kustomize{
		dir:          buildDir,
		packagedFile: defaultPackagedFile,
		kust: &types.Kustomization{
			CommonLabels: map[string]string{
				cosmov1alpha1.LabelKeyInstanceName: template.DefaultVarsInstance,
				cosmov1alpha1.LabelKeyTemplateName: template.DefaultVarsTemplate,
			},
			Namespace: template.DefaultVarsNamespace,
			Resources: []string{
				defaultPackagedFile,
			},
			NamePrefix: template.DefaultVarsInstance + "-",
		},
	}
	return k
}

func (k *Kustomize) WithDisableNamePrefix() {
	k.kust.NamePrefix = ""
}

func (k *Kustomize) WithPatchesStrategicMerges(files ...types.PatchStrategicMerge) {
	if len(k.kust.PatchesStrategicMerge) == 0 {
		k.kust.PatchesStrategicMerge = files
	} else {
		k.kust.PatchesStrategicMerge = append(k.kust.PatchesStrategicMerge, files...)
	}
}

func (k *Kustomize) WithResources(files ...string) {
	if len(k.kust.Resources) == 0 {
		k.kust.Resources = files
	} else {
		k.kust.Resources = append(k.kust.Resources, files...)
	}
}

func (k *Kustomize) WritePackagedFile(resourcesYAML []byte) ([]unstructured.Unstructured, error) {
	dummies, err := dummyTemplateBuild(string(resourcesYAML))
	if err != nil {
		return nil, fmt.Errorf("failed to pre-build template: %w", err)
	}

	err = CreateFile(k.dir, k.packagedFile, resourcesYAML)
	if err != nil {
		return nil, fmt.Errorf("failed to create packaged file: %w", err)
	}

	return dummies, nil
}

func dummyTemplateBuild(rawTmpl string) ([]unstructured.Unstructured, error) {
	var inst cosmov1alpha1.Instance
	inst.SetName("dummy")
	inst.SetNamespace("dummy")

	builder := template.NewRawYAMLBuilder(rawTmpl, &inst)
	return builder.Build()
}

func (k *Kustomize) Build(ctx context.Context) ([]byte, error) {
	log := clog.FromContext(ctx).WithCaller()

	cmd, err := kustomizeBuildCmd()
	if err != nil {
		return nil, err
	}
	log.Debug().Info("kustomize cmd", "cmd", cmd)

	kustYaml, err := yaml.Marshal(k.kust)
	if err != nil {
		return nil, err
	}
	log.Debug().Info(string(kustYaml), "obj", KustomizationFile)

	// create kustomization.yaml
	if err := CreateFile(k.dir, KustomizationFile, kustYaml); err != nil {
		return nil, err
	}
	defer RemoveFile(k.dir, KustomizationFile)

	// run kustomize build
	cmd = append(cmd, k.dir)

	out, err := exec.CommandContext(ctx, cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to exec kustomize : %w : %s", err, out)
	}
	return out, nil
}

func kustomizeBuildCmd() ([]string, error) {
	kust, kustErr := exec.LookPath("kustomize")
	if kustErr != nil {
		kctl, kctlErr := exec.LookPath("kubectl")
		if kctlErr != nil {
			return nil, fmt.Errorf("kubectl nor kustomize found: kustmizr=%v, kubectl=%v", kustErr, kctlErr)
		}
		return []string{kctl, "kustomize"}, nil
	}
	return []string{kust, "build"}, nil
}
