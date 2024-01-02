// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/external_plugin/external_plugin_connect.proto

package inbox_marketing

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

type ExternalPluginConnectBodyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PluginName string  `protobuf:"bytes,1,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	PluginType string  `protobuf:"bytes,2,opt,name=plugin_type,json=pluginType,proto3" json:"plugin_type,omitempty"`
	Config     *Config `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *ExternalPluginConnectBodyRequest) Reset() {
	*x = ExternalPluginConnectBodyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExternalPluginConnectBodyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExternalPluginConnectBodyRequest) ProtoMessage() {}

func (x *ExternalPluginConnectBodyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExternalPluginConnectBodyRequest.ProtoReflect.Descriptor instead.
func (*ExternalPluginConnectBodyRequest) Descriptor() ([]byte, []int) {
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP(), []int{0}
}

func (x *ExternalPluginConnectBodyRequest) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

func (x *ExternalPluginConnectBodyRequest) GetPluginType() string {
	if x != nil {
		return x.PluginType
	}
	return ""
}

func (x *ExternalPluginConnectBodyRequest) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Incom  *Incom  `protobuf:"bytes,1,opt,name=incom,proto3" json:"incom,omitempty"`
	Abenla *Abenla `protobuf:"bytes,2,opt,name=abenla,proto3" json:"abenla,omitempty"`
	Fpt    *Fpt    `protobuf:"bytes,3,opt,name=fpt,proto3" json:"fpt,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP(), []int{1}
}

func (x *Config) GetIncom() *Incom {
	if x != nil {
		return x.Incom
	}
	return nil
}

func (x *Config) GetAbenla() *Abenla {
	if x != nil {
		return x.Abenla
	}
	return nil
}

func (x *Config) GetFpt() *Fpt {
	if x != nil {
		return x.Fpt
	}
	return nil
}

type Incom struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Status   bool   `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Incom) Reset() {
	*x = Incom{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Incom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Incom) ProtoMessage() {}

func (x *Incom) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Incom.ProtoReflect.Descriptor instead.
func (*Incom) Descriptor() ([]byte, []int) {
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP(), []int{2}
}

func (x *Incom) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Incom) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Incom) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type Fpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GrantType    string `protobuf:"bytes,1,opt,name=grant_type,json=grantType,proto3" json:"grant_type,omitempty"`
	ClientId     string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,3,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	Scope        string `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	Status       bool   `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Fpt) Reset() {
	*x = Fpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Fpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Fpt) ProtoMessage() {}

func (x *Fpt) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Fpt.ProtoReflect.Descriptor instead.
func (*Fpt) Descriptor() ([]byte, []int) {
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP(), []int{3}
}

func (x *Fpt) GetGrantType() string {
	if x != nil {
		return x.GrantType
	}
	return ""
}

func (x *Fpt) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Fpt) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *Fpt) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *Fpt) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type Abenla struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Status   bool   `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Abenla) Reset() {
	*x = Abenla{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Abenla) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Abenla) ProtoMessage() {}

func (x *Abenla) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Abenla.ProtoReflect.Descriptor instead.
func (*Abenla) Descriptor() ([]byte, []int) {
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP(), []int{4}
}

func (x *Abenla) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Abenla) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Abenla) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type ExternalPluginConnectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string           `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *structpb.Struct `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Id      string           `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ExternalPluginConnectResponse) Reset() {
	*x = ExternalPluginConnectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExternalPluginConnectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExternalPluginConnectResponse) ProtoMessage() {}

