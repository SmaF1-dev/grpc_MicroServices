package order

import (
	"context"
	"testing"

	order "github.com/SmaF1-dev/grpc_MicroServices/api/order"
	payment "github.com/SmaF1-dev/grpc_MicroServices/api/payment"
	"github.com/SmaF1-dev/grpc_MicroServices/internal/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockRepository struct {
	SaveFunc         func(order *model.Order) error
	UpdateStatusFunc func(id, status, transactionID string) error
	FindByIDFunc     func(id string) (*model.Order, error)
}

func (m *mockRepository) Save(order *model.Order) error {
	return m.SaveFunc(order)
}

func (m *mockRepository) UpdateStatus(id, status, transactionID string) error {
	return m.UpdateStatusFunc(id, status, transactionID)
}

func (m *mockRepository) FindByID(id string) (*model.Order, error) {
	return m.FindByIDFunc(id)
}

type mockPaymentClient struct {
	ProcessPaymentFunc func(ctx context.Context, req *payment.ProcessPaymentRequest, opts ...grpc.CallOption) (*payment.ProcessPaymentResponse, error)
}

func (m *mockPaymentClient) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentRequest, opts ...grpc.CallOption) (*payment.ProcessPaymentResponse, error) {
	return m.ProcessPaymentFunc(ctx, req)
}

func TestOrderServiceServer_CreateOrder_Success(t *testing.T) {
	mockRepo := &mockRepository{
		SaveFunc: func(order *model.Order) error {
			assert.Equal(t, "order-1", order.ID)
			assert.Equal(t, "pending", order.Status)
			return nil
		},
		UpdateStatusFunc: func(id, status, transactionID string) error {
			assert.Equal(t, "order-1", id)
			assert.Equal(t, "confirmed", status)
			assert.NotEmpty(t, transactionID)
			return nil
		},
	}

	mockPay := &mockPaymentClient{
		ProcessPaymentFunc: func(ctx context.Context, req *payment.ProcessPaymentRequest, opts ...grpc.CallOption) (*payment.ProcessPaymentResponse, error) {
			assert.Equal(t, "order-1", req.OrderId)
			assert.Equal(t, 100.0, req.Amount)
			return &payment.ProcessPaymentResponse{
				Success:       true,
				TransactionId: "txn-123",
				Message:       "OK",
			}, nil
		},
	}

	srv := &OrderServiceServer{
		paymentClient: mockPay,
		repo:          mockRepo,
	}

	req := &order.CreateOrderRequest{
		OrderId:       "order-1",
		UserId:        "user-1",
		Items:         []string{"item1"},
		TotalAmount:   100.0,
		PaymentMethod: "card",
	}

	resp, err := srv.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "confirmed", resp.Status)
}

func TestOrderServiceServer_CreateOrder_PaymentFailed(t *testing.T) {
	mockRepo := &mockRepository{
		SaveFunc: func(order *model.Order) error {
			return nil
		},
		UpdateStatusFunc: func(id, status, transactionID string) error {
			assert.Equal(t, "order-2", id)
			assert.Equal(t, "failed", status)
			return nil
		},
	}

	mockPay := &mockPaymentClient{
		ProcessPaymentFunc: func(ctx context.Context, req *payment.ProcessPaymentRequest, opts ...grpc.CallOption) (*payment.ProcessPaymentResponse, error) {
			return &payment.ProcessPaymentResponse{
				Success: false,
				Message: "Insufficient funds",
			}, nil
		},
	}

	srv := &OrderServiceServer{
		paymentClient: mockPay,
		repo:          mockRepo,
	}

	req := &order.CreateOrderRequest{
		OrderId:       "order-2",
		UserId:        "user-2",
		Items:         []string{"item1"},
		TotalAmount:   100.0,
		PaymentMethod: "card",
	}

	resp, err := srv.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "failed", resp.Status)
}

func TestOrderServiceServer_CreateOrder_SaveError(t *testing.T) {
	mockRepo := &mockRepository{
		SaveFunc: func(order *model.Order) error {
			return assert.AnError
		}}

	mockPay := &mockPaymentClient{}

	srv := &OrderServiceServer{
		paymentClient: mockPay,
		repo:          mockRepo,
	}

	req := &order.CreateOrderRequest{
		OrderId:       "order-3",
		UserId:        "user-3",
		Items:         []string{"item1"},
		TotalAmount:   100.0,
		PaymentMethod: "card",
	}

	resp, err := srv.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Nil(t, resp)
}
