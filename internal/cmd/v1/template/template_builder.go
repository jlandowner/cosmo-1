package template

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	cmdutil "github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	"github.com/cosmo-workspace/cosmo/pkg/kubeutil"
	"github.com/cosmo-workspace/cosmo/pkg/template"
	"github.com/cosmo-workspace/cosmo/pkg/workspace"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type TemplateObjectBuilder struct {
	tmpl              cosmov1alpha1.TemplateObject
	disableNamePrefix bool
}

func NewTemplateObjectBuilder(isClusterScope bool) *TemplateObjectBuilder {
	var b TemplateObjectBuilder

	if isClusterScope {
		b.tmpl = &cosmov1alpha1.ClusterTemplate{}
	} else {
		b.tmpl = &cosmov1alpha1.Template{}
	}

	scheme := runtime.NewScheme()
	if err := cosmov1alpha1.AddToScheme(scheme); err != nil {
		panic(err)
	}

	gvk, err := apiutil.GVKForObject(b.tmpl, scheme)
	if err != nil {
		panic(err)
	}
	b.tmpl.SetGroupVersionKind(gvk)
	return &b
}

func (b *TemplateObjectBuilder) Name(name string) *TemplateObjectBuilder {
	b.tmpl.SetName(name)
	return b
}

func (b *TemplateObjectBuilder) Description(desc string) *TemplateObjectBuilder {
	b.tmpl.GetSpec().Description = desc
	return b
}

func (b *TemplateObjectBuilder) RequiredVars(vars []string) *TemplateObjectBuilder {
	if len(vars) > 0 {
		vv := make([]cosmov1alpha1.RequiredVarSpec, 0, len(vars))
		for _, v := range vars {
			vcol := strings.Split(v, ":")
			varSpec := cosmov1alpha1.RequiredVarSpec{Var: vcol[0]}
			if len(vcol) > 1 {
				varSpec.Default = vcol[1]
			}
			vv = append(vv, varSpec)
		}
		b.tmpl.GetSpec().RequiredVars = vv
	}
	return b
}

func (b *TemplateObjectBuilder) SetUserRoles(roles []string) *TemplateObjectBuilder {
	if len(roles) > 0 {
		kubeutil.SetAnnotation(b.tmpl, cosmov1alpha1.TemplateAnnKeyUserRoles, strings.Join(roles, ","))
	}
	return b
}

func (b *TemplateObjectBuilder) SetRequiredAddons(addons []string) *TemplateObjectBuilder {
	if len(addons) > 0 {
		kubeutil.SetAnnotation(b.tmpl, cosmov1alpha1.TemplateAnnKeyRequiredAddons, strings.Join(addons, ","))
	}
	return b
}

func (b *TemplateObjectBuilder) DisableNamePrefix() *TemplateObjectBuilder {
	kubeutil.SetAnnotation(b.tmpl, cosmov1alpha1.TemplateAnnKeyDisableNamePrefix, strconv.FormatBool(true))
	b.disableNamePrefix = true
	return b
}

func (b *TemplateObjectBuilder) TypeUserAddon(setDefault bool) *TemplateObjectBuilder {
	template.SetTemplateType(b.tmpl, cosmov1alpha1.TemplateLabelEnumTypeUserAddon)
	if setDefault {
		kubeutil.SetAnnotation(b.tmpl, cosmov1alpha1.UserAddonTemplateAnnKeyDefaultUserAddon, strconv.FormatBool(true))
	}

	b.DisableNamePrefix()
	return b
}

func (b *TemplateObjectBuilder) TypeWorkspace(wscfg cosmov1alpha1.Config) *TemplateObjectBuilder {
	template.SetTemplateType(b.tmpl, cosmov1alpha1.TemplateLabelEnumTypeWorkspace)
	workspace.SetConfigOnTemplateAnnotations(b.tmpl, wscfg)
	return b
}

func (b *TemplateObjectBuilder) RawYAML(rawYAML string) *TemplateObjectBuilder {
	b.tmpl.GetSpec().RawYaml = rawYAML
	return b
}

func (b *TemplateObjectBuilder) KustomizeBuilder() *KustomizeBuilder {
	builder := NewKustomizeBuilder()
	if b.disableNamePrefix {
		builder.DisableNamePrefix()
	}
	return builder
}

func (b *TemplateObjectBuilder) Build() ([]byte, error) {
	return yaml.Marshal(b.tmpl)
}

func ExecKustomizeToRawYAML(ctx context.Context, builder *KustomizeBuilder, rawYAML string) (string, error) {
	const packagedFile = "packaged.yaml"

	// create tmp dir
	tmpDir, err := os.MkdirTemp(os.TempDir(), "cosmoctl-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir : %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// save it as "packaged" file
	if err := cmdutil.CreateFile(tmpDir, packagedFile, []byte(rawYAML)); err != nil {
		return "", err
	}

	builder.Resources(packagedFile)
	kust := builder.Build()

	out, err := cmdutil.ExecKustomize(ctx, tmpDir, kust)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type KustomizeBuilder struct {
	kust types.Kustomization
}

func NewKustomizeBuilder() *KustomizeBuilder {
	var b KustomizeBuilder
	b.kust = types.Kustomization{
		CommonLabels: map[string]string{
			cosmov1alpha1.LabelKeyInstanceName: template.DefaultVarsInstance,
			cosmov1alpha1.LabelKeyTemplateName: template.DefaultVarsTemplate,
		},
		Namespace:  template.DefaultVarsNamespace,
		Resources:  []string{},
		NamePrefix: template.DefaultVarsInstance + "-",
	}
	return &b
}

func (b *KustomizeBuilder) Resources(fileNames ...string) *KustomizeBuilder {
	b.kust.Resources = append(b.kust.Resources, fileNames...)
	return b
}

func (b *KustomizeBuilder) DisableNamePrefix() *KustomizeBuilder {
	b.kust.NamePrefix = ""
	return b
}

func (b *KustomizeBuilder) AddPatchesStrategicMerge(files ...types.PatchStrategicMerge) *KustomizeBuilder {
	if b.kust.PatchesStrategicMerge == nil {
		b.kust.PatchesStrategicMerge = files
	} else {
		b.kust.PatchesStrategicMerge = append(b.kust.PatchesStrategicMerge, files...)
	}
	return b
}

func (b *KustomizeBuilder) Build() *types.Kustomization {
	return &b.kust
}
