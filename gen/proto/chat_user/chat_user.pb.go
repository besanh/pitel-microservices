// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: proto/chat_user/chat_user.proto

package pb

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/structpb"
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

type PostChatUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Email    string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Level    string `protobuf:"bytes,4,opt,name=level,proto3" json:"level,omitempty"`
	Status   bool   `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
	FullName string `protobuf:"bytes,6,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	RoleId   string `protobuf:"bytes,7,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *PostChatUserRequest) Reset() {
	*x = PostChatUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_user_chat_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostChatUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostChatUserRequest) ProtoMessage() {}

func (x *PostChatUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_user_chat_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostChatUserRequest.ProtoReflect.Descriptor instead.
func (*PostChatUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_user_chat_user_proto_rawDescGZIP(), []int{0}
}

func (x *PostChatUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *PostChatUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *PostChatUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *PostChatUserRequest) GetLevel() string {
	if x != nil {
		return x.Level
	}
	return ""
}

func (x *PostChatUserRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *PostChatUserRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *PostChatUserRequest) GetRoleId() string {
	if x != nil {
		return x.RoleId
	}
	return ""
}

type PostChatUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Id      string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PostChatUserResponse) Reset() {
	*x = PostChatUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_user_chat_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostChatUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostChatUserResponse) ProtoMessage() {}

