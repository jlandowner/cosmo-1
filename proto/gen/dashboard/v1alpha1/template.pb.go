//
//Cosmo Dashboard API
//Manipulate cosmo dashboard resource API

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: dashboard/v1alpha1/template.proto

package dashboardv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TemplateRequiredVars struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VarName      string `protobuf:"bytes,1,opt,name=var_name,json=varName,proto3" json:"var_name,omitempty"`
	DefaultValue string `protobuf:"bytes,2,opt,name=default_value,json=defaultValue,proto3" json:"default_value,omitempty"`
}

func (x *TemplateRequiredVars) Reset() {
	*x = TemplateRequiredVars{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dashboard_v1alpha1_template_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateRequiredVars) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateRequiredVars) ProtoMessage() {}

func (x *TemplateRequiredVars) ProtoReflect() protoreflect.Message {
	mi := &file_dashboard_v1alpha1_template_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateRequiredVars.ProtoReflect.Descriptor instead.
func (*TemplateRequiredVars) Descriptor() ([]byte, []int) {
	return file_dashboard_v1alpha1_template_proto_rawDescGZIP(), []int{0}
}

func (x *TemplateRequiredVars) GetVarName() string {
	if x != nil {
		return x.VarName
	}
	return ""
}

func (x *TemplateRequiredVars) GetDefaultValue() string {
	if x != nil {
		return x.DefaultValue
	}
	return ""
}

type Template struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name               string                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description        string                  `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	RequiredVars       []*TemplateRequiredVars `protobuf:"bytes,3,rep,name=required_vars,json=requiredVars,proto3" json:"required_vars,omitempty"`
	IsDefaultUserAddon *bool                   `protobuf:"varint,4,opt,name=is_default_user_addon,json=isDefaultUserAddon,proto3,oneof" json:"is_default_user_addon,omitempty"`
	IsClusterScope     bool                    `protobuf:"varint,5,opt,name=is_cluster_scope,json=isClusterScope,proto3" json:"is_cluster_scope,omitempty"`
	RequiredUseraddons []string                `protobuf:"bytes,6,rep,name=required_useraddons,json=requiredUseraddons,proto3" json:"required_useraddons,omitempty"`
	Userroles          []string                `protobuf:"bytes,7,rep,name=userroles,proto3" json:"userroles,omitempty"`
	Raw                *string                 `protobuf:"bytes,8,opt,name=raw,proto3,oneof" json:"raw,omitempty"`
}

func (x *Template) Reset() {
	*x = Template{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dashboard_v1alpha1_template_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Template) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Template) ProtoMessage() {}

func (x *Template) ProtoReflect() protoreflect.Message {
	mi := &file_dashboard_v1alpha1_template_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Template.ProtoReflect.Descriptor instead.
func (*Template) Descriptor() ([]byte, []int) {
	return file_dashboard_v1alpha1_template_proto_rawDescGZIP(), []int{1}
}

func (x *Template) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Template) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Template) GetRequiredVars() []*TemplateRequiredVars {
	if x != nil {
		return x.RequiredVars
	}
	return nil
}

func (x *Template) GetIsDefaultUserAddon() bool {
	if x != nil && x.IsDefaultUserAddon != nil {
		return *x.IsDefaultUserAddon
	}
	return false
}

func (x *Template) GetIsClusterScope() bool {
	if x != nil {
		return x.IsClusterScope
	}
	return false
}

func (x *Template) GetRequiredUseraddons() []string {
	if x != nil {
		return x.RequiredUseraddons
	}
	return nil
}

func (x *Template) GetUserroles() []string {
	if x != nil {
		return x.Userroles
	}
	return nil
}

func (x *Template) GetRaw() string {
	if x != nil && x.Raw != nil {
		return *x.Raw
	}
	return ""
}

var File_dashboard_v1alpha1_template_proto protoreflect.FileDescriptor

var file_dashboard_v1alpha1_template_proto_rawDesc = []byte{
	0x0a, 0x21, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x12, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x22, 0x56, 0x0a, 0x14, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x56, 0x61, 0x72, 0x73, 0x12,
	0x19, 0x0a, 0x08, 0x76, 0x61, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x61, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0xf9, 0x02, 0x0a, 0x08, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x4d, 0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x5f, 0x76,
	0x61, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x64, 0x61, 0x73, 0x68,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x54,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x56,
	0x61, 0x72, 0x73, 0x52, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x56, 0x61, 0x72,
	0x73, 0x12, 0x36, 0x0a, 0x15, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x48, 0x00, 0x52, 0x12, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x41, 0x64, 0x64, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x28, 0x0a, 0x10, 0x69, 0x73, 0x5f,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x73, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x63,
	0x6f, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x13, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x61, 0x64, 0x64, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x12, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x61, 0x64,
	0x64, 0x6f, 0x6e, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65,
	0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c,
	0x65, 0x73, 0x12, 0x15, 0x0a, 0x03, 0x72, 0x61, 0x77, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x03, 0x72, 0x61, 0x77, 0x88, 0x01, 0x01, 0x42, 0x18, 0x0a, 0x16, 0x5f, 0x69, 0x73,
	0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x64,
	0x64, 0x6f, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x72, 0x61, 0x77, 0x42, 0xe1, 0x01, 0x0a, 0x16,
	0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0d, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x2d, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x44, 0x58, 0x58, 0xaa, 0x02,
	0x12, 0x44, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0xca, 0x02, 0x12, 0x44, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x5c,
	0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x1e, 0x44, 0x61, 0x73, 0x68, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x44, 0x61, 0x73, 0x68,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dashboard_v1alpha1_template_proto_rawDescOnce sync.Once
	file_dashboard_v1alpha1_template_proto_rawDescData = file_dashboard_v1alpha1_template_proto_rawDesc
)

func file_dashboard_v1alpha1_template_proto_rawDescGZIP() []byte {
	file_dashboard_v1alpha1_template_proto_rawDescOnce.Do(func() {
		file_dashboard_v1alpha1_template_proto_rawDescData = protoimpl.X.CompressGZIP(file_dashboard_v1alpha1_template_proto_rawDescData)
	})
	return file_dashboard_v1alpha1_template_proto_rawDescData
}

var file_dashboard_v1alpha1_template_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_dashboard_v1alpha1_template_proto_goTypes = []interface{}{
	(*TemplateRequiredVars)(nil), // 0: dashboard.v1alpha1.TemplateRequiredVars
	(*Template)(nil),             // 1: dashboard.v1alpha1.Template
}
var file_dashboard_v1alpha1_template_proto_depIdxs = []int32{
	0, // 0: dashboard.v1alpha1.Template.required_vars:type_name -> dashboard.v1alpha1.TemplateRequiredVars
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_dashboard_v1alpha1_template_proto_init() }
func file_dashboard_v1alpha1_template_proto_init() {
	if File_dashboard_v1alpha1_template_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dashboard_v1alpha1_template_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateRequiredVars); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_dashboard_v1alpha1_template_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Template); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_dashboard_v1alpha1_template_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dashboard_v1alpha1_template_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dashboard_v1alpha1_template_proto_goTypes,
		DependencyIndexes: file_dashboard_v1alpha1_template_proto_depIdxs,
		MessageInfos:      file_dashboard_v1alpha1_template_proto_msgTypes,
	}.Build()
	File_dashboard_v1alpha1_template_proto = out.File
	file_dashboard_v1alpha1_template_proto_rawDesc = nil
	file_dashboard_v1alpha1_template_proto_goTypes = nil
	file_dashboard_v1alpha1_template_proto_depIdxs = nil
}
