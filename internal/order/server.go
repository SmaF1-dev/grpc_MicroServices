package order

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	order "github.com/SmaF1-dev/grpc_MicroServices/api/order"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
	"google.golang.org/grpc"
)

func Run(port int, paymentAddr string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	orderServiceServer, err := NewOrderServiceServer(paymentAddr)
	if err != nil {
		return fmt.Errorf("failed to create order service: %w", err)
	}
	order.RegisterOrderServiceServer(grpcServer, orderServiceServer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Order Service is running on port %d", port)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve: %v", err)
		}
	}()

	<-quit
	logger.Info("Shutting down Order Service...")
	grpcServer.GracefulStop()

	if err := orderServiceServer.Close(); err != nil {
		logger.Error("Error closing payment client connection: %v", err)
	}

	logger.Info("Order Service stopped")
	return nil
}
