package model

import "time"

type Order struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	Items         []string  `db:"items" json:"items"`
	TotalAmount   float64   `db:"total_amount" json:"total_amount"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	Status        string    `db:"status" json:"status"` // pending, confirmed, failed
	TransactionID string    `db:"transaction_id" json:"transaction_id,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
