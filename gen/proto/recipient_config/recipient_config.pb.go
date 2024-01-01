// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/recipient_config/recipient_config.proto

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

type RecipientConfigBodyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Recipient     string `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	RecipientType string `protobuf:"bytes,2,opt,name=recipient_type,json=recipientType,proto3" json:"recipient_type,omitempty"`
	Provider      string `protobuf:"bytes,3,opt,name=provider,proto3" json:"provider,omitempty"`
	Priority      string `protobuf:"bytes,4,opt,name=priority,proto3" json:"priority,omitempty"`
	Status        bool   `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *RecipientConfigBodyRequest) Reset() {
	*x = RecipientConfigBodyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecipientConfigBodyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipientConfigBodyRequest) ProtoMessage() {}

func (x *RecipientConfigBodyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipientConfigBodyRequest.ProtoReflect.Descriptor instead.
func (*RecipientConfigBodyRequest) Descriptor() ([]byte, []int) {
	return file_proto_recipient_config_recipient_config_proto_rawDescGZIP(), []int{0}
}

func (x *RecipientConfigBodyRequest) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

func (x *RecipientConfigBodyRequest) GetRecipientType() string {
	if x != nil {
		return x.RecipientType
	}
	return ""
}

func (x *RecipientConfigBodyRequest) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *RecipientConfigBodyRequest) GetPriority() string {
	if x != nil {
		return x.Priority
	}
	return ""
}

func (x *RecipientConfigBodyRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type RecipientConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit     int32    `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset    int32    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Recipient []string `protobuf:"bytes,3,rep,name=recipient,proto3" json:"recipient,omitempty"`
	Priority  []string `protobuf:"bytes,4,rep,name=priority,proto3" json:"priority,omitempty"`
	Provider  []string `protobuf:"bytes,5,rep,name=provider,proto3" json:"provider,omitempty"`
	Status    string   `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *RecipientConfigRequest) Reset() {
	*x = RecipientConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecipientConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipientConfigRequest) ProtoMessage() {}

