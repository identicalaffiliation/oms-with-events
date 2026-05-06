package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OrderEvent struct {
	ID        uuid.UUID    `db:"id"`
	OrderID   uuid.UUID    `db:"order_id"`
	EventType string       `db:"event_type"`
	Payload   sql.RawBytes `db:"payload" json:"-"`
	SentAt    *time.Time   `db:"sent_at"`
	CreatedAt time.Time    `db:"created_at"`
}

type DispatcherEvent struct {
	EventID uuid.UUID `json:"id"`
	Payload Payload   `json:"payload"`
}

type Payload struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  Status    `json:"status"`
	Amount  int       `json:"amount"`
}
