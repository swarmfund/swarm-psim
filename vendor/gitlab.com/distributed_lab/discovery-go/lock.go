package discovery

import (
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
)

// Lock abstraction on top of consul KV lock allowing to set lock key value
// with nicer API
type Lock struct {
	client *Client
	ctx    context.Context
	lockKey    string
	valueKey string
	lock   *api.Lock
}

func (c *Client) Lock(ctx context.Context, key string) *Lock {
	return &Lock{
		client: c,
		ctx:    ctx,
		valueKey:    key,
		lockKey: fmt.Sprintf("%s/lock", key),
	}
}

func (l *Lock) Lock() (<-chan struct{}, error) {
	opts := api.LockOptions{
		Key:         l.lockKey,
		LockTryOnce: true,
	}
	lock, err := l.client.consul.LockOpts(&opts)
	if err != nil {
		return nil, err
	}
	l.lock = lock
	stopCh := make(chan struct{})
	go func() {
		<-l.ctx.Done()
		close(stopCh)
	}()
	return lock.Lock(stopCh)
}

func (l *Lock) Unlock() error {
	return l.lock.Unlock()
}

func (l *Lock) Get() ([]byte, error){
	kv, err := l.client.Get(l.valueKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key")
	}

	if kv == nil || kv.Value == nil {
		return nil, nil
	}

	return kv.Value, nil
}

func (l *Lock) Set(value interface{}) (err error) {
	bytes, ok := value.([]byte)
	if !ok {
		bytes, err = json.Marshal(value)
		if err != nil {
			return errors.Wrap(err, "failed to marshal value")
		}
	}
	err = l.client.Set(KVPair{
		Key: l.valueKey,
		Value: bytes,
	})
	if err != nil {
		return errors.Wrap(err, "failed to")
	}
	return nil
}
