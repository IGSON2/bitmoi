// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: service_bitmoi.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_bitmoi_proto protoreflect.FileDescriptor

var file_service_bitmoi_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x62, 0x69, 0x74, 0x6d, 0x6f, 0x69,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0d, 0x63, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x88, 0x02, 0x0a,
	0x06, 0x42, 0x69, 0x74, 0x6d, 0x6f, 0x69, 0x12, 0x51, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x12, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x43,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x3a, 0x01, 0x2a, 0x22, 0x0b, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x12, 0x46, 0x0a, 0x09, 0x50, 0x6f,
	0x73, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x63, 0x6f,
	0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0e, 0x3a, 0x01, 0x2a, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x63, 0x6f,
	0x72, 0x65, 0x12, 0x63, 0x0a, 0x0f, 0x41, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x1a, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x6e, 0x6f, 0x74, 0x68,
	0x65, 0x72, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01,
	0x2a, 0x22, 0x14, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x42, 0x18, 0x5a, 0x16, 0x62, 0x69, 0x74, 0x6d, 0x6f,
	0x69, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x67, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_bitmoi_proto_goTypes = []interface{}{
	(*CandlesRequest)(nil),         // 0: pb.CandlesRequest
	(*ScoreRequest)(nil),           // 1: pb.ScoreRequest
	(*AnotherIntervalRequest)(nil), // 2: pb.AnotherIntervalRequest
	(*CandlesResponse)(nil),        // 3: pb.CandlesResponse
	(*ScoreResponse)(nil),          // 4: pb.ScoreResponse
}
var file_service_bitmoi_proto_depIdxs = []int32{
	0, // 0: pb.Bitmoi.RequestCandles:input_type -> pb.CandlesRequest
	1, // 1: pb.Bitmoi.PostScore:input_type -> pb.ScoreRequest
	2, // 2: pb.Bitmoi.AnotherInterval:input_type -> pb.AnotherIntervalRequest
	3, // 3: pb.Bitmoi.RequestCandles:output_type -> pb.CandlesResponse
	4, // 4: pb.Bitmoi.PostScore:output_type -> pb.ScoreResponse
	3, // 5: pb.Bitmoi.AnotherInterval:output_type -> pb.CandlesResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_bitmoi_proto_init() }
func file_service_bitmoi_proto_init() {
	if File_service_bitmoi_proto != nil {
		return
	}
	file_candles_proto_init()
	file_score_proto_init()
	file_another_interval_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_bitmoi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_bitmoi_proto_goTypes,
		DependencyIndexes: file_service_bitmoi_proto_depIdxs,
	}.Build()
	File_service_bitmoi_proto = out.File
	file_service_bitmoi_proto_rawDesc = nil
	file_service_bitmoi_proto_goTypes = nil
	file_service_bitmoi_proto_depIdxs = nil
}
