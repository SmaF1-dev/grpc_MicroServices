package main

import (
	"github.com/SmaF1-dev/grpc_MicroServices/internal/order"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
)

func main() {
	if err := order.Run(); err != nil {
		logger.Fatal("Order service failed: %v", err)
	}
}
