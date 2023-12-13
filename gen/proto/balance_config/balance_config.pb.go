// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/balance_config/balance_config.proto

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

type BalanceConfigBodyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Weight      string `protobuf:"bytes,1,opt,name=weight,proto3" json:"weight,omitempty"`
	BalanceType string `protobuf:"bytes,2,opt,name=balance_type,json=balanceType,proto3" json:"balance_type,omitempty"`
	Priority    string `protobuf:"bytes,3,opt,name=priority,proto3" json:"priority,omitempty"`
	Provider    string `protobuf:"bytes,4,opt,name=provider,proto3" json:"provider,omitempty"`
	Status      bool   `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *BalanceConfigBodyRequest) Reset() {
	*x = BalanceConfigBodyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_config_balance_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceConfigBodyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceConfigBodyRequest) ProtoMessage() {}

func (x *BalanceConfigBodyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_config_balance_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceConfigBodyRequest.ProtoReflect.Descriptor instead.
func (*BalanceConfigBodyRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_config_balance_config_proto_rawDescGZIP(), []int{0}
}

func (x *BalanceConfigBodyRequest) GetWeight() string {
	if x != nil {
		return x.Weight
	}
	return ""
}

func (x *BalanceConfigBodyRequest) GetBalanceType() string {
	if x != nil {
		return x.BalanceType
	}
	return ""
}

func (x *BalanceConfigBodyRequest) GetPriority() string {
	if x != nil {
		return x.Priority
	}
	return ""
}

func (x *BalanceConfigBodyRequest) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *BalanceConfigBodyRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type BalanceConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit    int32    `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset   int32    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Weight   []string `protobuf:"bytes,3,rep,name=weight,proto3" json:"weight,omitempty"`
	Priority []string `protobuf:"bytes,4,rep,name=priority,proto3" json:"priority,omitempty"`
	Provider []string `protobuf:"bytes,5,rep,name=provider,proto3" json:"provider,omitempty"`
	Status   string   `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *BalanceConfigRequest) Reset() {
	*x = BalanceConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_config_balance_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceConfigRequest) ProtoMessage() {}

func (x *BalanceConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_config_balance_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceConfigRequest.ProtoReflect.Descriptor instead.
func (*BalanceConfigRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_config_balance_config_proto_rawDescGZIP(), []int{1}
}

func (x *BalanceConfigRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *BalanceConfigRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *BalanceConfigRequest) GetWeight() []string {
	if x != nil {
		return x.Weight
	}
	return nil
}

func (x *BalanceConfigRequest) GetPriority() []string {
	if x != nil {
		return x.Priority
	}
	return nil
}

func (x *BalanceConfigRequest) GetProvider() []string {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *BalanceConfigRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type BalanceConfigByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *BalanceConfigByIdRequest) Reset() {
	*x = BalanceConfigByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_config_balance_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceConfigByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceConfigByIdRequest) ProtoMessage() {}

func (x *BalanceConfigByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_config_balance_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceConfigByIdRequest.ProtoReflect.Descriptor instead.
func (*BalanceConfigByIdRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_config_balance_config_proto_rawDescGZIP(), []int{2}
}

func (x *BalanceConfigByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type BalanceConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string             `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string             `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*structpb.Struct `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	Total   int32              `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *BalanceConfigResponse) Reset() {
	*x = BalanceConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_config_balance_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceConfigResponse) ProtoMessage() {}

func (x *BalanceConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_config_balance_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceConfigResponse.ProtoReflect.Descriptor instead.
func (*BalanceConfigResponse) Descriptor() ([]byte, []int) {
	return file_proto_balance_config_balance_config_proto_rawDescGZIP(), []int{3}
}

func (x *BalanceConfigResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *BalanceConfigResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *BalanceConfigResponse) GetData() []*structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *BalanceConfigResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type BalanceConfigPutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Balance     string `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	BalanceType string `protobuf:"bytes,3,opt,name=balance_type,json=balanceType,proto3" json:"balance_type,omitempty"`
	Provider    string `protobuf:"bytes,4,opt,name=provider,proto3" json:"provider,omitempty"`
	Priority    string `protobuf:"bytes,5,opt,name=priority,proto3" json:"priority,omitempty"`
	Status      bool   `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *BalanceConfigPutRequest) Reset() {
	*x = BalanceConfigPutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_config_balance_config_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceConfigPutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceConfigPutRequest) ProtoMessage() {}

func (x *BalanceConfigPutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_config_balance_config_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceConfigPutRequest.ProtoReflect.Descriptor instead.
func (*BalanceConfigPutRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_config_balance_config_proto_rawDescGZIP(), []int{4}
}

func (x *BalanceConfigPutRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BalanceConfigPutRequest) GetBalance() string {
	if x != nil {
		return x.Balance
	}
	return ""
}

func (x *BalanceConfigPutRequest) GetBalanceType() string {
	if x != nil {
		return x.BalanceType
	}
	return ""
}

func (x *BalanceConfigPutRequest) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *BalanceConfigPutRequest) GetPriority() string {
	if x != nil {
		return x.Priority
	}
	return ""
}

func (x *BalanceConfigPutRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type BalanceConfigByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string           `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *structpb.Struct `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *BalanceConfigByIdResponse) Reset() {
	*x = BalanceConfigByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_config_balance_config_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceConfigByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceConfigByIdResponse) ProtoMessage() {}

func (x *BalanceConfigByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_config_balance_config_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceConfigByIdResponse.ProtoReflect.Descriptor instead.
func (*BalanceConfigByIdResponse) Descriptor() ([]byte, []int) {
	return file_proto_balance_config_balance_config_proto_rawDescGZIP(), []int{5}
}

func (x *BalanceConfigByIdResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *BalanceConfigByIdResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *BalanceConfigByIdResponse) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_balance_config_balance_config_proto protoreflect.FileDescriptor

var file_proto_balance_config_balance_config_proto_rawDesc = []byte{
	0x0a, 0x29, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa5, 0x01, 0x0a, 0x18,
	0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x6f, 0x64,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0xac, 0x01, 0x0a, 0x14, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x2a, 0x0a, 0x18, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x88,
	0x01, 0x0a, 0x15, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0xb6, 0x01, 0x0a, 0x17, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x75, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x76, 0x0a, 0x19, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x9a, 0x06, 0x0a, 0x0d, 0x42,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x97, 0x01, 0x0a,
	0x11, 0x50, 0x6f, 0x73, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x3a, 0x01, 0x2a, 0x22, 0x16,
	0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2d,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x8c, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x2a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f,
	0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2d, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x9c, 0x01, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x12, 0x2e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31,
	0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x12, 0x9e, 0x01, 0x0a, 0x14, 0x50, 0x75, 0x74, 0x42, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x12, 0x2d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x26, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x20, 0x3a, 0x01, 0x2a, 0x1a, 0x1b, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76,
	0x31, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x9f, 0x01, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49,
	0x64, 0x12, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x2a, 0x1b, 0x2f, 0x62, 0x73, 0x73,
	0x2f, 0x76, 0x31, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2d, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6c, 0x34, 0x76, 0x6e, 0x2f, 0x66, 0x69, 0x6e,
	0x73, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_balance_config_balance_config_proto_rawDescOnce sync.Once
	file_proto_balance_config_balance_config_proto_rawDescData = file_proto_balance_config_balance_config_proto_rawDesc
)

func file_proto_balance_config_balance_config_proto_rawDescGZIP() []byte {
	file_proto_balance_config_balance_config_proto_rawDescOnce.Do(func() {
		file_proto_balance_config_balance_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_balance_config_balance_config_proto_rawDescData)
	})
	return file_proto_balance_config_balance_config_proto_rawDescData
}

var file_proto_balance_config_balance_config_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_balance_config_balance_config_proto_goTypes = []interface{}{
	(*BalanceConfigBodyRequest)(nil),  // 0: proto.balance_config.BalanceConfigBodyRequest
	(*BalanceConfigRequest)(nil),      // 1: proto.balance_config.BalanceConfigRequest
	(*BalanceConfigByIdRequest)(nil),  // 2: proto.balance_config.BalanceConfigByIdRequest
	(*BalanceConfigResponse)(nil),     // 3: proto.balance_config.BalanceConfigResponse
	(*BalanceConfigPutRequest)(nil),   // 4: proto.balance_config.BalanceConfigPutRequest
	(*BalanceConfigByIdResponse)(nil), // 5: proto.balance_config.BalanceConfigByIdResponse
	(*structpb.Struct)(nil),           // 6: google.protobuf.Struct
}
var file_proto_balance_config_balance_config_proto_depIdxs = []int32{
	6, // 0: proto.balance_config.BalanceConfigResponse.data:type_name -> google.protobuf.Struct
	6, // 1: proto.balance_config.BalanceConfigByIdResponse.data:type_name -> google.protobuf.Struct
	0, // 2: proto.balance_config.BalanceConfig.PostBalanceConfig:input_type -> proto.balance_config.BalanceConfigBodyRequest
	1, // 3: proto.balance_config.BalanceConfig.GetBalanceConfigs:input_type -> proto.balance_config.BalanceConfigRequest
	2, // 4: proto.balance_config.BalanceConfig.GetBalanceConfigById:input_type -> proto.balance_config.BalanceConfigByIdRequest
	4, // 5: proto.balance_config.BalanceConfig.PutBalanceConfigById:input_type -> proto.balance_config.BalanceConfigPutRequest
	2, // 6: proto.balance_config.BalanceConfig.DeleteBalanceConfigById:input_type -> proto.balance_config.BalanceConfigByIdRequest
	5, // 7: proto.balance_config.BalanceConfig.PostBalanceConfig:output_type -> proto.balance_config.BalanceConfigByIdResponse
	3, // 8: proto.balance_config.BalanceConfig.GetBalanceConfigs:output_type -> proto.balance_config.BalanceConfigResponse
	5, // 9: proto.balance_config.BalanceConfig.GetBalanceConfigById:output_type -> proto.balance_config.BalanceConfigByIdResponse
	5, // 10: proto.balance_config.BalanceConfig.PutBalanceConfigById:output_type -> proto.balance_config.BalanceConfigByIdResponse
	5, // 11: proto.balance_config.BalanceConfig.DeleteBalanceConfigById:output_type -> proto.balance_config.BalanceConfigByIdResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_balance_config_balance_config_proto_init() }
func file_proto_balance_config_balance_config_proto_init() {
	if File_proto_balance_config_balance_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_balance_config_balance_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceConfigBodyRequest); i {
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
		file_proto_balance_config_balance_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceConfigRequest); i {
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
		file_proto_balance_config_balance_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceConfigByIdRequest); i {
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
		file_proto_balance_config_balance_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceConfigResponse); i {
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
		file_proto_balance_config_balance_config_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceConfigPutRequest); i {
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
		file_proto_balance_config_balance_config_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceConfigByIdResponse); i {
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
			RawDescriptor: file_proto_balance_config_balance_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_balance_config_balance_config_proto_goTypes,
		DependencyIndexes: file_proto_balance_config_balance_config_proto_depIdxs,
		MessageInfos:      file_proto_balance_config_balance_config_proto_msgTypes,
	}.Build()
	File_proto_balance_config_balance_config_proto = out.File
	file_proto_balance_config_balance_config_proto_rawDesc = nil
	file_proto_balance_config_balance_config_proto_goTypes = nil
	file_proto_balance_config_balance_config_proto_depIdxs = nil
}
