package main

import (
	"github.com/SmaF1-dev/grpc_MicroServices/internal/payment"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/config"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
)

func main() {
	port := config.GetEnvAsInt("PAYMENT_PORT", 50051)

	if err := payment.Run(port); err != nil {
		logger.Fatal("Payment service failed: %v", err)
	}
}