func (x *RecipientConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipientConfigRequest.ProtoReflect.Descriptor instead.
func (*RecipientConfigRequest) Descriptor() ([]byte, []int) {
	return file_proto_recipient_config_recipient_config_proto_rawDescGZIP(), []int{1}
}

func (x *RecipientConfigRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *RecipientConfigRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *RecipientConfigRequest) GetRecipient() []string {
	if x != nil {
		return x.Recipient
	}
	return nil
}

func (x *RecipientConfigRequest) GetPriority() []string {
	if x != nil {
		return x.Priority
	}
	return nil
}

func (x *RecipientConfigRequest) GetProvider() []string {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *RecipientConfigRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type RecipientConfigByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RecipientConfigByIdRequest) Reset() {
	*x = RecipientConfigByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecipientConfigByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipientConfigByIdRequest) ProtoMessage() {}

func (x *RecipientConfigByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipientConfigByIdRequest.ProtoReflect.Descriptor instead.
func (*RecipientConfigByIdRequest) Descriptor() ([]byte, []int) {
	return file_proto_recipient_config_recipient_config_proto_rawDescGZIP(), []int{2}
}

func (x *RecipientConfigByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RecipientConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string             `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string             `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*structpb.Struct `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	Total   int32              `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *RecipientConfigResponse) Reset() {
	*x = RecipientConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecipientConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipientConfigResponse) ProtoMessage() {}

func (x *RecipientConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipientConfigResponse.ProtoReflect.Descriptor instead.
func (*RecipientConfigResponse) Descriptor() ([]byte, []int) {
	return file_proto_recipient_config_recipient_config_proto_rawDescGZIP(), []int{3}
}

func (x *RecipientConfigResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *RecipientConfigResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RecipientConfigResponse) GetData() []*structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RecipientConfigResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type RecipientConfigPutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Recipient     string `protobuf:"bytes,2,opt,name=recipient,proto3" json:"recipient,omitempty"`
	RecipientType string `protobuf:"bytes,3,opt,name=recipient_type,json=recipientType,proto3" json:"recipient_type,omitempty"`
	Provider      string `protobuf:"bytes,4,opt,name=provider,proto3" json:"provider,omitempty"`
	Priority      string `protobuf:"bytes,5,opt,name=priority,proto3" json:"priority,omitempty"`
	Status        bool   `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *RecipientConfigPutRequest) Reset() {
	*x = RecipientConfigPutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecipientConfigPutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipientConfigPutRequest) ProtoMessage() {}

func (x *RecipientConfigPutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipientConfigPutRequest.ProtoReflect.Descriptor instead.
func (*RecipientConfigPutRequest) Descriptor() ([]byte, []int) {
	return file_proto_recipient_config_recipient_config_proto_rawDescGZIP(), []int{4}
}

func (x *RecipientConfigPutRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RecipientConfigPutRequest) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

func (x *RecipientConfigPutRequest) GetRecipientType() string {
	if x != nil {
		return x.RecipientType
	}
	return ""
}

func (x *RecipientConfigPutRequest) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *RecipientConfigPutRequest) GetPriority() string {
	if x != nil {
		return x.Priority
	}
	return ""
}

func (x *RecipientConfigPutRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type RecipientConfigByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string           `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *structpb.Struct `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Id      string           `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RecipientConfigByIdResponse) Reset() {
	*x = RecipientConfigByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecipientConfigByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipientConfigByIdResponse) ProtoMessage() {}

func (x *RecipientConfigByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_recipient_config_recipient_config_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipientConfigByIdResponse.ProtoReflect.Descriptor instead.
func (*RecipientConfigByIdResponse) Descriptor() ([]byte, []int) {
	return file_proto_recipient_config_recipient_config_proto_rawDescGZIP(), []int{5}
}

func (x *RecipientConfigByIdResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *RecipientConfigByIdResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RecipientConfigByIdResponse) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RecipientConfigByIdResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_proto_recipient_config_recipient_config_proto protoreflect.FileDescriptor

var file_proto_recipient_config_recipient_config_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xb1, 0x01, 0x0a, 0x1a, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x25,
	0x0a, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xb4, 0x01, 0x0a, 0x16, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2c, 0x0a, 0x1a,
	0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42,
	0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x8a, 0x01, 0x0a, 0x17, 0x52,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0xc0, 0x01, 0x0a, 0x19, 0x52, 0x65, 0x63, 0x69,
	0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x75, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x88, 0x01, 0x0a, 0x1b, 0x52,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79,
	0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x32, 0xd8, 0x06, 0x0a, 0x0f, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0xa3, 0x01, 0x0a, 0x13, 0x50, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79,
	0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1d, 0x3a, 0x01, 0x2a, 0x22, 0x18, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x72,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x98, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a,
	0x12, 0x18, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0xa8, 0x01, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x42, 0x79, 0x49, 0x64, 0x12, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79,
	0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x25,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x12, 0x1d, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f,
	0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0xaa, 0x01, 0x0a, 0x16, 0x50, 0x75, 0x74, 0x52, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64,
	0x12, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x63, 0x69,
	0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22,
	0x3a, 0x01, 0x2a, 0x1a, 0x1d, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x12, 0xab, 0x01, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64,
	0x12, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1f, 0x2a, 0x1d, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x7b, 0x69, 0x64, 0x7d,
	0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x65, 0x6c, 0x34, 0x76, 0x6e, 0x2f, 0x66, 0x69, 0x6e, 0x73, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_recipient_config_recipient_config_proto_rawDescOnce sync.Once
	file_proto_recipient_config_recipient_config_proto_rawDescData = file_proto_recipient_config_recipient_config_proto_rawDesc
)

func file_proto_recipient_config_recipient_config_proto_rawDescGZIP() []byte {
	file_proto_recipient_config_recipient_config_proto_rawDescOnce.Do(func() {
		file_proto_recipient_config_recipient_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_recipient_config_recipient_config_proto_rawDescData)
	})
	return file_proto_recipient_config_recipient_config_proto_rawDescData
}

var file_proto_recipient_config_recipient_config_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_recipient_config_recipient_config_proto_goTypes = []interface{}{
	(*RecipientConfigBodyRequest)(nil),  // 0: proto.recipient_config.RecipientConfigBodyRequest
	(*RecipientConfigRequest)(nil),      // 1: proto.recipient_config.RecipientConfigRequest
	(*RecipientConfigByIdRequest)(nil),  // 2: proto.recipient_config.RecipientConfigByIdRequest
	(*RecipientConfigResponse)(nil),     // 3: proto.recipient_config.RecipientConfigResponse
	(*RecipientConfigPutRequest)(nil),   // 4: proto.recipient_config.RecipientConfigPutRequest
	(*RecipientConfigByIdResponse)(nil), // 5: proto.recipient_config.RecipientConfigByIdResponse
	(*structpb.Struct)(nil),             // 6: google.protobuf.Struct
}
var file_proto_recipient_config_recipient_config_proto_depIdxs = []int32{
	6, // 0: proto.recipient_config.RecipientConfigResponse.data:type_name -> google.protobuf.Struct
	6, // 1: proto.recipient_config.RecipientConfigByIdResponse.data:type_name -> google.protobuf.Struct
	0, // 2: proto.recipient_config.RecipientConfig.PostRecipientConfig:input_type -> proto.recipient_config.RecipientConfigBodyRequest
	1, // 3: proto.recipient_config.RecipientConfig.GetRecipientConfigs:input_type -> proto.recipient_config.RecipientConfigRequest
	2, // 4: proto.recipient_config.RecipientConfig.GetRecipientConfigById:input_type -> proto.recipient_config.RecipientConfigByIdRequest
	4, // 5: proto.recipient_config.RecipientConfig.PutRecipientConfigById:input_type -> proto.recipient_config.RecipientConfigPutRequest
	2, // 6: proto.recipient_config.RecipientConfig.DeleteRecipientConfigById:input_type -> proto.recipient_config.RecipientConfigByIdRequest
	5, // 7: proto.recipient_config.RecipientConfig.PostRecipientConfig:output_type -> proto.recipient_config.RecipientConfigByIdResponse
	3, // 8: proto.recipient_config.RecipientConfig.GetRecipientConfigs:output_type -> proto.recipient_config.RecipientConfigResponse
	5, // 9: proto.recipient_config.RecipientConfig.GetRecipientConfigById:output_type -> proto.recipient_config.RecipientConfigByIdResponse
	5, // 10: proto.recipient_config.RecipientConfig.PutRecipientConfigById:output_type -> proto.recipient_config.RecipientConfigByIdResponse
	5, // 11: proto.recipient_config.RecipientConfig.DeleteRecipientConfigById:output_type -> proto.recipient_config.RecipientConfigByIdResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_recipient_config_recipient_config_proto_init() }
func file_proto_recipient_config_recipient_config_proto_init() {
	if File_proto_recipient_config_recipient_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_recipient_config_recipient_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecipientConfigBodyRequest); i {
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
		file_proto_recipient_config_recipient_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecipientConfigRequest); i {
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
		file_proto_recipient_config_recipient_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecipientConfigByIdRequest); i {
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
		file_proto_recipient_config_recipient_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecipientConfigResponse); i {
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
		file_proto_recipient_config_recipient_config_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecipientConfigPutRequest); i {
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
		file_proto_recipient_config_recipient_config_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecipientConfigByIdResponse); i {
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
			RawDescriptor: file_proto_recipient_config_recipient_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_recipient_config_recipient_config_proto_goTypes,
		DependencyIndexes: file_proto_recipient_config_recipient_config_proto_depIdxs,
		MessageInfos:      file_proto_recipient_config_recipient_config_proto_msgTypes,
	}.Build()
	File_proto_recipient_config_recipient_config_proto = out.File
	file_proto_recipient_config_recipient_config_proto_rawDesc = nil
	file_proto_recipient_config_recipient_config_proto_goTypes = nil
	file_proto_recipient_config_recipient_config_proto_depIdxs = nil
}
