package middleware

import (
	"context"

	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC error", "method", info.FullMethod, "code", st.Code(), "message", st.Message())
		}

		return resp, err
	}
}
