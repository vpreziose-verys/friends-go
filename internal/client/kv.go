package client

import (
	"time"

	"github.com/BethesdaNet/friends-go/internal/platform"
)

type KV interface {
	Name() string
	Value() interface{}
}

type KVTimed interface {
	KV
	Expire() time.Time
	SetExpire(t time.Time) KV
}

type kv struct {
	k      string
	v      platform.Key
	expire time.Time
}

func (k *kv) Name() string { return k.k }
func (k *kv) Value() interface{} {
	if k == nil {
		return nil
	}
	return &k.v
}
func (k *kv) Expire() time.Time       { return k.expire }
func (k kv) SetExpire(t time.Time) KV { k.expire = t; return &k }
