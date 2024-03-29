// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: orbis/ring/v1alpha1/ring.proto

package ringv1alpha1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	RingService_ListRings_FullMethodName       = "/orbis.ring.v1alpha1.RingService/ListRings"
	RingService_GetRing_FullMethodName         = "/orbis.ring.v1alpha1.RingService/GetRing"
	RingService_CreateRing_FullMethodName      = "/orbis.ring.v1alpha1.RingService/CreateRing"
	RingService_DeleteRing_FullMethodName      = "/orbis.ring.v1alpha1.RingService/DeleteRing"
	RingService_PublicKey_FullMethodName       = "/orbis.ring.v1alpha1.RingService/PublicKey"
	RingService_Refresh_FullMethodName         = "/orbis.ring.v1alpha1.RingService/Refresh"
	RingService_State_FullMethodName           = "/orbis.ring.v1alpha1.RingService/State"
	RingService_ListSecrets_FullMethodName     = "/orbis.ring.v1alpha1.RingService/ListSecrets"
	RingService_StoreSecret_FullMethodName     = "/orbis.ring.v1alpha1.RingService/StoreSecret"
	RingService_ReencryptSecret_FullMethodName = "/orbis.ring.v1alpha1.RingService/ReencryptSecret"
	RingService_DeleteSecret_FullMethodName    = "/orbis.ring.v1alpha1.RingService/DeleteSecret"
)

