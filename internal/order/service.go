package order

import (
	"context"
	"time"

	order "github.com/SmaF1-dev/grpc_MicroServices/api/order"
	payment "github.com/SmaF1-dev/grpc_MicroServices/api/payment"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderServiceServer struct {
	order.UnimplementedOrderServiceServer
	paymentClient payment.PaymentServiceClient
	conn          *grpc.ClientConn
}

func NewOrderServiceServer(paymentAddr string) (*OrderServiceServer, error) {
	conn, err := grpc.NewClient(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	paymentClient := payment.NewPaymentServiceClient(conn)

	return &OrderServiceServer{
		paymentClient: paymentClient,
		conn:          conn,
	}, nil
}

func (s *OrderServiceServer) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	logger.Info("Creating order: %+v", req)

	// Упрощенная валидация

	if req.OrderId == "" || req.TotalAmount <= 0 {
		return &order.CreateOrderResponse{
			OrderId: req.OrderId,
			Status:  "failed",
			Message: "Invalid order data",
		}, nil
	}

	ctxPayment, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	paymentResp, err := s.paymentClient.ProcessPayment(ctxPayment, &payment.ProcessPaymentRequest{
		OrderId: req.OrderId,
		Amount:  req.TotalAmount,
		Method:  req.PaymentMethod,
	})

	if err != nil {
		logger.Error("Payment service call failed: %v", err)
		return &order.CreateOrderResponse{
			OrderId: req.OrderId,
			Status:  "failed",
			Message: "payment service error: " + err.Error(),
		}, nil
	}

	if !paymentResp.Success {
		logger.Info("Payment rejected for order %s: %s", req.OrderId, paymentResp.Message)
		return &order.CreateOrderResponse{
			OrderId: req.OrderId,
			Status:  "failed",
			Message: "payment failed: " + paymentResp.Message,
		}, nil
	}

	logger.Info("Order %s created succesfully, transaction: %s", req.OrderId, paymentResp.TransactionId)
	return &order.CreateOrderResponse{
		OrderId: req.OrderId,
		Status:  "confirmed",
		Message: "order created and payment succeeded",
	}, nil

}
