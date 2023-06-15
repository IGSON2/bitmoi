// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.9
// source: order.proto

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

type OrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mode         string  `protobuf:"bytes,1,opt,name=mode,proto3" json:"mode,omitempty"`
	UserId       string  `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name         string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Stage        int32   `protobuf:"varint,4,opt,name=stage,proto3" json:"stage,omitempty"`
	IsLong       bool    `protobuf:"varint,5,opt,name=is_long,json=isLong,proto3" json:"is_long,omitempty"`
	EntryPrice   float64 `protobuf:"fixed64,6,opt,name=entry_price,json=entryPrice,proto3" json:"entry_price,omitempty"`
	Quantity     float64 `protobuf:"fixed64,7,opt,name=quantity,proto3" json:"quantity,omitempty"`
	QuantityRate float64 `protobuf:"fixed64,8,opt,name=quantity_rate,json=quantityRate,proto3" json:"quantity_rate,omitempty"`
	ProfitPrice  float64 `protobuf:"fixed64,9,opt,name=profit_price,json=profitPrice,proto3" json:"profit_price,omitempty"`
	LossPrice    float64 `protobuf:"fixed64,10,opt,name=loss_price,json=lossPrice,proto3" json:"loss_price,omitempty"`
	Leverage     int32   `protobuf:"varint,11,opt,name=leverage,proto3" json:"leverage,omitempty"`
	Balance      float64 `protobuf:"fixed64,12,opt,name=balance,proto3" json:"balance,omitempty"`
	Identifier   string  `protobuf:"bytes,13,opt,name=identifier,proto3" json:"identifier,omitempty"`
	ScoreId      string  `protobuf:"bytes,14,opt,name=score_id,json=scoreId,proto3" json:"score_id,omitempty"`
	WaitingTerm  int32   `protobuf:"varint,15,opt,name=waiting_term,json=waitingTerm,proto3" json:"waiting_term,omitempty"`
}

func (x *OrderRequest) Reset() {
	*x = OrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderRequest) ProtoMessage() {}

func (x *OrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderRequest.ProtoReflect.Descriptor instead.
func (*OrderRequest) Descriptor() ([]byte, []int) {
	return file_order_proto_rawDescGZIP(), []int{0}
}

func (x *OrderRequest) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

func (x *OrderRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *OrderRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OrderRequest) GetStage() int32 {
	if x != nil {
		return x.Stage
	}
	return 0
}

func (x *OrderRequest) GetIsLong() bool {
	if x != nil {
		return x.IsLong
	}
	return false
}

func (x *OrderRequest) GetEntryPrice() float64 {
	if x != nil {
		return x.EntryPrice
	}
	return 0
}

func (x *OrderRequest) GetQuantity() float64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *OrderRequest) GetQuantityRate() float64 {
	if x != nil {
		return x.QuantityRate
	}
	return 0
}

func (x *OrderRequest) GetProfitPrice() float64 {
	if x != nil {
		return x.ProfitPrice
	}
	return 0
}

func (x *OrderRequest) GetLossPrice() float64 {
	if x != nil {
		return x.LossPrice
	}
	return 0
}

func (x *OrderRequest) GetLeverage() int32 {
	if x != nil {
		return x.Leverage
	}
	return 0
}

func (x *OrderRequest) GetBalance() float64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

func (x *OrderRequest) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

func (x *OrderRequest) GetScoreId() string {
	if x != nil {
		return x.ScoreId
	}
	return ""
}

func (x *OrderRequest) GetWaitingTerm() int32 {
	if x != nil {
		return x.WaitingTerm
	}
	return 0
}

type OrderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginChart *CandleData `protobuf:"bytes,1,opt,name=origin_chart,json=originChart,proto3" json:"origin_chart,omitempty"`
	ResultChart *CandleData `protobuf:"bytes,2,opt,name=result_chart,json=resultChart,proto3" json:"result_chart,omitempty"`
	Score       *Score      `protobuf:"bytes,3,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *OrderResponse) Reset() {
	*x = OrderResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderResponse) ProtoMessage() {}

