// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.11.2
// source: push_service.proto

package pb

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

// 客户端订阅请求
type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_push_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_push_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_push_service_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

// 服务端主动推送的消息
type PushMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PushMessage) Reset() {
	*x = PushMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_push_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMessage) ProtoMessage() {}

func (x *PushMessage) ProtoReflect() protoreflect.Message {
	mi := &file_push_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMessage.ProtoReflect.Descriptor instead.
func (*PushMessage) Descriptor() ([]byte, []int) {
	return file_push_service_proto_rawDescGZIP(), []int{1}
}

func (x *PushMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_push_service_proto protoreflect.FileDescriptor

var file_push_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x75, 0x73, 0x68, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x70, 0x75, 0x73, 0x68, 0x22, 0x2f, 0x0a, 0x10, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x27, 0x0a, 0x0b, 0x50,
	0x75, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0x47, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x12, 0x16, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x2e,
	0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x30, 0x01, 0x42, 0x07, 0x5a,
	0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_push_service_proto_rawDescOnce sync.Once
	file_push_service_proto_rawDescData = file_push_service_proto_rawDesc
)

func file_push_service_proto_rawDescGZIP() []byte {
	file_push_service_proto_rawDescOnce.Do(func() {
		file_push_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_push_service_proto_rawDescData)
	})
	return file_push_service_proto_rawDescData
}

var file_push_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_push_service_proto_goTypes = []interface{}{
	(*SubscribeRequest)(nil), // 0: push.SubscribeRequest
	(*PushMessage)(nil),      // 1: push.PushMessage
}
var file_push_service_proto_depIdxs = []int32{
	0, // 0: push.PushService.Subscribe:input_type -> push.SubscribeRequest
	1, // 1: push.PushService.Subscribe:output_type -> push.PushMessage
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_push_service_proto_init() }
func file_push_service_proto_init() {
	if File_push_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_push_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeRequest); i {
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
		file_push_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMessage); i {
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
			RawDescriptor: file_push_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_push_service_proto_goTypes,
		DependencyIndexes: file_push_service_proto_depIdxs,
		MessageInfos:      file_push_service_proto_msgTypes,
	}.Build()
	File_push_service_proto = out.File
	file_push_service_proto_rawDesc = nil
	file_push_service_proto_goTypes = nil
	file_push_service_proto_depIdxs = nil
}
