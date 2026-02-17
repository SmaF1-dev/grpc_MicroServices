package interceptor

import (
	"context"
	"time"

	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		st, _ := status.FromError(err)
		logger.Error("gRPC call %s falled: %v (code=%s, took=%v)", info.FullMethod, err, st.Code(), duration)
	} else {
		logger.Info("gRPC call %s succeeded (took=%v)", info.FullMethod, duration)
	}
	return resp, err
}
