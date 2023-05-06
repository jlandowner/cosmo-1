package transformer

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/template"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestNetworkTransformer_Transform(t *testing.T) {
	prefix := netv1.PathTypePrefix
	type fields struct {
		netSpec  *cosmov1alpha1.NetworkOverrideSpec
		instName string
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
			name: "Override service",
			fields: fields{
				netSpec: &cosmov1alpha1.NetworkOverrideSpec{
					Ingress: []cosmov1alpha1.IngressOverrideSpec{
						{
							TargetName: "workspace",
							Annotations: map[string]string{
								"key": "val",
							},
							Rules: []netv1.IngressRule{
								{
									Host: "example2.com",
									IngressRuleValue: netv1.IngressRuleValue{
										HTTP: &netv1.HTTPIngressRuleValue{
											Paths: []netv1.HTTPIngressPath{
												{
													Backend: netv1.IngressBackend{
														Service: &netv1.IngressServiceBackend{
															Name: "test2",
															Port: netv1.ServiceBackendPort{
																Number: 8081,
															},
														},
													},
													Path:     "/test",
													PathType: &prefix,
												},
											},
										},
									},
								},
							},
						},
					},
					Service: []cosmov1alpha1.ServiceOverrideSpec{
						{
							TargetName: "workspace",
							Ports: []corev1.ServicePort{
								{
									Name:       "port2",
									Port:       int32(8081),
									Protocol:   "TCP",
									TargetPort: intstr.FromInt(8081),
								},
							},
						},
					},
				},
				instName: "instance",
			},
			args: args{
				src: `apiVersion: v1
kind: Service
metadata:
  name: instance-workspace
  namespace: default
spec:
  ports:
  - name: port1
    port: 8080
    protocol: TCP
  type: ClusterIP
`,
			},
			want: `apiVersion: v1
kind: Service
metadata:
  name: instance-workspace
  namespace: default
spec:
  ports:
  - name: port1
    port: 8080
    protocol: TCP
  - name: port2
    port: 8081
    protocol: TCP
    targetPort: 8081
  type: ClusterIP
`,
			wantErr: false,
		},
		{
			name: "Override ingress",
			fields: fields{
				netSpec: &cosmov1alpha1.NetworkOverrideSpec{
					Ingress: []cosmov1alpha1.IngressOverrideSpec{
						{
							TargetName: "workspace",
							Annotations: map[string]string{
								"key": "val",
							},
							Rules: []netv1.IngressRule{
								{
									Host: "example2.com",
									IngressRuleValue: netv1.IngressRuleValue{
										HTTP: &netv1.HTTPIngressRuleValue{
											Paths: []netv1.HTTPIngressPath{
												{
													Backend: netv1.IngressBackend{
														Service: &netv1.IngressServiceBackend{
															Name: "test2",
															Port: netv1.ServiceBackendPort{
																Number: 8081,
															},
														},
													},
													Path:     "/test",
													PathType: &prefix,
												},
											},
										},
									},
								},
							},
						},
					},
					Service: []cosmov1alpha1.ServiceOverrideSpec{
						{
							TargetName: "workspace",
							Ports: []corev1.ServicePort{
								{
									Name:       "port2",
									Port:       int32(8081),
									Protocol:   "TCP",
									TargetPort: intstr.FromInt(8081),
								},
							},
						},
					},
				},
				instName: "instance",
			},
			args: args{
				src: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: instance-workspace
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port:
              number: 8080
        path: /*
        pathType: Exact
`,
			},
			want: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    key: val
  name: instance-workspace
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port:
              number: 8080
        path: /*
        pathType: Exact
  - host: example2.com
    http:
      paths:
      - backend:
          service:
            name: test2
            port:
              number: 8081
        path: /test
        pathType: Prefix
`,
			wantErr: false,
		},
		{
			name: "netSpec nil",
			fields: fields{
				netSpec:  nil,
				instName: "instance",
			},
			args: args{
				src: `apiVersion: v1
kind: Service
metadata:
  name: instance-test
  namespace: default
spec:
  ports:
  - name: port1
    port: 8080
    protocol: TCP
  - name: port2
    port: 8081
    protocol: TCP
  type: ClusterIP
`,
			},
			want: `apiVersion: v1
kind: Service
metadata:
  name: instance-test
  namespace: default
spec:
  ports:
  - name: port1
    port: 8080
    protocol: TCP
  - name: port2
    port: 8081
    protocol: TCP
  type: ClusterIP
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewNetworkTransformer(tt.fields.netSpec, tt.fields.instName)
			_, obj, err := template.StringToUnstructured(tt.args.src)
			if err != nil {
				t.Errorf("NetworkTransformer.Transform() template.StringToUnstructured error = %v", err)
				return
			}
			gotObj, err := tr.Transform(obj)
			snaps.MatchSnapshot(t, err)
			snaps.MatchJSON(t, gotObj)
		})
	}
}

func Test_overrideAnnotations(t *testing.T) {
	type args struct {
		obj string
		ann map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    key1: val
  name: test-ing
  namespace: default
spec:
  rules:
  # rule[0] is WebIDE main page
  # host must not be set
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port:
              number: 8080 # MUST BE NUMBER
        path: /*
        pathType: Exact`,
				ann: map[string]string{
					"key2": "VAL1", "KEY": "VAL2",
				},
			},
		},
		{
			name: "Empty annotation",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ing
  namespace: default
spec:
  rules:
  # rule[0] is WebIDE main page
  # host must not be set
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port:
              number: 8080 # MUST BE NUMBER
        path: /*
        pathType: Exact`,
				ann: map[string]string{
					"key": "val",
				},
			},
		},
		{
			name: "nil",
			args: args{
				obj: ``,
				ann: map[string]string{
					"key2": "VAL1", "KEY": "VAL2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj *unstructured.Unstructured
			var err error
			if tt.args.obj != "" {
				_, obj, err = template.StringToUnstructured(tt.args.obj)
				if err != nil {
					t.Errorf("overrideAnnotations() template.StringToUnstructured error = %v", err)
					return
				}
			} else {
				obj = nil
			}
			overrideAnnotations(obj, tt.args.ann)
			snaps.MatchJSON(t, obj)
		})
	}
}

func Test_overrideIngressRules(t *testing.T) {
	prefix := netv1.PathTypePrefix
	type args struct {
		obj      string
		ingRules []netv1.IngressRule
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "append",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ing
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port: 
              number: 8080 # MUST BE NUMBER
        path: /*
        pathType: Exact`,
				ingRules: []netv1.IngressRule{
					{
						Host: "example2.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test2",
												Port: netv1.ServiceBackendPort{
													Number: 8081,
												},
											},
										},
										Path:     "/test",
										PathType: &prefix,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "override",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ing
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port: 
              number: 8080
        path: /
        pathType: Prefix`,
				ingRules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8081,
												},
											},
										},
										Path:     "/",
										PathType: &prefix,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "path",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ing
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port: 
              number: 8080
        path: /
        pathType: Prefix`,
				ingRules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8081,
												},
											},
										},
										Path:     "/",
										PathType: &prefix,
									},
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8082,
												},
											},
										},
										Path:     "/test",
										PathType: &prefix,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "rule and path",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ing
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port: 
              number: 8080
        path: /
        pathType: Prefix`,
				ingRules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8081,
												},
											},
										},
										Path:     "/",
										PathType: &prefix,
									},
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8082,
												},
											},
										},
										Path:     "/test",
										PathType: &prefix,
									},
								},
							},
						},
					},
					{
						Host: "example2.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8081,
												},
											},
										},
										Path:     "/",
										PathType: &prefix,
									},
								},
							},
						},
					},
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test",
												Port: netv1.ServiceBackendPort{
													Number: 8083,
												},
											},
										},
										Path:     "/add",
										PathType: &prefix,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "No rules",
			args: args{
				obj: `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ing
  namespace: default
spec:
  rules:
  # rule[0] is WebIDE main page
  # host must not be set
  - host: example.com
    http:
      paths:
      - backend:
          service:
            name: test
            port: 
              number: 8080 # MUST BE NUMBER
        path: /*
        pathType: Exact`,
				ingRules: []netv1.IngressRule{},
			},
		},
		{
			name: "nil",
			args: args{
				obj: ``,
				ingRules: []netv1.IngressRule{
					{
						Host: "example2.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test2",
												Port: netv1.ServiceBackendPort{
													Number: 8081,
												},
											},
										},
										Path: "/test",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj *unstructured.Unstructured
			var err error
			if tt.args.obj != "" {
				_, obj, err = template.StringToUnstructured(tt.args.obj)
				if err != nil {
					t.Errorf("overrideIngressRules() template.StringToUnstructured error = %v", err)
					return
				}
			} else {
				obj = nil
			}
			overrideIngressRules(obj, tt.args.ingRules)
			snaps.MatchJSON(t, obj)
		})
	}
}

func Test_overrideServicePort(t *testing.T) {
	type args struct {
		obj      string
		svcPorts []corev1.ServicePort
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "append",
			args: args{
				obj: `apiVersion: v1
kind: Service
metadata:
  name: test-svc
  namespace: default
spec:
  type: ClusterIP
  ports:
  # rule[0] is WebIDE main page
  - name: port1
    port: 8080
    protocol: TCP`,
				svcPorts: []corev1.ServicePort{
					{
						Name:       "port2",
						Port:       int32(8081),
						Protocol:   "TCP",
						TargetPort: intstr.FromInt(8081),
					},
					{
						Name:       "port3",
						Port:       int32(8082),
						Protocol:   "UDP",
						TargetPort: intstr.FromInt(8082),
					},
				},
			},
		},
		{
			name: "override",
			args: args{
				obj: `apiVersion: v1
kind: Service
metadata:
  name: test-svc
  namespace: default
spec:
  type: ClusterIP
  ports:
  # rule[0] is WebIDE main page
  - name: port1
    port: 8080
    protocol: TCP`,
				svcPorts: []corev1.ServicePort{
					{
						Name:       "port2",
						Port:       int32(8081),
						Protocol:   "TCP",
						TargetPort: intstr.FromInt(8081),
					},
					{
						Name:       "port1",
						Port:       int32(8082),
						Protocol:   "UDP",
						TargetPort: intstr.FromInt(8082),
					},
				},
			},
		},
		{
			name: "No ports",
			args: args{
				obj: `apiVersion: v1
kind: Service
metadata:
  name: test-svc
  namespace: default
spec:
  type: ClusterIP
  ports:
  # rule[0] is WebIDE main page
  - name: port1
    port: 8080
    protocol: TCP`,
				svcPorts: []corev1.ServicePort{},
			},
		},
		{
			name: "nil",
			args: args{
				obj: ``,
				svcPorts: []corev1.ServicePort{
					{
						Name:       "port2",
						Port:       int32(8081),
						Protocol:   "TCP",
						TargetPort: intstr.FromInt(8081),
					},
					{
						Name:       "port3",
						Port:       int32(8082),
						Protocol:   "UDP",
						TargetPort: intstr.FromInt(8082),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj *unstructured.Unstructured
			var err error
			if tt.args.obj != "" {
				_, obj, err = template.StringToUnstructured(tt.args.obj)
				if err != nil {
					t.Errorf("overrideServicePort() template.StringToUnstructured error = %v", err)
					return
				}
			} else {
				obj = nil
			}
			overrideServicePort(obj, tt.args.svcPorts)
			snaps.MatchJSON(t, obj)
		})
	}
}
