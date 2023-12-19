// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/plugin_config/plugin_config.proto

package pb

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PluginConfigBodyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PluginName string `protobuf:"bytes,1,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	PluginType string `protobuf:"bytes,2,opt,name=plugin_type,json=pluginType,proto3" json:"plugin_type,omitempty"`
	Status     bool   `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PluginConfigBodyRequest) Reset() {
	*x = PluginConfigBodyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_config_plugin_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginConfigBodyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginConfigBodyRequest) ProtoMessage() {}

func (x *PluginConfigBodyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_config_plugin_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginConfigBodyRequest.ProtoReflect.Descriptor instead.
func (*PluginConfigBodyRequest) Descriptor() ([]byte, []int) {
	return file_proto_plugin_config_plugin_config_proto_rawDescGZIP(), []int{0}
}

func (x *PluginConfigBodyRequest) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

func (x *PluginConfigBodyRequest) GetPluginType() string {
	if x != nil {
		return x.PluginType
	}
	return ""
}

func (x *PluginConfigBodyRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type PluginConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit      int32    `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset     int32    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	PluginName []string `protobuf:"bytes,3,rep,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	PluginType []string `protobuf:"bytes,4,rep,name=plugin_type,json=pluginType,proto3" json:"plugin_type,omitempty"`
	Status     string   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PluginConfigRequest) Reset() {
	*x = PluginConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_config_plugin_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginConfigRequest) ProtoMessage() {}

func (x *PluginConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_config_plugin_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginConfigRequest.ProtoReflect.Descriptor instead.
func (*PluginConfigRequest) Descriptor() ([]byte, []int) {
	return file_proto_plugin_config_plugin_config_proto_rawDescGZIP(), []int{1}
}

func (x *PluginConfigRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *PluginConfigRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *PluginConfigRequest) GetPluginName() []string {
	if x != nil {
		return x.PluginName
	}
	return nil
}

func (x *PluginConfigRequest) GetPluginType() []string {
	if x != nil {
		return x.PluginType
	}
	return nil
}

func (x *PluginConfigRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type PluginConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string             `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string             `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*structpb.Struct `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	Total   int32              `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *PluginConfigResponse) Reset() {
	*x = PluginConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_config_plugin_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginConfigResponse) ProtoMessage() {}

func (x *PluginConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_config_plugin_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginConfigResponse.ProtoReflect.Descriptor instead.
func (*PluginConfigResponse) Descriptor() ([]byte, []int) {
	return file_proto_plugin_config_plugin_config_proto_rawDescGZIP(), []int{2}
}

func (x *PluginConfigResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PluginConfigResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *PluginConfigResponse) GetData() []*structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PluginConfigResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_proto_plugin_config_plugin_config_proto protoreflect.FileDescriptor

var file_proto_plugin_config_plugin_config_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x17, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x9d, 0x01,
	0x0a, 0x13, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f,
	0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x87, 0x01,
	0x0a, 0x14, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x32, 0xa7, 0x02, 0x0a, 0x0c, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x8d, 0x01, 0x0a, 0x10, 0x50, 0x6f, 0x73,
	0x74, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x2c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01,
	0x2a, 0x22, 0x15, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x86, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x28, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x62, 0x73, 0x73,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x65, 0x6c, 0x34, 0x76, 0x6e, 0x2f, 0x66, 0x69, 0x6e, 0x73, 0x2d, 0x6d, 0x69, 0x63, 0x72,
	0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_plugin_config_plugin_config_proto_rawDescOnce sync.Once
	file_proto_plugin_config_plugin_config_proto_rawDescData = file_proto_plugin_config_plugin_config_proto_rawDesc
)

func file_proto_plugin_config_plugin_config_proto_rawDescGZIP() []byte {
	file_proto_plugin_config_plugin_config_proto_rawDescOnce.Do(func() {
		file_proto_plugin_config_plugin_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_plugin_config_plugin_config_proto_rawDescData)
	})
	return file_proto_plugin_config_plugin_config_proto_rawDescData
}

var file_proto_plugin_config_plugin_config_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_plugin_config_plugin_config_proto_goTypes = []interface{}{
	(*PluginConfigBodyRequest)(nil), // 0: proto.plugin_config.PluginConfigBodyRequest
	(*PluginConfigRequest)(nil),     // 1: proto.plugin_config.PluginConfigRequest
	(*PluginConfigResponse)(nil),    // 2: proto.plugin_config.PluginConfigResponse
	(*structpb.Struct)(nil),         // 3: google.protobuf.Struct
}
var file_proto_plugin_config_plugin_config_proto_depIdxs = []int32{
	3, // 0: proto.plugin_config.PluginConfigResponse.data:type_name -> google.protobuf.Struct
	0, // 1: proto.plugin_config.PluginConfig.PostPluginConfig:input_type -> proto.plugin_config.PluginConfigBodyRequest
	1, // 2: proto.plugin_config.PluginConfig.GetPluginConfigs:input_type -> proto.plugin_config.PluginConfigRequest
	2, // 3: proto.plugin_config.PluginConfig.PostPluginConfig:output_type -> proto.plugin_config.PluginConfigResponse
	2, // 4: proto.plugin_config.PluginConfig.GetPluginConfigs:output_type -> proto.plugin_config.PluginConfigResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_plugin_config_plugin_config_proto_init() }
func file_proto_plugin_config_plugin_config_proto_init() {
	if File_proto_plugin_config_plugin_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_plugin_config_plugin_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginConfigBodyRequest); i {
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
		file_proto_plugin_config_plugin_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginConfigRequest); i {
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
		file_proto_plugin_config_plugin_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginConfigResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_plugin_config_plugin_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_plugin_config_plugin_config_proto_goTypes,
		DependencyIndexes: file_proto_plugin_config_plugin_config_proto_depIdxs,
		MessageInfos:      file_proto_plugin_config_plugin_config_proto_msgTypes,
	}.Build()
	File_proto_plugin_config_plugin_config_proto = out.File
	file_proto_plugin_config_plugin_config_proto_rawDesc = nil
	file_proto_plugin_config_plugin_config_proto_goTypes = nil
	file_proto_plugin_config_plugin_config_proto_depIdxs = nil
}
