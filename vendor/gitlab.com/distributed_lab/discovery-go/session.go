package discovery

import (
	"time"

	consul "github.com/hashicorp/consul/api"
)

type Session struct {
	ID      string
	client  *Client
	release chan bool
}

func NewSession(client *Client) (*Session, error) {
	sid, _, err := client.consul.Session().Create(&consul.SessionEntry{
		LockDelay: 0,
		TTL:       "10s",
	}, nil)
	if err != nil {
		return nil, err
	}
	session := Session{
		client:  client,
		ID:      sid,
		release: make(chan bool),
	}
	return &session, nil
}

func (s *Session) EndlessRenew() {
	go func() {
		for {
			select {
			case <-time.NewTicker(5 * time.Second).C:
				// TODO multiplex errors to client
				s.client.consul.Session().Renew(s.ID, nil)
			case <-s.release:
				return
			}
		}
	}()
}

func (s *Session) Release() {
	s.release <- true
}