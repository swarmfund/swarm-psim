package q

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/notificator-server/types"
	"time"
)

type RequestQInterface interface {
	Insert(request *types.Request) error
	ByHash(hash string) (*types.Request, error)
	ByToken(token string) (*types.Request, error)
	NextWindow(requestType types.RequestTypeID, token string, limit int, from time.Duration) (*time.Duration, error)
	GetHead() ([]types.Request, error)
	LowerPriority(id int64) error
	MarkCompleted(id int64) error
}

type RequestQ struct {
	db *sqlx.DB
}

func (r *RequestQ) Insert(request *types.Request) error {
	_, err := r.db.NamedExec(`
		insert into requests (type, payload, priority, token, hash)
		values (:type, :payload, :priority, :token, :hash)`, request)
	return err
}

func (r *RequestQ) ByHash(hash string) (*types.Request, error) {
	result := new(types.Request)
	err := r.db.Get(result, `
		select * from requests
		where hash = $1`, hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (r *RequestQ) ByToken(token string) (*types.Request, error) {
	result := new(types.Request)
	err := r.db.Get(result, `
		select * from requests
		where token = $1`, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (r *RequestQ) NextWindow(requestType types.RequestTypeID, token string, limit int, interval time.Duration) (*time.Duration, error) {
	windowStart := new(time.Time)
	err := r.db.Get(&windowStart, `
		  select created_at from requests
		  where created_at > (timestamp 'now' - $1*interval '1')
		    and token = $2
		    and type = $3
		  order by created_at desc
		  offset $4
		  limit 1`, interval.Seconds(), token, requestType, limit-1)
	delta := time.Now().Sub(*windowStart)
	if err != nil || windowStart.Nanosecond() == 0 || delta > interval {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	result := interval - delta
	return &result, err
}

func (r *RequestQ) GetHead() ([]types.Request, error) {
	var result []types.Request
	err := r.db.Select(&result, `
		select * from requests
		where completed_at is null
		order by priority desc
		limit 100`)
	return result, err
}

func (r *RequestQ) LowerPriority(id int64) error {
	_, err := r.db.Exec(`
		update requests
		set priority = priority - 1
		where id = $1`, id)
	return err
}

func (r *RequestQ) MarkCompleted(id int64) error {
	_, err := r.db.Exec(`
		update requests
		set completed_at = timestamp 'now'
		where id = $1`, id)
	return err
}
