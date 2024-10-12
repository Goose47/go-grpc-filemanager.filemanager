package interceptors

import (
	"context"
	"filemanager/internal/lib/semaphore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewUnaryLimiter(maxConnections int) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	sem := semaphore.New(maxConnections)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !sem.Acquire() {
			return nil, status.Error(codes.FailedPrecondition, "unary connection limit exceeded")
		}
		defer sem.Release()

		return handler(ctx, req)
	}
}

func NewStreamLimiter(maxConnections int) func(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	sem := semaphore.New(maxConnections)

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if !sem.Acquire() {
			return status.Error(codes.FailedPrecondition, "unary connection limit exceeded")
		}
		defer sem.Release()

		return handler(srv, ss)
	}
}
