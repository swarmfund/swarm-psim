package q

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/notificator-server/auth"
)

type AuthQInterface interface {
	Insert(pair *auth.Pair) error
	ByPublic(token string) (*auth.Pair, error)
}

type AuthQ struct {
	db *sqlx.DB
}

func (q *AuthQ) Insert(pair *auth.Pair) error {
	_, err := q.db.NamedExec(`
		insert into pairs (public, secret) values (:public, :secret)
	`, pair)
	return err
}

func (q *AuthQ) ByPublic(token string) (*auth.Pair, error) {
	result := new(auth.Pair)
	err := q.db.Get(result, `
		select * from pairs
		where public = $1
	`, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}
