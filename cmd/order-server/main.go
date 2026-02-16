package main

import (
	"github.com/SmaF1-dev/grpc_MicroServices/internal/order"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/config"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
)

func main() {
	port := config.GetEnvAsInt("ORDER_PORT", 50052)
	paymentAddr := config.GetEnv("PAYMENT_SERVICE_ADDR", "localhost:50051")

	if err := order.Run(port, paymentAddr); err != nil {
		logger.Fatal("Order service failed: %v", err)
	}
}
