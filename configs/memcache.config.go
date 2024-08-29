package configs

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheClient struct {
	client *memcache.Client
}

var Memcache *MemcacheClient

func InitMemeCache(server string) {
	Memcache = &MemcacheClient{client: memcache.New(server)}
}

func (m *MemcacheClient) Set(key string, value []byte, expiration time.Duration) error {
	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(expiration.Seconds()),
	})
}

func (m *MemcacheClient) Get(key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (m *MemcacheClient) Delete(key string) error {
	return m.client.Delete(key)
}

func (m *MemcacheClient) FlushAll() error {
	return m.client.FlushAll()
}
