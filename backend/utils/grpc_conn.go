package utils

import (
	"context"

	"google.golang.org/grpc"
)

type GrpcConn struct {
	host string
}

func NewGrpcConn() *GrpcConn {
	return &GrpcConn{
		host: "localhost:8080",
	}
}

// Invoke performs a unary RPC and returns after the response is received
// into reply.
func (g *GrpcConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	conn, err := grpc.NewClient(g.host)
	if err != nil {
		return err
	}
	err = conn.Invoke(ctx, method, args, reply, opts...)
	return err
}

// NewStream begins a streaming RPC.
func (g *GrpcConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	conn, err := grpc.NewClient(g.host)
	if err != nil {
		return nil, err
	}
	return conn.NewStream(ctx, desc, method, opts...)
}
