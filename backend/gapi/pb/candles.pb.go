// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.9
// source: candles.proto

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

type GetCandlesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
	Mode  string   `protobuf:"bytes,2,opt,name=mode,proto3" json:"mode,omitempty"`
}

func (x *GetCandlesRequest) Reset() {
	*x = GetCandlesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_candles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCandlesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCandlesRequest) ProtoMessage() {}

func (x *GetCandlesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_candles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCandlesRequest.ProtoReflect.Descriptor instead.
func (*GetCandlesRequest) Descriptor() ([]byte, []int) {
	return file_candles_proto_rawDescGZIP(), []int{0}
}

func (x *GetCandlesRequest) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

func (x *GetCandlesRequest) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

type GetCandlesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Onechart   *Candles `protobuf:"bytes,2,opt,name=onechart,proto3" json:"onechart,omitempty"`
	Btcratio   float64  `protobuf:"fixed64,3,opt,name=btcratio,proto3" json:"btcratio,omitempty"`
	Entrytime  string   `protobuf:"bytes,4,opt,name=entrytime,proto3" json:"entrytime,omitempty"`
	Identifier string   `protobuf:"bytes,5,opt,name=identifier,proto3" json:"identifier,omitempty"`
}

func (x *GetCandlesResponse) Reset() {
	*x = GetCandlesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_candles_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCandlesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCandlesResponse) ProtoMessage() {}

func (x *GetCandlesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_candles_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCandlesResponse.ProtoReflect.Descriptor instead.
func (*GetCandlesResponse) Descriptor() ([]byte, []int) {
	return file_candles_proto_rawDescGZIP(), []int{1}
}

func (x *GetCandlesResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetCandlesResponse) GetOnechart() *Candles {
	if x != nil {
		return x.Onechart
	}
	return nil
}

func (x *GetCandlesResponse) GetBtcratio() float64 {
	if x != nil {
		return x.Btcratio
	}
	return 0
}

func (x *GetCandlesResponse) GetEntrytime() string {
	if x != nil {
		return x.Entrytime
	}
	return ""
}

func (x *GetCandlesResponse) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

type Candles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PData []*PriceData  `protobuf:"bytes,1,rep,name=pData,proto3" json:"pData,omitempty"`
	VData []*VolumeData `protobuf:"bytes,2,rep,name=vData,proto3" json:"vData,omitempty"`
}

func (x *Candles) Reset() {
	*x = Candles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_candles_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Candles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Candles) ProtoMessage() {}

func (x *Candles) ProtoReflect() protoreflect.Message {
	mi := &file_candles_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Candles.ProtoReflect.Descriptor instead.
func (*Candles) Descriptor() ([]byte, []int) {
	return file_candles_proto_rawDescGZIP(), []int{2}
}

func (x *Candles) GetPData() []*PriceData {
	if x != nil {
		return x.PData
	}
	return nil
}

func (x *Candles) GetVData() []*VolumeData {
	if x != nil {
		return x.VData
	}
	return nil
}

type PriceData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string  `protobuf:"bytes,1,opt,name=name,json=-,proto3" json:"name,omitempty"`
	Open  float64 `protobuf:"fixed64,2,opt,name=open,proto3" json:"open,omitempty"`
	Close float64 `protobuf:"fixed64,3,opt,name=close,proto3" json:"close,omitempty"`
	High  float64 `protobuf:"fixed64,4,opt,name=high,proto3" json:"high,omitempty"`
	Low   float64 `protobuf:"fixed64,5,opt,name=low,proto3" json:"low,omitempty"`
	Time  int64   `protobuf:"varint,6,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *PriceData) Reset() {
	*x = PriceData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_candles_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PriceData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PriceData) ProtoMessage() {}

func (x *PriceData) ProtoReflect() protoreflect.Message {
	mi := &file_candles_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PriceData.ProtoReflect.Descriptor instead.
func (*PriceData) Descriptor() ([]byte, []int) {
	return file_candles_proto_rawDescGZIP(), []int{3}
}

func (x *PriceData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PriceData) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *PriceData) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *PriceData) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *PriceData) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *PriceData) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type VolumeData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
	Time  int64   `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Color string  `protobuf:"bytes,3,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *VolumeData) Reset() {
	*x = VolumeData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_candles_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VolumeData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VolumeData) ProtoMessage() {}

func (x *VolumeData) ProtoReflect() protoreflect.Message {
	mi := &file_candles_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VolumeData.ProtoReflect.Descriptor instead.
func (*VolumeData) Descriptor() ([]byte, []int) {
	return file_candles_proto_rawDescGZIP(), []int{4}
}

func (x *VolumeData) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *VolumeData) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *VolumeData) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

var File_candles_proto protoreflect.FileDescriptor

var file_candles_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x3d, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a,
	0x08, 0x6f, 0x6e, 0x65, 0x63, 0x68, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x52, 0x08, 0x6f, 0x6e,
	0x65, 0x63, 0x68, 0x61, 0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x74, 0x63, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x62, 0x74, 0x63, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x74, 0x69, 0x6d, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x22, 0x54, 0x0a, 0x07, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x05, 0x70,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x62, 0x2e,
	0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x70, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x24, 0x0a, 0x05, 0x76, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x05, 0x76, 0x44, 0x61, 0x74, 0x61, 0x22, 0x80, 0x01, 0x0a, 0x09, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x0f, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x01, 0x2d, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f,
	0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x68,
	0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x4c, 0x0a, 0x0a, 0x56, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x42, 0x18, 0x5a, 0x16, 0x62, 0x69, 0x74, 0x6d, 0x6f,
	0x69, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x67, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_candles_proto_rawDescOnce sync.Once
	file_candles_proto_rawDescData = file_candles_proto_rawDesc
)

func file_candles_proto_rawDescGZIP() []byte {
	file_candles_proto_rawDescOnce.Do(func() {
		file_candles_proto_rawDescData = protoimpl.X.CompressGZIP(file_candles_proto_rawDescData)
	})
	return file_candles_proto_rawDescData
}

var file_candles_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_candles_proto_goTypes = []interface{}{
	(*GetCandlesRequest)(nil),  // 0: pb.GetCandlesRequest
	(*GetCandlesResponse)(nil), // 1: pb.GetCandlesResponse
	(*Candles)(nil),            // 2: pb.Candles
	(*PriceData)(nil),          // 3: pb.PriceData
	(*VolumeData)(nil),         // 4: pb.VolumeData
}
var file_candles_proto_depIdxs = []int32{
	2, // 0: pb.GetCandlesResponse.onechart:type_name -> pb.Candles
	3, // 1: pb.Candles.pData:type_name -> pb.PriceData
	4, // 2: pb.Candles.vData:type_name -> pb.VolumeData
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_candles_proto_init() }
func file_candles_proto_init() {
	if File_candles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_candles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCandlesRequest); i {
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
		file_candles_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCandlesResponse); i {
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
		file_candles_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Candles); i {
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
		file_candles_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PriceData); i {
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
		file_candles_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VolumeData); i {
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
			RawDescriptor: file_candles_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_candles_proto_goTypes,
		DependencyIndexes: file_candles_proto_depIdxs,
		MessageInfos:      file_candles_proto_msgTypes,
	}.Build()
	File_candles_proto = out.File
	file_candles_proto_rawDesc = nil
	file_candles_proto_goTypes = nil
	file_candles_proto_depIdxs = nil
}
