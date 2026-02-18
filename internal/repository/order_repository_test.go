package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SmaF1-dev/grpc_MicroServices/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPostgresOrderRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	now := time.Now()
	order := model.Order{
		ID:            "test-123",
		UserID:        "user-1",
		Items:         []string{"item1", "item2"},
		TotalAmount:   99.99,
		PaymentMethod: "card",
		Status:        "pending",
		TransactionID: "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	itemsJSON, _ := json.Marshal(order.Items)

	mock.ExpectExec(`INSERT INTO orders`).WithArgs(
		order.ID,
		order.UserID,
		itemsJSON,
		order.TotalAmount,
		order.PaymentMethod,
		order.Status,
		order.TransactionID,
		order.CreatedAt,
		order.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(&order)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOrderRepository_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	orderID := "test-123"
	status := "confirmed"
	txnID := "txn-456"

	mock.ExpectExec(`UPDATE orders SET status = \$1, transaction_id = \$2, updated_at = \$3 WHERE id = \$4`).WithArgs(
		status, txnID, sqlmock.AnyArg(), orderID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateStatus(orderID, status, txnID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOrderRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresOrderRepository(db)

	orderID := "test-123"
	items := []string{"item1", "item2"}
	itemsJSON, _ := json.Marshal(items)
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "items", "total_amount", "payment_method",
		"status", "transaction_id", "created_at", "updated_at",
	}).AddRow(
		orderID,
		"user-1",
		itemsJSON,
		99.99,
		"card",
		"confirmed",
		"txn-456",
		now,
		now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM orders WHERE id = \$1`).WithArgs(orderID).WillReturnRows(rows)

	order, err := repo.FindByID(orderID)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orderID, order.ID)
	assert.Equal(t, items, order.Items)
	assert.NoError(t, mock.ExpectationsWereMet())
}
