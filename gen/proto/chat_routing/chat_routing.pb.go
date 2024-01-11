// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/chat_routing/chat_routing.proto

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

type ChatRoutingBodyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoutingName string `protobuf:"bytes,1,opt,name=routing_name,json=routingName,proto3" json:"routing_name,omitempty"`
	Status      bool   `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ChatRoutingBodyRequest) Reset() {
	*x = ChatRoutingBodyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_routing_chat_routing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatRoutingBodyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatRoutingBodyRequest) ProtoMessage() {}

func (x *ChatRoutingBodyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_routing_chat_routing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatRoutingBodyRequest.ProtoReflect.Descriptor instead.
func (*ChatRoutingBodyRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_routing_chat_routing_proto_rawDescGZIP(), []int{0}
}

func (x *ChatRoutingBodyRequest) GetRoutingName() string {
	if x != nil {
		return x.RoutingName
	}
	return ""
}

func (x *ChatRoutingBodyRequest) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type ChatRoutingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ChatRoutingResponse) Reset() {
	*x = ChatRoutingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chat_routing_chat_routing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatRoutingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatRoutingResponse) ProtoMessage() {}

func (x *ChatRoutingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_routing_chat_routing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatRoutingResponse.ProtoReflect.Descriptor instead.
func (*ChatRoutingResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_routing_chat_routing_proto_rawDescGZIP(), []int{1}
}

func (x *ChatRoutingResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *ChatRoutingResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_chat_routing_chat_routing_proto protoreflect.FileDescriptor

var file_proto_chat_routing_chat_routing_proto_rawDesc = []byte{
	0x0a, 0x25, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x72, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x16, 0x43, 0x68, 0x61, 0x74, 0x52, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x21, 0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x43, 0x0a, 0x13, 0x43, 0x68,
	0x61, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32,
	0x9c, 0x01, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x8c, 0x01, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74, 0x52, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x67, 0x12, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x5f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x72, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e,
	0x3a, 0x01, 0x2a, 0x22, 0x19, 0x2f, 0x62, 0x73, 0x73, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x36,
	0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6c,
	0x34, 0x76, 0x6e, 0x2f, 0x66, 0x69, 0x6e, 0x73, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_chat_routing_chat_routing_proto_rawDescOnce sync.Once
	file_proto_chat_routing_chat_routing_proto_rawDescData = file_proto_chat_routing_chat_routing_proto_rawDesc
)

func file_proto_chat_routing_chat_routing_proto_rawDescGZIP() []byte {
	file_proto_chat_routing_chat_routing_proto_rawDescOnce.Do(func() {
		file_proto_chat_routing_chat_routing_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_chat_routing_chat_routing_proto_rawDescData)
	})
	return file_proto_chat_routing_chat_routing_proto_rawDescData
}

var file_proto_chat_routing_chat_routing_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_chat_routing_chat_routing_proto_goTypes = []interface{}{
	(*ChatRoutingBodyRequest)(nil), // 0: proto.chat_routing.ChatRoutingBodyRequest
	(*ChatRoutingResponse)(nil),    // 1: proto.chat_routing.ChatRoutingResponse
}
var file_proto_chat_routing_chat_routing_proto_depIdxs = []int32{
	0, // 0: proto.chat_routing.ChatRouting.PostChatRouting:input_type -> proto.chat_routing.ChatRoutingBodyRequest
	1, // 1: proto.chat_routing.ChatRouting.PostChatRouting:output_type -> proto.chat_routing.ChatRoutingResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_chat_routing_chat_routing_proto_init() }
func file_proto_chat_routing_chat_routing_proto_init() {
	if File_proto_chat_routing_chat_routing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_chat_routing_chat_routing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatRoutingBodyRequest); i {
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
		file_proto_chat_routing_chat_routing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatRoutingResponse); i {
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
			RawDescriptor: file_proto_chat_routing_chat_routing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chat_routing_chat_routing_proto_goTypes,
		DependencyIndexes: file_proto_chat_routing_chat_routing_proto_depIdxs,
		MessageInfos:      file_proto_chat_routing_chat_routing_proto_msgTypes,
	}.Build()
	File_proto_chat_routing_chat_routing_proto = out.File
	file_proto_chat_routing_chat_routing_proto_rawDesc = nil
	file_proto_chat_routing_chat_routing_proto_goTypes = nil
	file_proto_chat_routing_chat_routing_proto_depIdxs = nil
}
