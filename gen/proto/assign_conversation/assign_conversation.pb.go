// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: proto/assign_conversation/assign_conversation.proto

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

// Post to InsertUserInQueue
type InsertUserInQueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId string `protobuf:"bytes,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	Status         string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	UserId         string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	QueueId        string `protobuf:"bytes,4,opt,name=queue_id,json=queueId,proto3" json:"queue_id,omitempty"`
}

func (x *InsertUserInQueueRequest) Reset() {
	*x = InsertUserInQueueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InsertUserInQueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsertUserInQueueRequest) ProtoMessage() {}

func (x *InsertUserInQueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsertUserInQueueRequest.ProtoReflect.Descriptor instead.
func (*InsertUserInQueueRequest) Descriptor() ([]byte, []int) {
	return file_proto_assign_conversation_assign_conversation_proto_rawDescGZIP(), []int{0}
}

func (x *InsertUserInQueueRequest) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

func (x *InsertUserInQueueRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *InsertUserInQueueRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *InsertUserInQueueRequest) GetQueueId() string {
	if x != nil {
		return x.QueueId
	}
	return ""
}

type AnyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *structpb.Struct `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Code int32            `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *AnyResponse) Reset() {
	*x = AnyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnyResponse) ProtoMessage() {}

func (x *AnyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnyResponse.ProtoReflect.Descriptor instead.
func (*AnyResponse) Descriptor() ([]byte, []int) {
	return file_proto_assign_conversation_assign_conversation_proto_rawDescGZIP(), []int{1}
}

func (x *AnyResponse) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AnyResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

// Get user in queue
type GetUserInQueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId            string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
	OaId             string `protobuf:"bytes,2,opt,name=oa_id,json=oaId,proto3" json:"oa_id,omitempty"`
	ConversationId   string `protobuf:"bytes,3,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	ConversationType string `protobuf:"bytes,4,opt,name=conversation_type,json=conversationType,proto3" json:"conversation_type,omitempty"`
	Status           string `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *GetUserInQueueRequest) Reset() {
	*x = GetUserInQueueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserInQueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInQueueRequest) ProtoMessage() {}

func (x *GetUserInQueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInQueueRequest.ProtoReflect.Descriptor instead.
func (*GetUserInQueueRequest) Descriptor() ([]byte, []int) {
	return file_proto_assign_conversation_assign_conversation_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserInQueueRequest) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *GetUserInQueueRequest) GetOaId() string {
	if x != nil {
		return x.OaId
	}
	return ""
}

func (x *GetUserInQueueRequest) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

func (x *GetUserInQueueRequest) GetConversationType() string {
	if x != nil {
		return x.ConversationType
	}
	return ""
}

func (x *GetUserInQueueRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// Get user assigned
type GetUserAssignedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *GetUserAssignedRequest) Reset() {
	*x = GetUserAssignedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserAssignedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserAssignedRequest) ProtoMessage() {}

func (x *GetUserAssignedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_assign_conversation_assign_conversation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserAssignedRequest.ProtoReflect.Descriptor instead.
func (*GetUserAssignedRequest) Descriptor() ([]byte, []int) {
	return file_proto_assign_conversation_assign_conversation_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserAssignedRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetUserAssignedRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_assign_conversation_assign_conversation_proto protoreflect.FileDescriptor

var file_proto_assign_conversation_assign_conversation_proto_rawDesc = []byte{
	0x0a, 0x33, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x63,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x73, 0x73,
	0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75,
	0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f, 0x01, 0x0a, 0x18, 0x49, 0x6e,
	0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72,
	0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x19, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x64, 0x22, 0x4e, 0x0a, 0x0b, 0x41,
	0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xb1, 0x01, 0x0a, 0x15,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x13, 0x0a, 0x05,
	0x6f, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6f, 0x61, 0x49,
	0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x11, 0x63, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x40, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x32, 0x94, 0x04, 0x0a, 0x19, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0xa7, 0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x12, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3b, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x35, 0x12, 0x33, 0x2f, 0x62, 0x73, 0x73, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f,
	0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x2d, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72,
	0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2d, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0xa0, 0x01, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76, 0x65,
	0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x12, 0x2e, 0x2f, 0x62,
	0x73, 0x73, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x2d, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x2d, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0xa9, 0x01, 0x0a,
	0x11, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x12, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x49, 0x6e,
	0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x33, 0x3a, 0x01, 0x2a, 0x22, 0x2e, 0x2f, 0x62, 0x73, 0x73, 0x2d, 0x63,
	0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x2d, 0x63, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2d,
	0x69, 0x6e, 0x2d, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6c, 0x34, 0x76, 0x6e, 0x2f, 0x66, 0x69,
	0x6e, 0x73, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_assign_conversation_assign_conversation_proto_rawDescOnce sync.Once
	file_proto_assign_conversation_assign_conversation_proto_rawDescData = file_proto_assign_conversation_assign_conversation_proto_rawDesc
)

func file_proto_assign_conversation_assign_conversation_proto_rawDescGZIP() []byte {
	file_proto_assign_conversation_assign_conversation_proto_rawDescOnce.Do(func() {
		file_proto_assign_conversation_assign_conversation_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_assign_conversation_assign_conversation_proto_rawDescData)
	})
	return file_proto_assign_conversation_assign_conversation_proto_rawDescData
}

var file_proto_assign_conversation_assign_conversation_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_assign_conversation_assign_conversation_proto_goTypes = []any{
	(*InsertUserInQueueRequest)(nil), // 0: proto.assignConversation.InsertUserInQueueRequest
	(*AnyResponse)(nil),              // 1: proto.assignConversation.AnyResponse
	(*GetUserInQueueRequest)(nil),    // 2: proto.assignConversation.GetUserInQueueRequest
	(*GetUserAssignedRequest)(nil),   // 3: proto.assignConversation.GetUserAssignedRequest
	(*structpb.Struct)(nil),          // 4: google.protobuf.Struct
}
var file_proto_assign_conversation_assign_conversation_proto_depIdxs = []int32{
	4, // 0: proto.assignConversation.AnyResponse.data:type_name -> google.protobuf.Struct
	3, // 1: proto.assignConversation.AssignConversationService.GetUserAssigned:input_type -> proto.assignConversation.GetUserAssignedRequest
	2, // 2: proto.assignConversation.AssignConversationService.GetUserInQueue:input_type -> proto.assignConversation.GetUserInQueueRequest
	0, // 3: proto.assignConversation.AssignConversationService.InsertUserInQueue:input_type -> proto.assignConversation.InsertUserInQueueRequest
	1, // 4: proto.assignConversation.AssignConversationService.GetUserAssigned:output_type -> proto.assignConversation.AnyResponse
	1, // 5: proto.assignConversation.AssignConversationService.GetUserInQueue:output_type -> proto.assignConversation.AnyResponse
	1, // 6: proto.assignConversation.AssignConversationService.InsertUserInQueue:output_type -> proto.assignConversation.AnyResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_assign_conversation_assign_conversation_proto_init() }
func file_proto_assign_conversation_assign_conversation_proto_init() {
	if File_proto_assign_conversation_assign_conversation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_assign_conversation_assign_conversation_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*InsertUserInQueueRequest); i {
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
		file_proto_assign_conversation_assign_conversation_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AnyResponse); i {
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
		file_proto_assign_conversation_assign_conversation_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GetUserInQueueRequest); i {
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
		file_proto_assign_conversation_assign_conversation_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetUserAssignedRequest); i {
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
			RawDescriptor: file_proto_assign_conversation_assign_conversation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_assign_conversation_assign_conversation_proto_goTypes,
		DependencyIndexes: file_proto_assign_conversation_assign_conversation_proto_depIdxs,
		MessageInfos:      file_proto_assign_conversation_assign_conversation_proto_msgTypes,
	}.Build()
	File_proto_assign_conversation_assign_conversation_proto = out.File
	file_proto_assign_conversation_assign_conversation_proto_rawDesc = nil
	file_proto_assign_conversation_assign_conversation_proto_goTypes = nil
	file_proto_assign_conversation_assign_conversation_proto_depIdxs = nil
}