func (x *PostChatUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_user_chat_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostChatUserResponse.ProtoReflect.Descriptor instead.
func (*PostChatUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_user_chat_user_proto_rawDescGZIP(), []int{1}
}

func (x *PostChatUserResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PostChatUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *PostChatUserResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PutChatUserStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PutChatUserStatusRequest) Reset() {
	*x = PutChatUserStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_user_chat_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutChatUserStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutChatUserStatusRequest) ProtoMessage() {}

func (x *PutChatUserStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_user_chat_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutChatUserStatusRequest.ProtoReflect.Descriptor instead.
func (*PutChatUserStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_user_chat_user_proto_rawDescGZIP(), []int{2}
}

func (x *PutChatUserStatusRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type PutChatUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PutChatUserResponse) Reset() {
	*x = PutChatUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_user_chat_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutChatUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutChatUserResponse) ProtoMessage() {}

func (x *PutChatUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_user_chat_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutChatUserResponse.ProtoReflect.Descriptor instead.
func (*PutChatUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_user_chat_user_proto_rawDescGZIP(), []int{3}
}

func (x *PutChatUserResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PutChatUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetChatUserStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetChatUserStatusRequest) Reset() {
	*x = GetChatUserStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_user_chat_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChatUserStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChatUserStatusRequest) ProtoMessage() {}

func (x *GetChatUserStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_user_chat_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChatUserStatusRequest.ProtoReflect.Descriptor instead.
func (*GetChatUserStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_user_chat_user_proto_rawDescGZIP(), []int{4}
}

type GetChatUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Status  string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *GetChatUserResponse) Reset() {
	*x = GetChatUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_user_chat_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChatUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChatUserResponse) ProtoMessage() {}

func (x *GetChatUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_user_chat_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChatUserResponse.ProtoReflect.Descriptor instead.
func (*GetChatUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_user_chat_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetChatUserResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *GetChatUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *GetChatUserResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_chat_user_chat_user_proto protoreflect.FileDescriptor

var file_proto_chat_user_chat_user_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc7, 0x01, 0x0a, 0x13,
	0x50, 0x6f, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x6f, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x54, 0x0a, 0x14, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x68, 0x61,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x32, 0x0a, 0x18, 0x50,
	0x75, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x43, 0x0a, 0x13, 0x50, 0x75, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x1a, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x5b, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0xb5, 0x03,
	0x0a, 0x0f, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x7c, 0x0a, 0x0c, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x21, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x1b, 0x3a, 0x01, 0x2a, 0x22, 0x16, 0x2f, 0x62, 0x73, 0x73, 0x2d, 0x63, 0x68,
	0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x12,
	0x93, 0x01, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x49, 0x64, 0x12, 0x28, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x50, 0x75,
	0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x22, 0x3a, 0x01, 0x2a, 0x1a, 0x1d, 0x2f, 0x62, 0x73, 0x73, 0x2d, 0x63, 0x68, 0x61,
	0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x8d, 0x01, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x49, 0x64, 0x12,
	0x28, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x25,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x12, 0x1d, 0x2f, 0x62, 0x73, 0x73, 0x2d, 0x63, 0x68, 0x61,
	0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6c, 0x34, 0x76, 0x6e, 0x2f, 0x70, 0x69, 0x74, 0x65, 0x6c,
	0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_chat_user_chat_user_proto_rawDescOnce sync.Once
	file_proto_chat_user_chat_user_proto_rawDescData = file_proto_chat_user_chat_user_proto_rawDesc
)

func file_proto_chat_user_chat_user_proto_rawDescGZIP() []byte {
	file_proto_chat_user_chat_user_proto_rawDescOnce.Do(func() {
		file_proto_chat_user_chat_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_chat_user_chat_user_proto_rawDescData)
	})
	return file_proto_chat_user_chat_user_proto_rawDescData
}

var file_proto_chat_user_chat_user_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_chat_user_chat_user_proto_goTypes = []interface{}{
	(*PostChatUserRequest)(nil),      // 0: proto.chatUser.PostChatUserRequest
	(*PostChatUserResponse)(nil),     // 1: proto.chatUser.PostChatUserResponse
	(*PutChatUserStatusRequest)(nil), // 2: proto.chatUser.PutChatUserStatusRequest
	(*PutChatUserResponse)(nil),      // 3: proto.chatUser.PutChatUserResponse
	(*GetChatUserStatusRequest)(nil), // 4: proto.chatUser.GetChatUserStatusRequest
	(*GetChatUserResponse)(nil),      // 5: proto.chatUser.GetChatUserResponse
}
var file_proto_chat_user_chat_user_proto_depIdxs = []int32{
	0, // 0: proto.chatUser.ChatUserService.PostChatUser:input_type -> proto.chatUser.PostChatUserRequest
	2, // 1: proto.chatUser.ChatUserService.UpdateChatUserStatusById:input_type -> proto.chatUser.PutChatUserStatusRequest
	4, // 2: proto.chatUser.ChatUserService.GetChatUserStatusById:input_type -> proto.chatUser.GetChatUserStatusRequest
	1, // 3: proto.chatUser.ChatUserService.PostChatUser:output_type -> proto.chatUser.PostChatUserResponse
	3, // 4: proto.chatUser.ChatUserService.UpdateChatUserStatusById:output_type -> proto.chatUser.PutChatUserResponse
	5, // 5: proto.chatUser.ChatUserService.GetChatUserStatusById:output_type -> proto.chatUser.GetChatUserResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_chat_user_chat_user_proto_init() }
func file_proto_chat_user_chat_user_proto_init() {
	if File_proto_chat_user_chat_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_chat_user_chat_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostChatUserRequest); i {
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
		file_proto_chat_user_chat_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostChatUserResponse); i {
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
		file_proto_chat_user_chat_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutChatUserStatusRequest); i {
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
		file_proto_chat_user_chat_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutChatUserResponse); i {
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
		file_proto_chat_user_chat_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChatUserStatusRequest); i {
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
		file_proto_chat_user_chat_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChatUserResponse); i {
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
			RawDescriptor: file_proto_chat_user_chat_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chat_user_chat_user_proto_goTypes,
		DependencyIndexes: file_proto_chat_user_chat_user_proto_depIdxs,
		MessageInfos:      file_proto_chat_user_chat_user_proto_msgTypes,
	}.Build()
	File_proto_chat_user_chat_user_proto = out.File
	file_proto_chat_user_chat_user_proto_rawDesc = nil
	file_proto_chat_user_chat_user_proto_goTypes = nil
	file_proto_chat_user_chat_user_proto_depIdxs = nil
}
