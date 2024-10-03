// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: api/v1/collector_server.proto

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

const (
	CollectionInternalService_CreateUsage_FullMethodName = "/llmariner.apiusage.server.v1.CollectionInternalService/CreateUsage"
)

// CollectionInternalServiceClient is the client API for CollectionInternalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CollectionInternalServiceClient interface {
	CreateUsage(ctx context.Context, in *CreateUsageRequest, opts ...grpc.CallOption) (*Usage, error)
}

type collectionInternalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCollectionInternalServiceClient(cc grpc.ClientConnInterface) CollectionInternalServiceClient {
	return &collectionInternalServiceClient{cc}
}

func (c *collectionInternalServiceClient) CreateUsage(ctx context.Context, in *CreateUsageRequest, opts ...grpc.CallOption) (*Usage, error) {
	out := new(Usage)
	err := c.cc.Invoke(ctx, CollectionInternalService_CreateUsage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CollectionInternalServiceServer is the server API for CollectionInternalService service.
// All implementations must embed UnimplementedCollectionInternalServiceServer
// for forward compatibility
type CollectionInternalServiceServer interface {
	CreateUsage(context.Context, *CreateUsageRequest) (*Usage, error)
	mustEmbedUnimplementedCollectionInternalServiceServer()
}

// UnimplementedCollectionInternalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCollectionInternalServiceServer struct {
}

func (UnimplementedCollectionInternalServiceServer) CreateUsage(context.Context, *CreateUsageRequest) (*Usage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUsage not implemented")
}
func (UnimplementedCollectionInternalServiceServer) mustEmbedUnimplementedCollectionInternalServiceServer() {
}

// UnsafeCollectionInternalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CollectionInternalServiceServer will
// result in compilation errors.
type UnsafeCollectionInternalServiceServer interface {
	mustEmbedUnimplementedCollectionInternalServiceServer()
}

func RegisterCollectionInternalServiceServer(s grpc.ServiceRegistrar, srv CollectionInternalServiceServer) {
	s.RegisterService(&CollectionInternalService_ServiceDesc, srv)
}

func _CollectionInternalService_CreateUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUsageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionInternalServiceServer).CreateUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectionInternalService_CreateUsage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionInternalServiceServer).CreateUsage(ctx, req.(*CreateUsageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CollectionInternalService_ServiceDesc is the grpc.ServiceDesc for CollectionInternalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CollectionInternalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "llmariner.apiusage.server.v1.CollectionInternalService",
	HandlerType: (*CollectionInternalServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUsage",
			Handler:    _CollectionInternalService_CreateUsage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/collector_server.proto",
}
