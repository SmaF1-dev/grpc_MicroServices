package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SmaF1-dev/grpc_MicroServices/internal/model"
)

type OrderRepository interface {
	Save(order *model.Order) error
	UpdateStatus(id, status, transactionID string) error
	FindByID(id string) (*model.Order, error)
}

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Save(order *model.Order) error {
	query := `
	INSERT INTO orders (id, user_id, items, total_amount, payment_method, status, transaction_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	_, err = r.db.Exec(query,
		order.ID,
		order.UserID,
		itemsJSON,
		order.TotalAmount,
		order.PaymentMethod,
		order.Status,
		order.TransactionID,
		order.CreatedAt,
		order.UpdatedAt,
	)
	return err
}

func (r *PostgresOrderRepository) UpdateStatus(id, status, transactionID string) error {
	query := `UPDATE orders SET status = $1, transaction_id = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, status, transactionID, time.Now(), id)
	return err
}

func (r *PostgresOrderRepository) FindByID(id string) (*model.Order, error) {
	var order model.Order
	var itemsJSON []byte
	query := `SELECT id, user_id, items, total_amount, payment_method, status, transaction_id, created_at, updated_at
	FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&itemsJSON,
		&order.TotalAmount,
		&order.PaymentMethod,
		&order.Status,
		&order.TransactionID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}
	return &order, nil
}
