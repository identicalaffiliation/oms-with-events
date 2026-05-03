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
	Payload   sql.RawBytes `db:"payload"`
	SentAt    *time.Time   `db:"sent_at"`
	CreatedAt time.Time    `db:"created_at"`
}
