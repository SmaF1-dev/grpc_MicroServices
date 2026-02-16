package payment

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	payment "github.com/SmaF1-dev/grpc_MicroServices/api/payment"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
	"google.golang.org/grpc"
)

func Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	paymentServiceServer := NewPaymentServiceServer()
	payment.RegisterPaymentServiceServer(grpcServer, paymentServiceServer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Payment Service is running on port %d", port)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve: %v", err)
		}
	}()

	<-quit
	logger.Info("Shutting down Payment Service")
	grpcServer.GracefulStop()
	logger.Info("Payment service stopped")
	return nil
}
