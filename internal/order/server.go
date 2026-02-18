package order

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	order "github.com/SmaF1-dev/grpc_MicroServices/api/order"
	"github.com/SmaF1-dev/grpc_MicroServices/internal/interceptor"
	"github.com/SmaF1-dev/grpc_MicroServices/internal/repository"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/config"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
	"google.golang.org/grpc"
)

func Run() error {
	port := config.GetEnvAsInt("ORDER_PORT", 50052)
	paymentAddr := config.GetEnv("PAYMENT_SERVICE_ADDR", "localhost:50051")

	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbUser := config.GetEnv("DB_USER", "postgres")
	dbPass := config.GetEnv("DB_PASSWORD", "postgres")
	dbName := config.GetEnv("DB_NAME", "orders")
	dbSSL := config.GetEnv("DB_SSLMODE", "disable")

	db, err := repository.NewDB(dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	repo := repository.NewPostgresOrderRepository(db)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
	)

	orderServiceServer, err := NewOrderServiceServer(paymentAddr, repo)
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

	logger.Info("Total succesfully processed orders: %d", orderServiceServer.GetSuccessCount())

	if err := orderServiceServer.Close(); err != nil {
		logger.Error("Error closing payment client connection: %v", err)
	}

	logger.Info("Order Service stopped")
	return nil
}