// RingServiceClient is the client API for RingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RingServiceClient interface {
	ListRings(ctx context.Context, in *ListRingsRequest, opts ...grpc.CallOption) (*ListRingsResponse, error)
	GetRing(ctx context.Context, in *GetRingRequest, opts ...grpc.CallOption) (*GetRingResponse, error)
	CreateRing(ctx context.Context, in *CreateRingRequest, opts ...grpc.CallOption) (*CreateRingResponse, error)
	DeleteRing(ctx context.Context, in *DeleteRingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PublicKey(ctx context.Context, in *PublicKeyRequest, opts ...grpc.CallOption) (*PublicKeyResponse, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
	State(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*StateResponse, error)
	ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error)
	StoreSecret(ctx context.Context, in *StoreSecretRequest, opts ...grpc.CallOption) (*StoreSecretResponse, error)
	ReencryptSecret(ctx context.Context, in *ReencryptSecretRequest, opts ...grpc.CallOption) (*ReencryptSecretResponse, error)
	DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type ringServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRingServiceClient(cc grpc.ClientConnInterface) RingServiceClient {
	return &ringServiceClient{cc}
}

func (c *ringServiceClient) ListRings(ctx context.Context, in *ListRingsRequest, opts ...grpc.CallOption) (*ListRingsResponse, error) {
	out := new(ListRingsResponse)
	err := c.cc.Invoke(ctx, RingService_ListRings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) GetRing(ctx context.Context, in *GetRingRequest, opts ...grpc.CallOption) (*GetRingResponse, error) {
	out := new(GetRingResponse)
	err := c.cc.Invoke(ctx, RingService_GetRing_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) CreateRing(ctx context.Context, in *CreateRingRequest, opts ...grpc.CallOption) (*CreateRingResponse, error) {
	out := new(CreateRingResponse)
	err := c.cc.Invoke(ctx, RingService_CreateRing_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) DeleteRing(ctx context.Context, in *DeleteRingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RingService_DeleteRing_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) PublicKey(ctx context.Context, in *PublicKeyRequest, opts ...grpc.CallOption) (*PublicKeyResponse, error) {
	out := new(PublicKeyResponse)
	err := c.cc.Invoke(ctx, RingService_PublicKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	out := new(RefreshResponse)
	err := c.cc.Invoke(ctx, RingService_Refresh_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) State(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*StateResponse, error) {
	out := new(StateResponse)
	err := c.cc.Invoke(ctx, RingService_State_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error) {
	out := new(ListSecretsResponse)
	err := c.cc.Invoke(ctx, RingService_ListSecrets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) StoreSecret(ctx context.Context, in *StoreSecretRequest, opts ...grpc.CallOption) (*StoreSecretResponse, error) {
	out := new(StoreSecretResponse)
	err := c.cc.Invoke(ctx, RingService_StoreSecret_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) ReencryptSecret(ctx context.Context, in *ReencryptSecretRequest, opts ...grpc.CallOption) (*ReencryptSecretResponse, error) {
	out := new(ReencryptSecretResponse)
	err := c.cc.Invoke(ctx, RingService_ReencryptSecret_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringServiceClient) DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RingService_DeleteSecret_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RingServiceServer is the server API for RingService service.
// All implementations must embed UnimplementedRingServiceServer
// for forward compatibility
type RingServiceServer interface {
	ListRings(context.Context, *ListRingsRequest) (*ListRingsResponse, error)
	GetRing(context.Context, *GetRingRequest) (*GetRingResponse, error)
	CreateRing(context.Context, *CreateRingRequest) (*CreateRingResponse, error)
	DeleteRing(context.Context, *DeleteRingRequest) (*emptypb.Empty, error)
	PublicKey(context.Context, *PublicKeyRequest) (*PublicKeyResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	State(context.Context, *StateRequest) (*StateResponse, error)
	ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error)
	StoreSecret(context.Context, *StoreSecretRequest) (*StoreSecretResponse, error)
	ReencryptSecret(context.Context, *ReencryptSecretRequest) (*ReencryptSecretResponse, error)
	DeleteSecret(context.Context, *DeleteSecretRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedRingServiceServer()
}

// UnimplementedRingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRingServiceServer struct {
}

func (UnimplementedRingServiceServer) ListRings(context.Context, *ListRingsRequest) (*ListRingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRings not implemented")
}
func (UnimplementedRingServiceServer) GetRing(context.Context, *GetRingRequest) (*GetRingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRing not implemented")
}
func (UnimplementedRingServiceServer) CreateRing(context.Context, *CreateRingRequest) (*CreateRingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRing not implemented")
}
func (UnimplementedRingServiceServer) DeleteRing(context.Context, *DeleteRingRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRing not implemented")
}
func (UnimplementedRingServiceServer) PublicKey(context.Context, *PublicKeyRequest) (*PublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublicKey not implemented")
}
func (UnimplementedRingServiceServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedRingServiceServer) State(context.Context, *StateRequest) (*StateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method State not implemented")
}
func (UnimplementedRingServiceServer) ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecrets not implemented")
}
func (UnimplementedRingServiceServer) StoreSecret(context.Context, *StoreSecretRequest) (*StoreSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreSecret not implemented")
}
func (UnimplementedRingServiceServer) ReencryptSecret(context.Context, *ReencryptSecretRequest) (*ReencryptSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReencryptSecret not implemented")
}
func (UnimplementedRingServiceServer) DeleteSecret(context.Context, *DeleteSecretRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}
func (UnimplementedRingServiceServer) mustEmbedUnimplementedRingServiceServer() {}

// UnsafeRingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RingServiceServer will
// result in compilation errors.
type UnsafeRingServiceServer interface {
	mustEmbedUnimplementedRingServiceServer()
}

func RegisterRingServiceServer(s grpc.ServiceRegistrar, srv RingServiceServer) {
	s.RegisterService(&RingService_ServiceDesc, srv)
}

func _RingService_ListRings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).ListRings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_ListRings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).ListRings(ctx, req.(*ListRingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_GetRing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).GetRing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_GetRing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).GetRing(ctx, req.(*GetRingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_CreateRing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).CreateRing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_CreateRing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).CreateRing(ctx, req.(*CreateRingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_DeleteRing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).DeleteRing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_DeleteRing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).DeleteRing(ctx, req.(*DeleteRingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_PublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).PublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_PublicKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).PublicKey(ctx, req.(*PublicKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_Refresh_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_State_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).State(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_State_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).State(ctx, req.(*StateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_ListSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).ListSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_ListSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).ListSecrets(ctx, req.(*ListSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_StoreSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).StoreSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_StoreSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).StoreSecret(ctx, req.(*StoreSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_ReencryptSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReencryptSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).ReencryptSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_ReencryptSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).ReencryptSecret(ctx, req.(*ReencryptSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingService_DeleteSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServiceServer).DeleteSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingService_DeleteSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServiceServer).DeleteSecret(ctx, req.(*DeleteSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RingService_ServiceDesc is the grpc.ServiceDesc for RingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "orbis.ring.v1alpha1.RingService",
	HandlerType: (*RingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRings",
			Handler:    _RingService_ListRings_Handler,
		},
		{
			MethodName: "GetRing",
			Handler:    _RingService_GetRing_Handler,
		},
		{
			MethodName: "CreateRing",
			Handler:    _RingService_CreateRing_Handler,
		},
		{
			MethodName: "DeleteRing",
			Handler:    _RingService_DeleteRing_Handler,
		},
		{
			MethodName: "PublicKey",
			Handler:    _RingService_PublicKey_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _RingService_Refresh_Handler,
		},
		{
			MethodName: "State",
			Handler:    _RingService_State_Handler,
		},
		{
			MethodName: "ListSecrets",
			Handler:    _RingService_ListSecrets_Handler,
		},
		{
			MethodName: "StoreSecret",
			Handler:    _RingService_StoreSecret_Handler,
		},
		{
			MethodName: "ReencryptSecret",
			Handler:    _RingService_ReencryptSecret_Handler,
		},
		{
			MethodName: "DeleteSecret",
			Handler:    _RingService_DeleteSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "orbis/ring/v1alpha1/ring.proto",
}