func (x *OrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_order_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderResponse.ProtoReflect.Descriptor instead.
func (*OrderResponse) Descriptor() ([]byte, []int) {
	return file_order_proto_rawDescGZIP(), []int{1}
}

func (x *OrderResponse) GetOriginChart() *CandleData {
	if x != nil {
		return x.OriginChart
	}
	return nil
}

func (x *OrderResponse) GetResultChart() *CandleData {
	if x != nil {
		return x.ResultChart
	}
	return nil
}

func (x *OrderResponse) GetScore() *Score {
	if x != nil {
		return x.Score
	}
	return nil
}

type Score struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stage        int32   `protobuf:"varint,1,opt,name=stage,proto3" json:"stage,omitempty"`
	Name         string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Entrytime    string  `protobuf:"bytes,3,opt,name=entrytime,proto3" json:"entrytime,omitempty"`
	Leverage     int32   `protobuf:"varint,4,opt,name=leverage,proto3" json:"leverage,omitempty"`
	EntryPrice   float64 `protobuf:"fixed64,5,opt,name=entry_price,json=entryPrice,proto3" json:"entry_price,omitempty"`
	EndPrice     float64 `protobuf:"fixed64,6,opt,name=end_price,json=endPrice,proto3" json:"end_price,omitempty"`
	OutTime      int32   `protobuf:"varint,7,opt,name=out_time,json=outTime,proto3" json:"out_time,omitempty"`
	Roe          float64 `protobuf:"fixed64,8,opt,name=roe,proto3" json:"roe,omitempty"`
	Pnl          float64 `protobuf:"fixed64,9,opt,name=pnl,proto3" json:"pnl,omitempty"`
	Commission   float64 `protobuf:"fixed64,10,opt,name=commission,proto3" json:"commission,omitempty"`
	IsLiquidated bool    `protobuf:"varint,11,opt,name=is_liquidated,json=isLiquidated,proto3" json:"is_liquidated,omitempty"`
}

func (x *Score) Reset() {
	*x = Score{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Score) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Score) ProtoMessage() {}

func (x *Score) ProtoReflect() protoreflect.Message {
	mi := &file_order_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Score.ProtoReflect.Descriptor instead.
func (*Score) Descriptor() ([]byte, []int) {
	return file_order_proto_rawDescGZIP(), []int{2}
}

func (x *Score) GetStage() int32 {
	if x != nil {
		return x.Stage
	}
	return 0
}

func (x *Score) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Score) GetEntrytime() string {
	if x != nil {
		return x.Entrytime
	}
	return ""
}

func (x *Score) GetLeverage() int32 {
	if x != nil {
		return x.Leverage
	}
	return 0
}

func (x *Score) GetEntryPrice() float64 {
	if x != nil {
		return x.EntryPrice
	}
	return 0
}

func (x *Score) GetEndPrice() float64 {
	if x != nil {
		return x.EndPrice
	}
	return 0
}

func (x *Score) GetOutTime() int32 {
	if x != nil {
		return x.OutTime
	}
	return 0
}

func (x *Score) GetRoe() float64 {
	if x != nil {
		return x.Roe
	}
	return 0
}

func (x *Score) GetPnl() float64 {
	if x != nil {
		return x.Pnl
	}
	return 0
}

func (x *Score) GetCommission() float64 {
	if x != nil {
		return x.Commission
	}
	return 0
}

func (x *Score) GetIsLiquidated() bool {
	if x != nil {
		return x.IsLiquidated
	}
	return false
}

var File_order_proto protoreflect.FileDescriptor

