// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.9
// source: service_bitmoi.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Bitmoi_RequestCandles_FullMethodName  = "/pb.Bitmoi/RequestCandles"
	Bitmoi_PostScore_FullMethodName       = "/pb.Bitmoi/PostScore"
	Bitmoi_AnotherInterval_FullMethodName = "/pb.Bitmoi/AnotherInterval"
)

// BitmoiClient is the client API for Bitmoi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BitmoiClient interface {
	RequestCandles(ctx context.Context, in *GetCandlesRequest, opts ...grpc.CallOption) (*GetCandlesResponse, error)
	PostScore(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	AnotherInterval(ctx context.Context, in *AnotherIntervalRequest, opts ...grpc.CallOption) (*GetCandlesResponse, error)
}

type bitmoiClient struct {
	cc grpc.ClientConnInterface
}

func NewBitmoiClient(cc grpc.ClientConnInterface) BitmoiClient {
	return &bitmoiClient{cc}
}

func (c *bitmoiClient) RequestCandles(ctx context.Context, in *GetCandlesRequest, opts ...grpc.CallOption) (*GetCandlesResponse, error) {
	out := new(GetCandlesResponse)
	err := c.cc.Invoke(ctx, Bitmoi_RequestCandles_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bitmoiClient) PostScore(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, Bitmoi_PostScore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bitmoiClient) AnotherInterval(ctx context.Context, in *AnotherIntervalRequest, opts ...grpc.CallOption) (*GetCandlesResponse, error) {
	out := new(GetCandlesResponse)
	err := c.cc.Invoke(ctx, Bitmoi_AnotherInterval_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BitmoiServer is the server API for Bitmoi service.
// All implementations must embed UnimplementedBitmoiServer
// for forward compatibility
type BitmoiServer interface {
	RequestCandles(context.Context, *GetCandlesRequest) (*GetCandlesResponse, error)
	PostScore(context.Context, *OrderRequest) (*OrderResponse, error)
	AnotherInterval(context.Context, *AnotherIntervalRequest) (*GetCandlesResponse, error)
	mustEmbedUnimplementedBitmoiServer()
}

// UnimplementedBitmoiServer must be embedded to have forward compatible implementations.
type UnimplementedBitmoiServer struct {
}

func (UnimplementedBitmoiServer) RequestCandles(context.Context, *GetCandlesRequest) (*GetCandlesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestCandles not implemented")
}
func (UnimplementedBitmoiServer) PostScore(context.Context, *OrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostScore not implemented")
}
func (UnimplementedBitmoiServer) AnotherInterval(context.Context, *AnotherIntervalRequest) (*GetCandlesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnotherInterval not implemented")
}
func (UnimplementedBitmoiServer) mustEmbedUnimplementedBitmoiServer() {}

// UnsafeBitmoiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BitmoiServer will
// result in compilation errors.
type UnsafeBitmoiServer interface {
	mustEmbedUnimplementedBitmoiServer()
}

func RegisterBitmoiServer(s grpc.ServiceRegistrar, srv BitmoiServer) {
	s.RegisterService(&Bitmoi_ServiceDesc, srv)
}

func _Bitmoi_RequestCandles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCandlesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitmoiServer).RequestCandles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bitmoi_RequestCandles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitmoiServer).RequestCandles(ctx, req.(*GetCandlesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bitmoi_PostScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitmoiServer).PostScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bitmoi_PostScore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitmoiServer).PostScore(ctx, req.(*OrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bitmoi_AnotherInterval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnotherIntervalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitmoiServer).AnotherInterval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bitmoi_AnotherInterval_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitmoiServer).AnotherInterval(ctx, req.(*AnotherIntervalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bitmoi_ServiceDesc is the grpc.ServiceDesc for Bitmoi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bitmoi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Bitmoi",
	HandlerType: (*BitmoiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestCandles",
			Handler:    _Bitmoi_RequestCandles_Handler,
		},
		{
			MethodName: "PostScore",
			Handler:    _Bitmoi_PostScore_Handler,
		},
		{
			MethodName: "AnotherInterval",
			Handler:    _Bitmoi_AnotherInterval_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_bitmoi.proto",
}
