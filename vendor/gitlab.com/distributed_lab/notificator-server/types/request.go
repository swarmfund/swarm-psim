package types

import (
	"time"
	"github.com/jmoiron/sqlx/types"
)

type Request struct {
	ID        int64          `db:"id"`
	Type      RequestTypeID  `db:"type"`
	Payload   types.JSONText `db:"payload"`
	Priority  int            `db:"priority"`
	Token     string         `db:"token"`
	Completed *time.Time     `db:"completed_at"`
	CreatedAt time.Time      `db:"created_at"`
	Hash      string         `db:"hash"`
}

func NewRequest(id RequestTypeID, priority int, payload string, token string, hash string) *Request {
	return &Request{
		Type:     id,
		Priority: priority,
		Payload:  types.JSONText(payload),
		Token:    token,
		Hash:     hash,
	}
}

