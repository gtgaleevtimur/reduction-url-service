// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/proto.proto

package proto

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

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	AddByText(ctx context.Context, in *StringForm, opts ...grpc.CallOption) (*CommonResponse, error)
	GetByHashURL(ctx context.Context, in *StringForm, opts ...grpc.CallOption) (*CommonResponse, error)
	Ping(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*IntForm, error)
	Stats(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*StatsResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*IntForm, error)
	GetUserURLs(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*GetUserURLsResponse, error)
	PostJSON(ctx context.Context, in *PostJSONRespReq, opts ...grpc.CallOption) (*PostJSONRespReq, error)
	PostBatch(ctx context.Context, in *PostBatchRequest, opts ...grpc.CallOption) (*PostBatchResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) AddByText(ctx context.Context, in *StringForm, opts ...grpc.CallOption) (*CommonResponse, error) {
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/AddByText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetByHashURL(ctx context.Context, in *StringForm, opts ...grpc.CallOption) (*CommonResponse, error) {
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetByHashURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Ping(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*IntForm, error) {
	out := new(IntForm)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Stats(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/Stats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*IntForm, error) {
	out := new(IntForm)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetUserURLs(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*GetUserURLsResponse, error) {
	out := new(GetUserURLsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) PostJSON(ctx context.Context, in *PostJSONRespReq, opts ...grpc.CallOption) (*PostJSONRespReq, error) {
	out := new(PostJSONRespReq)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/PostJSON", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) PostBatch(ctx context.Context, in *PostBatchRequest, opts ...grpc.CallOption) (*PostBatchResponse, error) {
	out := new(PostBatchResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/PostBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	AddByText(context.Context, *StringForm) (*CommonResponse, error)
	GetByHashURL(context.Context, *StringForm) (*CommonResponse, error)
	Ping(context.Context, *NoParam) (*IntForm, error)
	Stats(context.Context, *NoParam) (*StatsResponse, error)
	Delete(context.Context, *DeleteRequest) (*IntForm, error)
	GetUserURLs(context.Context, *NoParam) (*GetUserURLsResponse, error)
	PostJSON(context.Context, *PostJSONRespReq) (*PostJSONRespReq, error)
	PostBatch(context.Context, *PostBatchRequest) (*PostBatchResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) AddByText(context.Context, *StringForm) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddByText not implemented")
}
func (UnimplementedShortenerServer) GetByHashURL(context.Context, *StringForm) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByHashURL not implemented")
}
func (UnimplementedShortenerServer) Ping(context.Context, *NoParam) (*IntForm, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortenerServer) Stats(context.Context, *NoParam) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedShortenerServer) Delete(context.Context, *DeleteRequest) (*IntForm, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedShortenerServer) GetUserURLs(context.Context, *NoParam) (*GetUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
}
func (UnimplementedShortenerServer) PostJSON(context.Context, *PostJSONRespReq) (*PostJSONRespReq, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostJSON not implemented")
}
func (UnimplementedShortenerServer) PostBatch(context.Context, *PostBatchRequest) (*PostBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostBatch not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_AddByText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringForm)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).AddByText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/AddByText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).AddByText(ctx, req.(*StringForm))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetByHashURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringForm)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetByHashURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetByHashURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetByHashURL(ctx, req.(*StringForm))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Ping(ctx, req.(*NoParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/Stats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Stats(ctx, req.(*NoParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetUserURLs(ctx, req.(*NoParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_PostJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostJSONRespReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).PostJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/PostJSON",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).PostJSON(ctx, req.(*PostJSONRespReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_PostBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).PostBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/PostBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).PostBatch(ctx, req.(*PostBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddByText",
			Handler:    _Shortener_AddByText_Handler,
		},
		{
			MethodName: "GetByHashURL",
			Handler:    _Shortener_GetByHashURL_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Shortener_Ping_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Shortener_Stats_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Shortener_Delete_Handler,
		},
		{
			MethodName: "GetUserURLs",
			Handler:    _Shortener_GetUserURLs_Handler,
		},
		{
			MethodName: "PostJSON",
			Handler:    _Shortener_PostJSON_Handler,
		},
		{
			MethodName: "PostBatch",
			Handler:    _Shortener_PostBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/proto.proto",
}