package dto

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Status    string    `db:"status"`
	Amount    int       `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateOrderRequest struct {
	UserID uuid.UUID `json:"userId"`
	Status string    `json:"status"`
	Amount int       `json:"amount"`
}

type CreateOrderResponse struct {
	Order Order `json:"order"`
}
