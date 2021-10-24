package transformer

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/core/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/template"
)

func TestMetadataTransformer_Transform(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := cosmov1alpha1.AddToScheme(scheme); err != nil {
		t.Errorf("Failed to AddToScheme %v", err)
	}
	type fields struct {
		inst   *cosmov1alpha1.Instance
		tmpl   *cosmov1alpha1.Template
		scheme *runtime.Scheme
	}
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "New label",
			fields: fields{
				inst: &cosmov1alpha1.Instance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cs1",
						Namespace: "cosmo-user-tom",
					},
					Spec: cosmov1alpha1.InstanceSpec{
						Template: cosmov1alpha1.TemplateRef{
							Name: "code-server",
						},
						Override: cosmov1alpha1.OverrideSpec{},
						Vars:     map[string]string{"{{TEST}}": "OK"},
					},
				},
				tmpl: &cosmov1alpha1.Template{
					ObjectMeta: metav1.ObjectMeta{
						Name: "code-server",
					},
					Spec: cosmov1alpha1.TemplateSpec{
						RawYaml: "data",
					},
				},
				scheme: scheme,
			},
			args: args{
				src: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cosmo/ingress-patch-enable: "true"
    kubernetes.io/ingress.class: alb
  name: test
  namespace: cosmo-user-tom
spec:
  host: example.com
`,
			},
			want: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cosmo/ingress-patch-enable: "true"
    kubernetes.io/ingress.class: alb
  labels:
    cosmo/instance: cs1
    cosmo/template: code-server
  name: cs1-test
  namespace: cosmo-user-tom
  ownerReferences:
  - apiVersion: cosmo.cosmo-workspace.github.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: Instance
    name: cs1
    uid: ""
spec:
  host: example.com
`,
			wantErr: false,
		},
		{
			name: "Append label",
			fields: fields{
				inst: &cosmov1alpha1.Instance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cs1",
						Namespace: "cosmo-user-tom",
					},
					Spec: cosmov1alpha1.InstanceSpec{
						Template: cosmov1alpha1.TemplateRef{
							Name: "code-server",
						},
						Override: cosmov1alpha1.OverrideSpec{},
						Vars:     map[string]string{"{{TEST}}": "OK"},
					},
				},
				tmpl: &cosmov1alpha1.Template{
					ObjectMeta: metav1.ObjectMeta{
						Name: "code-server",
					},
					Spec: cosmov1alpha1.TemplateSpec{
						RawYaml: "data",
					},
				},
				scheme: scheme,
			},
			args: args{
				src: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cosmo/ingress-patch-enable: "true"
    kubernetes.io/ingress.class: alb
  labels:
    key: val
  name: test
  namespace: cosmo-user-tom
spec:
  host: example.com
`,
			},
			want: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cosmo/ingress-patch-enable: "true"
    kubernetes.io/ingress.class: alb
  labels:
    cosmo/instance: cs1
    cosmo/template: code-server
    key: val
  name: cs1-test
  namespace: cosmo-user-tom
  ownerReferences:
  - apiVersion: cosmo.cosmo-workspace.github.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: Instance
    name: cs1
    uid: ""
spec:
  host: example.com
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewMetadataTransformer(tt.fields.inst, tt.fields.tmpl, tt.fields.scheme)
			_, obj, err := template.StringToUnstructured(tt.args.src)
			if err != nil {
				t.Errorf("MetadataTransformer.Transform() template.StringToUnstructured error = %v", err)
				return
			}
			gotObj, err := tr.Transform(obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetadataTransformer.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var want unstructured.Unstructured
				err := yaml.Unmarshal([]byte(tt.want), &want)
				if err != nil {
					t.Errorf("yaml.Marshal err = %v", err)
				}
				if !equality.Semantic.DeepEqual(gotObj, &want) {
					t.Errorf("MetadataTransformer.Transform() = %v, want %v", *gotObj, want)
				}

			} else {
				if gotObj != nil {
					t.Errorf("MetadataTransformer.Transform() gotObj not nil %v", gotObj)
				}
			}
		})
	}
}