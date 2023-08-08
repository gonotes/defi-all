// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: v1/settings.proto

package v1

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

// SettingsServiceClient is the client API for SettingsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SettingsServiceClient interface {
	ResetSettings(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*ResetResponse, error)
	GetSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*GetSettingsResponse, error)
	UpdateSettings(ctx context.Context, in *UpdateSettingsRequest, opts ...grpc.CallOption) (*UpdateSettingsResponse, error)
}

type settingsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSettingsServiceClient(cc grpc.ClientConnInterface) SettingsServiceClient {
	return &settingsServiceClient{cc}
}

func (c *settingsServiceClient) ResetSettings(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*ResetResponse, error) {
	out := new(ResetResponse)
	err := c.cc.Invoke(ctx, "/settings.SettingsService/ResetSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingsServiceClient) GetSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*GetSettingsResponse, error) {
	out := new(GetSettingsResponse)
	err := c.cc.Invoke(ctx, "/settings.SettingsService/GetSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingsServiceClient) UpdateSettings(ctx context.Context, in *UpdateSettingsRequest, opts ...grpc.CallOption) (*UpdateSettingsResponse, error) {
	out := new(UpdateSettingsResponse)
	err := c.cc.Invoke(ctx, "/settings.SettingsService/UpdateSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SettingsServiceServer is the server API for SettingsService service.
// All implementations must embed UnimplementedSettingsServiceServer
// for forward compatibility
type SettingsServiceServer interface {
	ResetSettings(context.Context, *ResetRequest) (*ResetResponse, error)
	GetSettings(context.Context, *GetSettingsRequest) (*GetSettingsResponse, error)
	UpdateSettings(context.Context, *UpdateSettingsRequest) (*UpdateSettingsResponse, error)
	mustEmbedUnimplementedSettingsServiceServer()
}

// UnimplementedSettingsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSettingsServiceServer struct {
}

func (UnimplementedSettingsServiceServer) ResetSettings(context.Context, *ResetRequest) (*ResetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetSettings not implemented")
}
func (UnimplementedSettingsServiceServer) GetSettings(context.Context, *GetSettingsRequest) (*GetSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSettings not implemented")
}
func (UnimplementedSettingsServiceServer) UpdateSettings(context.Context, *UpdateSettingsRequest) (*UpdateSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSettings not implemented")
}
func (UnimplementedSettingsServiceServer) mustEmbedUnimplementedSettingsServiceServer() {}

// UnsafeSettingsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SettingsServiceServer will
// result in compilation errors.
type UnsafeSettingsServiceServer interface {
	mustEmbedUnimplementedSettingsServiceServer()
}

func RegisterSettingsServiceServer(s grpc.ServiceRegistrar, srv SettingsServiceServer) {
	s.RegisterService(&SettingsService_ServiceDesc, srv)
}

func _SettingsService_ResetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).ResetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/settings.SettingsService/ResetSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).ResetSettings(ctx, req.(*ResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingsService_GetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).GetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/settings.SettingsService/GetSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).GetSettings(ctx, req.(*GetSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingsService_UpdateSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).UpdateSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/settings.SettingsService/UpdateSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).UpdateSettings(ctx, req.(*UpdateSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SettingsService_ServiceDesc is the grpc.ServiceDesc for SettingsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SettingsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "settings.SettingsService",
	HandlerType: (*SettingsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetSettings",
			Handler:    _SettingsService_ResetSettings_Handler,
		},
		{
			MethodName: "GetSettings",
			Handler:    _SettingsService_GetSettings_Handler,
		},
		{
			MethodName: "UpdateSettings",
			Handler:    _SettingsService_UpdateSettings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/settings.proto",
}
