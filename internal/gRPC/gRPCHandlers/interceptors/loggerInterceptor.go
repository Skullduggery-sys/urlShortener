package interceptors

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func LoggerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := logger.WithFields(logrus.Fields{
			"method": info.FullMethod,
		})

		logger.Info("gRPC method called")

		resp, err := handler(ctx, req)

		if err != nil {
			logger.Errorf("gRPC method failed: %v", err)
		} else {
			logger.Info("gRPC method completed successfully")
		}

		return resp, err
	}
}
