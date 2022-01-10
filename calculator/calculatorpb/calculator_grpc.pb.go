// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package calculatorpb

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

// CalculatorClient is the client API for Calculator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalculatorClient interface {
	// Unary
	Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error)
	// Server Streaming
	PrimeNumberDecomposition(ctx context.Context, in *PrimeNumberDecompositionRequest, opts ...grpc.CallOption) (Calculator_PrimeNumberDecompositionClient, error)
	// Client Streaming
	ComputeAverage(ctx context.Context, opts ...grpc.CallOption) (Calculator_ComputeAverageClient, error)
	// Bidirectional streaming
	FindMaximum(ctx context.Context, opts ...grpc.CallOption) (Calculator_FindMaximumClient, error)
	// this RPC throw an INVALID_ARGUMENT if number is negative
	SquareRoot(ctx context.Context, in *SquareRootRequest, opts ...grpc.CallOption) (*SquareRootResponse, error)
}

type calculatorClient struct {
	cc grpc.ClientConnInterface
}

func NewCalculatorClient(cc grpc.ClientConnInterface) CalculatorClient {
	return &calculatorClient{cc}
}

func (c *calculatorClient) Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error) {
	out := new(SumResponse)
	err := c.cc.Invoke(ctx, "/calculator.Calculator/Sum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calculatorClient) PrimeNumberDecomposition(ctx context.Context, in *PrimeNumberDecompositionRequest, opts ...grpc.CallOption) (Calculator_PrimeNumberDecompositionClient, error) {
	stream, err := c.cc.NewStream(ctx, &Calculator_ServiceDesc.Streams[0], "/calculator.Calculator/PrimeNumberDecomposition", opts...)
	if err != nil {
		return nil, err
	}
	x := &calculatorPrimeNumberDecompositionClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Calculator_PrimeNumberDecompositionClient interface {
	Recv() (*PrimeNumberDecompositionResponse, error)
	grpc.ClientStream
}

type calculatorPrimeNumberDecompositionClient struct {
	grpc.ClientStream
}

func (x *calculatorPrimeNumberDecompositionClient) Recv() (*PrimeNumberDecompositionResponse, error) {
	m := new(PrimeNumberDecompositionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *calculatorClient) ComputeAverage(ctx context.Context, opts ...grpc.CallOption) (Calculator_ComputeAverageClient, error) {
	stream, err := c.cc.NewStream(ctx, &Calculator_ServiceDesc.Streams[1], "/calculator.Calculator/ComputeAverage", opts...)
	if err != nil {
		return nil, err
	}
	x := &calculatorComputeAverageClient{stream}
	return x, nil
}

type Calculator_ComputeAverageClient interface {
	Send(*AverageRequest) error
	CloseAndRecv() (*AverageResponse, error)
	grpc.ClientStream
}

type calculatorComputeAverageClient struct {
	grpc.ClientStream
}

func (x *calculatorComputeAverageClient) Send(m *AverageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *calculatorComputeAverageClient) CloseAndRecv() (*AverageResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(AverageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *calculatorClient) FindMaximum(ctx context.Context, opts ...grpc.CallOption) (Calculator_FindMaximumClient, error) {
	stream, err := c.cc.NewStream(ctx, &Calculator_ServiceDesc.Streams[2], "/calculator.Calculator/FindMaximum", opts...)
	if err != nil {
		return nil, err
	}
	x := &calculatorFindMaximumClient{stream}
	return x, nil
}

type Calculator_FindMaximumClient interface {
	Send(*MaximumRequest) error
	Recv() (*MaximumResponse, error)
	grpc.ClientStream
}

type calculatorFindMaximumClient struct {
	grpc.ClientStream
}

func (x *calculatorFindMaximumClient) Send(m *MaximumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *calculatorFindMaximumClient) Recv() (*MaximumResponse, error) {
	m := new(MaximumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *calculatorClient) SquareRoot(ctx context.Context, in *SquareRootRequest, opts ...grpc.CallOption) (*SquareRootResponse, error) {
	out := new(SquareRootResponse)
	err := c.cc.Invoke(ctx, "/calculator.Calculator/SquareRoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalculatorServer is the server API for Calculator service.
// All implementations must embed UnimplementedCalculatorServer
// for forward compatibility
type CalculatorServer interface {
	// Unary
	Sum(context.Context, *SumRequest) (*SumResponse, error)
	// Server Streaming
	PrimeNumberDecomposition(*PrimeNumberDecompositionRequest, Calculator_PrimeNumberDecompositionServer) error
	// Client Streaming
	ComputeAverage(Calculator_ComputeAverageServer) error
	// Bidirectional streaming
	FindMaximum(Calculator_FindMaximumServer) error
	// this RPC throw an INVALID_ARGUMENT if number is negative
	SquareRoot(context.Context, *SquareRootRequest) (*SquareRootResponse, error)
	mustEmbedUnimplementedCalculatorServer()
}

// UnimplementedCalculatorServer must be embedded to have forward compatible implementations.
type UnimplementedCalculatorServer struct {
}

func (UnimplementedCalculatorServer) Sum(context.Context, *SumRequest) (*SumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sum not implemented")
}
func (UnimplementedCalculatorServer) PrimeNumberDecomposition(*PrimeNumberDecompositionRequest, Calculator_PrimeNumberDecompositionServer) error {
	return status.Errorf(codes.Unimplemented, "method PrimeNumberDecomposition not implemented")
}
func (UnimplementedCalculatorServer) ComputeAverage(Calculator_ComputeAverageServer) error {
	return status.Errorf(codes.Unimplemented, "method ComputeAverage not implemented")
}
func (UnimplementedCalculatorServer) FindMaximum(Calculator_FindMaximumServer) error {
	return status.Errorf(codes.Unimplemented, "method FindMaximum not implemented")
}
func (UnimplementedCalculatorServer) SquareRoot(context.Context, *SquareRootRequest) (*SquareRootResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SquareRoot not implemented")
}
func (UnimplementedCalculatorServer) mustEmbedUnimplementedCalculatorServer() {}

// UnsafeCalculatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalculatorServer will
// result in compilation errors.
type UnsafeCalculatorServer interface {
	mustEmbedUnimplementedCalculatorServer()
}

func RegisterCalculatorServer(s grpc.ServiceRegistrar, srv CalculatorServer) {
	s.RegisterService(&Calculator_ServiceDesc, srv)
}

func _Calculator_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalculatorServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calculator.Calculator/Sum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalculatorServer).Sum(ctx, req.(*SumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calculator_PrimeNumberDecomposition_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PrimeNumberDecompositionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CalculatorServer).PrimeNumberDecomposition(m, &calculatorPrimeNumberDecompositionServer{stream})
}

type Calculator_PrimeNumberDecompositionServer interface {
	Send(*PrimeNumberDecompositionResponse) error
	grpc.ServerStream
}

type calculatorPrimeNumberDecompositionServer struct {
	grpc.ServerStream
}

func (x *calculatorPrimeNumberDecompositionServer) Send(m *PrimeNumberDecompositionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Calculator_ComputeAverage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CalculatorServer).ComputeAverage(&calculatorComputeAverageServer{stream})
}

type Calculator_ComputeAverageServer interface {
	SendAndClose(*AverageResponse) error
	Recv() (*AverageRequest, error)
	grpc.ServerStream
}

type calculatorComputeAverageServer struct {
	grpc.ServerStream
}

func (x *calculatorComputeAverageServer) SendAndClose(m *AverageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *calculatorComputeAverageServer) Recv() (*AverageRequest, error) {
	m := new(AverageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Calculator_FindMaximum_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CalculatorServer).FindMaximum(&calculatorFindMaximumServer{stream})
}

type Calculator_FindMaximumServer interface {
	Send(*MaximumResponse) error
	Recv() (*MaximumRequest, error)
	grpc.ServerStream
}

type calculatorFindMaximumServer struct {
	grpc.ServerStream
}

func (x *calculatorFindMaximumServer) Send(m *MaximumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *calculatorFindMaximumServer) Recv() (*MaximumRequest, error) {
	m := new(MaximumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Calculator_SquareRoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SquareRootRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalculatorServer).SquareRoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calculator.Calculator/SquareRoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalculatorServer).SquareRoot(ctx, req.(*SquareRootRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Calculator_ServiceDesc is the grpc.ServiceDesc for Calculator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calculator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "calculator.Calculator",
	HandlerType: (*CalculatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sum",
			Handler:    _Calculator_Sum_Handler,
		},
		{
			MethodName: "SquareRoot",
			Handler:    _Calculator_SquareRoot_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PrimeNumberDecomposition",
			Handler:       _Calculator_PrimeNumberDecomposition_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ComputeAverage",
			Handler:       _Calculator_ComputeAverage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "FindMaximum",
			Handler:       _Calculator_FindMaximum_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "calculator/calculatorpb/calculator.proto",
}