var file_order_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70,
	0x62, 0x1a, 0x0d, 0x63, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82, 0x05, 0x0a, 0x0c, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1c, 0xfa, 0x42, 0x19, 0x72, 0x17, 0x52,
	0x08, 0x70, 0x72, 0x61, 0x63, 0x74, 0x69, 0x63, 0x65, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x65,
	0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x1a,
	0x04, 0x18, 0x0a, 0x20, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x69, 0x73, 0x5f, 0x6c, 0x6f, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69,
	0x73, 0x4c, 0x6f, 0x6e, 0x67, 0x12, 0x2f, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x12,
	0x09, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x72,
	0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x12, 0x09, 0x21,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x3c, 0x0a, 0x0d, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x72,
	0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x42, 0x17, 0xfa, 0x42, 0x14, 0x12, 0x12,
	0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x40, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x52, 0x0c, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x31, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x12, 0x09, 0x29, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x2d, 0x0a, 0x0a, 0x6c, 0x6f, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x12, 0x09, 0x29, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x52, 0x09, 0x6c, 0x6f, 0x73, 0x73, 0x50, 0x72, 0x69,
	0x63, 0x65, 0x12, 0x25, 0x0a, 0x08, 0x6c, 0x65, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x1a, 0x04, 0x18, 0x64, 0x20, 0x00, 0x52,
	0x08, 0x6c, 0x65, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x62, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x12,
	0x09, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x12, 0x27, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01,
	0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x08,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52, 0x07, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x49, 0x64,
	0x12, 0x2c, 0x0a, 0x0c, 0x77, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x65, 0x72, 0x6d,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x1a, 0x04, 0x18, 0x1e, 0x20,
	0x00, 0x52, 0x0b, 0x77, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x96,
	0x01, 0x0a, 0x0d, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x31, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x43, 0x68,
	0x61, 0x72, 0x74, 0x12, 0x31, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x63, 0x68,
	0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x43,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x43, 0x68, 0x61, 0x72, 0x74, 0x12, 0x1f, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x63, 0x6f, 0x72, 0x65,
	0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0xad, 0x02, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x65,
	0x6e, 0x74, 0x72, 0x79, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x65, 0x6e, 0x74, 0x72, 0x79, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6c, 0x65, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x72,
	0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x64, 0x5f, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x75, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x72, 0x6f, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x72, 0x6f, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x70, 0x6e, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x70,
	0x6e, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x69, 0x73, 0x4c, 0x69, 0x71,
	0x75, 0x69, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x18, 0x5a, 0x16, 0x62, 0x69, 0x74, 0x6d, 0x6f,
	0x69, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x67, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_order_proto_rawDescOnce sync.Once
	file_order_proto_rawDescData = file_order_proto_rawDesc
)

func file_order_proto_rawDescGZIP() []byte {
	file_order_proto_rawDescOnce.Do(func() {
		file_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_order_proto_rawDescData)
	})
	return file_order_proto_rawDescData
}

var file_order_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_order_proto_goTypes = []interface{}{
	(*OrderRequest)(nil),  // 0: pb.OrderRequest
	(*OrderResponse)(nil), // 1: pb.OrderResponse
	(*Score)(nil),         // 2: pb.Score
	(*CandleData)(nil),    // 3: pb.CandleData
}
var file_order_proto_depIdxs = []int32{
	3, // 0: pb.OrderResponse.origin_chart:type_name -> pb.CandleData
	3, // 1: pb.OrderResponse.result_chart:type_name -> pb.CandleData
	2, // 2: pb.OrderResponse.score:type_name -> pb.Score
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_order_proto_init() }
func file_order_proto_init() {
	if File_order_proto != nil {
		return
	}
	file_candles_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderRequest); i {
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
		file_order_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderResponse); i {
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
		file_order_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Score); i {
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
			RawDescriptor: file_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_order_proto_goTypes,
		DependencyIndexes: file_order_proto_depIdxs,
		MessageInfos:      file_order_proto_msgTypes,
	}.Build()
	File_order_proto = out.File
	file_order_proto_rawDesc = nil
	file_order_proto_goTypes = nil
	file_order_proto_depIdxs = nil
}
