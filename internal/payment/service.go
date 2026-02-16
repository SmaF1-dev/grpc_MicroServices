package payment

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	payment "github.com/SmaF1-dev/grpc_MicroServices/api/payment"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
)

type PaymentServiceServer struct {
	payment.UnimplementedPaymentServiceServer
}

func NewPaymentServiceServer() *PaymentServiceServer {
	return &PaymentServiceServer{}
}

func (s *PaymentServiceServer) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentRequest) (*payment.ProcessPaymentResponse, error) {
	logger.Info("Processing payment for order %s, amount %.2f, method %s", req.OrderId, req.Amount, req.Method)

	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		logger.Error("Payment processing cancelled for order %s", req.OrderId)
		return nil, ctx.Err()
	}

	success := rand.Float32() <= 0.8

	if success {
		txnID := fmt.Sprintf("txn_%d", time.Now().UnixNano())
		logger.Info("Payment successful for order %s, txn=%s", req.OrderId, txnID)
		return &payment.ProcessPaymentResponse{
			Success:       true,
			TransactionId: txnID,
			Message:       "Payment processed successfully",
		}, nil
	}

	logger.Error("Payment failed for order %s", req.OrderId)
	return &payment.ProcessPaymentResponse{
		Success: false,
		Message: "Insufficient funds or other error",
	}, nil
}
