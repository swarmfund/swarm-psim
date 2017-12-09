package discovery

import (
	consul "github.com/hashicorp/consul/api"
)

type KVPair struct {
	Key         string
	Value       []byte
	Session     *Session
	ModifyIndex uint64
}

func NewKV(client *Client, ckv *consul.KVPair) *KVPair {
	return &KVPair{
		Key:         ckv.Key,
		Value:       ckv.Value,
		ModifyIndex: ckv.ModifyIndex,
	}
}

func (kv *KVPair) Release() {
	kv.Session.Release()
}

func (kv *KVPair) toConsul() *consul.KVPair {
	ckv := &consul.KVPair{
		Key:         kv.Key,
		Value:       kv.Value,
		ModifyIndex: kv.ModifyIndex,
	}
	if kv.Session != nil {
		ckv.Session = kv.Session.ID
	}
	return ckv
}
