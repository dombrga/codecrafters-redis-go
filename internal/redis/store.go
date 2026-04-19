package redis

import (
	"errors"
	"time"
)

type RedisEntryKey string
type RedisValueType string

type RedisEntryValue struct {
	value     RedisValueType
	expiresAt time.Time
}

type RedisStore struct {
	Items map[RedisEntryKey]RedisEntryValue
}

func (store *RedisStore) Set(key string, value string, ttl time.Duration) {
	entry := RedisEntryValue{
		value: RedisValueType(value),
	}

	if ttl > 0 {
		entry.expiresAt = time.Now().Add(ttl)
	}

	store.Items[RedisEntryKey(key)] = entry
}

func (store *RedisStore) Get(key string) (RedisValueType, error) {
	v, ok := store.Items[RedisEntryKey(key)]
	if !ok {
		return "", errors.New("key does not exist")
	}

	return v.value, nil
}