func (x *ExternalPluginConnectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_plugin_external_plugin_connect_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExternalPluginConnectResponse.ProtoReflect.Descriptor instead.
func (*ExternalPluginConnectResponse) Descriptor() ([]byte, []int) {
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP(), []int{5}
}

func (x *ExternalPluginConnectResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *ExternalPluginConnectResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ExternalPluginConnectResponse) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ExternalPluginConnectResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_proto_external_plugin_external_plugin_connect_proto protoreflect.FileDescriptor

var file_proto_external_plugin_external_plugin_connect_proto_rawDesc = []byte{
	0x0a, 0x33, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9b, 0x01, 0x0a, 0x20, 0x45, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x35, 0x0a,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x22, 0xa1, 0x01, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x32, 0x0a, 0x05, 0x69, 0x6e, 0x63, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x52, 0x05, 0x69, 0x6e,
	0x63, 0x6f, 0x6d, 0x12, 0x35, 0x0a, 0x06, 0x61, 0x62, 0x65, 0x6e, 0x6c, 0x61, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x41, 0x62, 0x65, 0x6e,
	0x6c, 0x61, 0x52, 0x06, 0x61, 0x62, 0x65, 0x6e, 0x6c, 0x61, 0x12, 0x2c, 0x0a, 0x03, 0x66, 0x70,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e,
	0x46, 0x70, 0x74, 0x52, 0x03, 0x66, 0x70, 0x74, 0x22, 0x57, 0x0a, 0x05, 0x49, 0x6e, 0x63, 0x6f,
	0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x94, 0x01, 0x0a, 0x03, 0x46, 0x70, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x61,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67,
	0x72, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x58, 0x0a, 0x06, 0x41, 0x62, 0x65, 0x6e,
	0x6c, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x8a, 0x01, 0x0a, 0x1d, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x32,
	0xd7, 0x01, 0x0a, 0x1c, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0xb6, 0x01, 0x0a, 0x11, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x37, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x45,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x34, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c, 0x3a, 0x01, 0x2a,
	0x22, 0x27, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2d, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2f, 0x6d, 0x69,
	0x63, 0x72, 0x6f, 0x2f, 0x76, 0x33, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x62,
	0x6f, 0x78, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_external_plugin_external_plugin_connect_proto_rawDescOnce sync.Once
	file_proto_external_plugin_external_plugin_connect_proto_rawDescData = file_proto_external_plugin_external_plugin_connect_proto_rawDesc
)

func file_proto_external_plugin_external_plugin_connect_proto_rawDescGZIP() []byte {
	file_proto_external_plugin_external_plugin_connect_proto_rawDescOnce.Do(func() {
		file_proto_external_plugin_external_plugin_connect_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_external_plugin_external_plugin_connect_proto_rawDescData)
	})
	return file_proto_external_plugin_external_plugin_connect_proto_rawDescData
}

var file_proto_external_plugin_external_plugin_connect_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_external_plugin_external_plugin_connect_proto_goTypes = []interface{}{
	(*ExternalPluginConnectBodyRequest)(nil), // 0: proto.external_plugin.ExternalPluginConnectBodyRequest
	(*Config)(nil),                           // 1: proto.external_plugin.Config
	(*Incom)(nil),                            // 2: proto.external_plugin.Incom
	(*Fpt)(nil),                              // 3: proto.external_plugin.Fpt
	(*Abenla)(nil),                           // 4: proto.external_plugin.Abenla
	(*ExternalPluginConnectResponse)(nil),    // 5: proto.external_plugin.ExternalPluginConnectResponse
	(*structpb.Struct)(nil),                  // 6: google.protobuf.Struct
}
var file_proto_external_plugin_external_plugin_connect_proto_depIdxs = []int32{
	1, // 0: proto.external_plugin.ExternalPluginConnectBodyRequest.config:type_name -> proto.external_plugin.Config
	2, // 1: proto.external_plugin.Config.incom:type_name -> proto.external_plugin.Incom
	4, // 2: proto.external_plugin.Config.abenla:type_name -> proto.external_plugin.Abenla
	3, // 3: proto.external_plugin.Config.fpt:type_name -> proto.external_plugin.Fpt
	6, // 4: proto.external_plugin.ExternalPluginConnectResponse.data:type_name -> google.protobuf.Struct
	0, // 5: proto.external_plugin.ExternalPluginConnectService.PostCreateConnect:input_type -> proto.external_plugin.ExternalPluginConnectBodyRequest
	5, // 6: proto.external_plugin.ExternalPluginConnectService.PostCreateConnect:output_type -> proto.external_plugin.ExternalPluginConnectResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_external_plugin_external_plugin_connect_proto_init() }
func file_proto_external_plugin_external_plugin_connect_proto_init() {
	if File_proto_external_plugin_external_plugin_connect_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_external_plugin_external_plugin_connect_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExternalPluginConnectBodyRequest); i {
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
		file_proto_external_plugin_external_plugin_connect_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_proto_external_plugin_external_plugin_connect_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Incom); i {
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
		file_proto_external_plugin_external_plugin_connect_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Fpt); i {
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
		file_proto_external_plugin_external_plugin_connect_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Abenla); i {
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
		file_proto_external_plugin_external_plugin_connect_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExternalPluginConnectResponse); i {
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
			RawDescriptor: file_proto_external_plugin_external_plugin_connect_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_external_plugin_external_plugin_connect_proto_goTypes,
		DependencyIndexes: file_proto_external_plugin_external_plugin_connect_proto_depIdxs,
		MessageInfos:      file_proto_external_plugin_external_plugin_connect_proto_msgTypes,
	}.Build()
	File_proto_external_plugin_external_plugin_connect_proto = out.File
	file_proto_external_plugin_external_plugin_connect_proto_rawDesc = nil
	file_proto_external_plugin_external_plugin_connect_proto_goTypes = nil
	file_proto_external_plugin_external_plugin_connect_proto_depIdxs = nil
}
