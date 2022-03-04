package model

import (
	"time"

	"github.com/Kolakanmi/grey_wallet/pkg/uuid"
)

type Base struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (b *Base) SetID() {
	b.ID = uuid.New()
}
