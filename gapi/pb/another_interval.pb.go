// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.26.1
// source: another_interval.proto

package pb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type AnotherIntervalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReqInterval string `protobuf:"bytes,1,opt,name=req_interval,json=reqInterval,proto3" json:"req_interval,omitempty"`
	Identifier  string `protobuf:"bytes,2,opt,name=identifier,proto3" json:"identifier,omitempty"`
	Mode        string `protobuf:"bytes,3,opt,name=mode,proto3" json:"mode,omitempty"`
	Stage       int32  `protobuf:"varint,4,opt,name=stage,proto3" json:"stage,omitempty"`
	UserId      string `protobuf:"bytes,5,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *AnotherIntervalRequest) Reset() {
	*x = AnotherIntervalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_another_interval_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnotherIntervalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnotherIntervalRequest) ProtoMessage() {}

func (x *AnotherIntervalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_another_interval_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnotherIntervalRequest.ProtoReflect.Descriptor instead.
func (*AnotherIntervalRequest) Descriptor() ([]byte, []int) {
	return file_another_interval_proto_rawDescGZIP(), []int{0}
}

func (x *AnotherIntervalRequest) GetReqInterval() string {
	if x != nil {
		return x.ReqInterval
	}
	return ""
}

func (x *AnotherIntervalRequest) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

func (x *AnotherIntervalRequest) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

func (x *AnotherIntervalRequest) GetStage() int32 {
	if x != nil {
		return x.Stage
	}
	return 0
}

func (x *AnotherIntervalRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

var File_another_interval_proto protoreflect.FileDescriptor

var file_another_interval_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xec, 0x01, 0x0a, 0x16, 0x41, 0x6e, 0x6f, 0x74, 0x68, 0x65,
	0x72, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x3d, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1a, 0xfa, 0x42, 0x17, 0x72, 0x15, 0x52, 0x02, 0x35,
	0x6d, 0x52, 0x03, 0x31, 0x35, 0x6d, 0x52, 0x02, 0x31, 0x68, 0x52, 0x02, 0x34, 0x68, 0x52, 0x02,
	0x31, 0x64, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12,
	0x27, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x0a, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1c, 0xfa, 0x42, 0x19, 0x72, 0x17, 0x52, 0x08, 0x70,
	0x72, 0x61, 0x63, 0x74, 0x69, 0x63, 0x65, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x74, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x1a, 0x04,
	0x18, 0x0a, 0x20, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x42, 0x18, 0x5a, 0x16, 0x62, 0x69, 0x74, 0x6d, 0x6f, 0x69, 0x2f, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x67, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_another_interval_proto_rawDescOnce sync.Once
	file_another_interval_proto_rawDescData = file_another_interval_proto_rawDesc
)

func file_another_interval_proto_rawDescGZIP() []byte {
	file_another_interval_proto_rawDescOnce.Do(func() {
		file_another_interval_proto_rawDescData = protoimpl.X.CompressGZIP(file_another_interval_proto_rawDescData)
	})
	return file_another_interval_proto_rawDescData
}

var file_another_interval_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_another_interval_proto_goTypes = []interface{}{
	(*AnotherIntervalRequest)(nil), // 0: pb.AnotherIntervalRequest
}
var file_another_interval_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_another_interval_proto_init() }
func file_another_interval_proto_init() {
	if File_another_interval_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_another_interval_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnotherIntervalRequest); i {
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
			RawDescriptor: file_another_interval_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_another_interval_proto_goTypes,
		DependencyIndexes: file_another_interval_proto_depIdxs,
		MessageInfos:      file_another_interval_proto_msgTypes,
	}.Build()
	File_another_interval_proto = out.File
	file_another_interval_proto_rawDesc = nil
	file_another_interval_proto_goTypes = nil
	file_another_interval_proto_depIdxs = nil
}
