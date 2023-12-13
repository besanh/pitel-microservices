// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/template_bss/template_bss.proto

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

type TemplateBssBodyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TemplateName string `protobuf:"bytes,1,opt,name=template_name,json=templateName,proto3" json:"template_name,omitempty"`
	TemplateCode string `protobuf:"bytes,2,opt,name=template_code,json=templateCode,proto3" json:"template_code,omitempty"`
	TemplateType string `protobuf:"bytes,4,opt,name=template_type,json=templateType,proto3" json:"template_type,omitempty"`
	Content      string `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	Status       bool   `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *TemplateBssBodyRequest) Reset() {
	*x = TemplateBssBodyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_template_bss_template_bss_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateBssBodyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateBssBodyRequest) ProtoMessage() {}

func (x *TemplateBssBodyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_template_bss_template_bss_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateBssBodyRequest.ProtoReflect.Descriptor instead.
func (*TemplateBssBodyRequest) Descriptor() ([]byte, []int) {
	return file_proto_template_bss_template_bss_proto_rawDescGZIP(), []int{0}
}

func (x *TemplateBssBodyRequest) GetTemplateName() string {
	if x != nil {
		return x.TemplateName
	}
	return ""
}

func (x *TemplateBssBodyRequest) GetTemplateCode() string {
	if x != nil {
		return x.TemplateCode
	}
	return ""
}

func (x *TemplateBssBodyRequest) GetTemplateType() string {
	if x != nil {
		return x.TemplateType
	}
	return ""
}

func (x *TemplateBssBodyRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *TemplateBssBodyRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type TemplateBssRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit        int32    `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset       int32    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	TemplateName string   `protobuf:"bytes,3,opt,name=template_name,json=templateName,proto3" json:"template_name,omitempty"`
	TemplateCode []string `protobuf:"bytes,4,rep,name=template_code,json=templateCode,proto3" json:"template_code,omitempty"`
	TemplateType []string `protobuf:"bytes,5,rep,name=template_type,json=templateType,proto3" json:"template_type,omitempty"`
	Content      string   `protobuf:"bytes,6,opt,name=content,proto3" json:"content,omitempty"`
	Status       string   `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *TemplateBssRequest) Reset() {
	*x = TemplateBssRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_template_bss_template_bss_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateBssRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateBssRequest) ProtoMessage() {}

func (x *TemplateBssRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_template_bss_template_bss_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateBssRequest.ProtoReflect.Descriptor instead.
func (*TemplateBssRequest) Descriptor() ([]byte, []int) {
	return file_proto_template_bss_template_bss_proto_rawDescGZIP(), []int{1}
}

func (x *TemplateBssRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TemplateBssRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *TemplateBssRequest) GetTemplateName() string {
	if x != nil {
		return x.TemplateName
	}
	return ""
}

func (x *TemplateBssRequest) GetTemplateCode() []string {
	if x != nil {
		return x.TemplateCode
	}
	return nil
}

func (x *TemplateBssRequest) GetTemplateType() []string {
	if x != nil {
		return x.TemplateType
	}
	return nil
}

func (x *TemplateBssRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *TemplateBssRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type TemplateBssByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TemplateBssByIdRequest) Reset() {
	*x = TemplateBssByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_template_bss_template_bss_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateBssByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateBssByIdRequest) ProtoMessage() {}

func (x *TemplateBssByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_template_bss_template_bss_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateBssByIdRequest.ProtoReflect.Descriptor instead.
func (*TemplateBssByIdRequest) Descriptor() ([]byte, []int) {
	return file_proto_template_bss_template_bss_proto_rawDescGZIP(), []int{2}
}

func (x *TemplateBssByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TemplateBssResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string             `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string             `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*structpb.Struct `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	Total   int32              `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *TemplateBssResponse) Reset() {
	*x = TemplateBssResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_template_bss_template_bss_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateBssResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateBssResponse) ProtoMessage() {}

func (x *TemplateBssResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_template_bss_template_bss_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateBssResponse.ProtoReflect.Descriptor instead.
func (*TemplateBssResponse) Descriptor() ([]byte, []int) {
	return file_proto_template_bss_template_bss_proto_rawDescGZIP(), []int{3}
}

func (x *TemplateBssResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *TemplateBssResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *TemplateBssResponse) GetData() []*structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *TemplateBssResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type TemplateBssPutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TemplateName string `protobuf:"bytes,2,opt,name=template_name,json=templateName,proto3" json:"template_name,omitempty"`
	TemplateCode string `protobuf:"bytes,3,opt,name=template_code,json=templateCode,proto3" json:"template_code,omitempty"`
	Content      string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Status       bool   `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *TemplateBssPutRequest) Reset() {
	*x = TemplateBssPutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_template_bss_template_bss_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateBssPutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateBssPutRequest) ProtoMessage() {}

func (x *TemplateBssPutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_template_bss_template_bss_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateBssPutRequest.ProtoReflect.Descriptor instead.
func (*TemplateBssPutRequest) Descriptor() ([]byte, []int) {
	return file_proto_template_bss_template_bss_proto_rawDescGZIP(), []int{4}
}

func (x *TemplateBssPutRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TemplateBssPutRequest) GetTemplateName() string {
	if x != nil {
		return x.TemplateName
	}
	return ""
}

func (x *TemplateBssPutRequest) GetTemplateCode() string {
	if x != nil {
		return x.TemplateCode
	}
	return ""
}

func (x *TemplateBssPutRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *TemplateBssPutRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type TemplateBssByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string           `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *structpb.Struct `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *TemplateBssByIdResponse) Reset() {
	*x = TemplateBssByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_template_bss_template_bss_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateBssByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateBssByIdResponse) ProtoMessage() {}

func (x *TemplateBssByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_template_bss_template_bss_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateBssByIdResponse.ProtoReflect.Descriptor instead.
func (*TemplateBssByIdResponse) Descriptor() ([]byte, []int) {
	return file_proto_template_bss_template_bss_proto_rawDescGZIP(), []int{5}
}

func (x *TemplateBssByIdResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *TemplateBssByIdResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *TemplateBssByIdResponse) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_template_bss_template_bss_proto protoreflect.FileDescriptor

var file_proto_template_bss_template_bss_proto_rawDesc = []byte{
	0x0a, 0x25, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x5f, 0x62, 0x73, 0x73, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb9, 0x01, 0x0a, 0x16, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0xe3, 0x01, 0x0a, 0x12, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x28, 0x0a, 0x16, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x86, 0x01, 0x0a, 0x13, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0xa3, 0x01, 0x0a, 0x15, 0x54,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x22, 0x74, 0x0a, 0x17, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42,
	0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xdd, 0x05, 0x0a, 0x0b, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x12, 0x8b, 0x01, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x54,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x12, 0x2a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x6f, 0x64, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e, 0x54, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a, 0x22, 0x14,
	0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x2d, 0x62, 0x73, 0x73, 0x12, 0x81, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x65, 0x73, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e, 0x54,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42,
	0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x16, 0x12, 0x14, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x2d, 0x62, 0x73, 0x73, 0x12, 0x90, 0x01, 0x0a, 0x12, 0x47, 0x65, 0x74,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64, 0x12,
	0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x5f, 0x62, 0x73, 0x73, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73,
	0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73,
	0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b,
	0x12, 0x19, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x2d, 0x62, 0x73, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x92, 0x01, 0x0a, 0x12,
	0x50, 0x75, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79,
	0x49, 0x64, 0x12, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x42, 0x73, 0x73, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62,
	0x73, 0x73, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79,
	0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1e, 0x3a, 0x01, 0x2a, 0x1a, 0x19, 0x2f, 0x62, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x62, 0x73, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d,
	0x12, 0x93, 0x01, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64, 0x12, 0x2a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x73, 0x73, 0x2e, 0x54, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x42, 0x73, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x2a, 0x19, 0x2f, 0x62, 0x73,
	0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x62, 0x73,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6c, 0x34, 0x76, 0x6e, 0x2f, 0x66, 0x69, 0x6e, 0x73,
	0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_template_bss_template_bss_proto_rawDescOnce sync.Once
	file_proto_template_bss_template_bss_proto_rawDescData = file_proto_template_bss_template_bss_proto_rawDesc
)

func file_proto_template_bss_template_bss_proto_rawDescGZIP() []byte {
	file_proto_template_bss_template_bss_proto_rawDescOnce.Do(func() {
		file_proto_template_bss_template_bss_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_template_bss_template_bss_proto_rawDescData)
	})
	return file_proto_template_bss_template_bss_proto_rawDescData
}

var file_proto_template_bss_template_bss_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_template_bss_template_bss_proto_goTypes = []interface{}{
	(*TemplateBssBodyRequest)(nil),  // 0: proto.template_bss.TemplateBssBodyRequest
	(*TemplateBssRequest)(nil),      // 1: proto.template_bss.TemplateBssRequest
	(*TemplateBssByIdRequest)(nil),  // 2: proto.template_bss.TemplateBssByIdRequest
	(*TemplateBssResponse)(nil),     // 3: proto.template_bss.TemplateBssResponse
	(*TemplateBssPutRequest)(nil),   // 4: proto.template_bss.TemplateBssPutRequest
	(*TemplateBssByIdResponse)(nil), // 5: proto.template_bss.TemplateBssByIdResponse
	(*structpb.Struct)(nil),         // 6: google.protobuf.Struct
}
var file_proto_template_bss_template_bss_proto_depIdxs = []int32{
	6, // 0: proto.template_bss.TemplateBssResponse.data:type_name -> google.protobuf.Struct
	6, // 1: proto.template_bss.TemplateBssByIdResponse.data:type_name -> google.protobuf.Struct
	0, // 2: proto.template_bss.TemplateBss.PostTemplateBss:input_type -> proto.template_bss.TemplateBssBodyRequest
	1, // 3: proto.template_bss.TemplateBss.GetTemplateBsses:input_type -> proto.template_bss.TemplateBssRequest
	2, // 4: proto.template_bss.TemplateBss.GetTemplateBssById:input_type -> proto.template_bss.TemplateBssByIdRequest
	4, // 5: proto.template_bss.TemplateBss.PutTemplateBssById:input_type -> proto.template_bss.TemplateBssPutRequest
	2, // 6: proto.template_bss.TemplateBss.DeleteTemplateBssById:input_type -> proto.template_bss.TemplateBssByIdRequest
	5, // 7: proto.template_bss.TemplateBss.PostTemplateBss:output_type -> proto.template_bss.TemplateBssByIdResponse
	3, // 8: proto.template_bss.TemplateBss.GetTemplateBsses:output_type -> proto.template_bss.TemplateBssResponse
	5, // 9: proto.template_bss.TemplateBss.GetTemplateBssById:output_type -> proto.template_bss.TemplateBssByIdResponse
	5, // 10: proto.template_bss.TemplateBss.PutTemplateBssById:output_type -> proto.template_bss.TemplateBssByIdResponse
	5, // 11: proto.template_bss.TemplateBss.DeleteTemplateBssById:output_type -> proto.template_bss.TemplateBssByIdResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_template_bss_template_bss_proto_init() }
func file_proto_template_bss_template_bss_proto_init() {
	if File_proto_template_bss_template_bss_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_template_bss_template_bss_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateBssBodyRequest); i {
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
		file_proto_template_bss_template_bss_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateBssRequest); i {
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
		file_proto_template_bss_template_bss_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateBssByIdRequest); i {
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
		file_proto_template_bss_template_bss_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateBssResponse); i {
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
		file_proto_template_bss_template_bss_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateBssPutRequest); i {
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
		file_proto_template_bss_template_bss_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemplateBssByIdResponse); i {
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
			RawDescriptor: file_proto_template_bss_template_bss_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_template_bss_template_bss_proto_goTypes,
		DependencyIndexes: file_proto_template_bss_template_bss_proto_depIdxs,
		MessageInfos:      file_proto_template_bss_template_bss_proto_msgTypes,
	}.Build()
	File_proto_template_bss_template_bss_proto = out.File
	file_proto_template_bss_template_bss_proto_rawDesc = nil
	file_proto_template_bss_template_bss_proto_goTypes = nil
	file_proto_template_bss_template_bss_proto_depIdxs = nil
}
