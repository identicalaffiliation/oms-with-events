package domain

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Created Status = "created"
	Paid    Status = "paid"
	Shipped Status = "shipped"
)

type Order struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Status    Status    `db:"status"`
	Amount    int       `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (o *Order) ToJSON() sql.RawBytes {
	bytes, _ := json.Marshal(o)

	return bytes
}
