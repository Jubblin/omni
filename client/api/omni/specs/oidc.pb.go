// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.24.4
// source: omni/specs/oidc.proto

package specs

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// JWTPublicKeySpec keeps the active set of JWT signing keys.
type JWTPublicKeySpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PKCS1 encoded RSA public key.
	PublicKey []byte `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// Expiration time (when it's ready to be cleaned up).
	Expiration *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiration,proto3" json:"expiration,omitempty"`
}

func (x *JWTPublicKeySpec) Reset() {
	*x = JWTPublicKeySpec{}
	mi := &file_omni_specs_oidc_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JWTPublicKeySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JWTPublicKeySpec) ProtoMessage() {}

func (x *JWTPublicKeySpec) ProtoReflect() protoreflect.Message {
	mi := &file_omni_specs_oidc_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JWTPublicKeySpec.ProtoReflect.Descriptor instead.
func (*JWTPublicKeySpec) Descriptor() ([]byte, []int) {
	return file_omni_specs_oidc_proto_rawDescGZIP(), []int{0}
}

func (x *JWTPublicKeySpec) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *JWTPublicKeySpec) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

var File_omni_specs_oidc_proto protoreflect.FileDescriptor

var file_omni_specs_oidc_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x6f, 0x69, 0x64,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x70, 0x65, 0x63, 0x73, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x6d, 0x0a, 0x10, 0x4a, 0x57, 0x54, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b,
	0x65, 0x79, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x64,
	0x65, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x73, 0x70, 0x65,
	0x63, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_omni_specs_oidc_proto_rawDescOnce sync.Once
	file_omni_specs_oidc_proto_rawDescData = file_omni_specs_oidc_proto_rawDesc
)

func file_omni_specs_oidc_proto_rawDescGZIP() []byte {
	file_omni_specs_oidc_proto_rawDescOnce.Do(func() {
		file_omni_specs_oidc_proto_rawDescData = protoimpl.X.CompressGZIP(file_omni_specs_oidc_proto_rawDescData)
	})
	return file_omni_specs_oidc_proto_rawDescData
}

var file_omni_specs_oidc_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_omni_specs_oidc_proto_goTypes = []any{
	(*JWTPublicKeySpec)(nil),      // 0: specs.JWTPublicKeySpec
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_omni_specs_oidc_proto_depIdxs = []int32{
	1, // 0: specs.JWTPublicKeySpec.expiration:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_omni_specs_oidc_proto_init() }
func file_omni_specs_oidc_proto_init() {
	if File_omni_specs_oidc_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_omni_specs_oidc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_omni_specs_oidc_proto_goTypes,
		DependencyIndexes: file_omni_specs_oidc_proto_depIdxs,
		MessageInfos:      file_omni_specs_oidc_proto_msgTypes,
	}.Build()
	File_omni_specs_oidc_proto = out.File
	file_omni_specs_oidc_proto_rawDesc = nil
	file_omni_specs_oidc_proto_goTypes = nil
	file_omni_specs_oidc_proto_depIdxs = nil
}
