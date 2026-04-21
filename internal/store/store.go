package store

import (
	"errors"
	"time"
)

const NULL_BULK_STRING = "$-1\r\n"

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
		return NULL_BULK_STRING, errors.New("key does not exist")
	}

	isExpired := store.IsEntryExpired(v)
	if isExpired {
		return NULL_BULK_STRING, errors.New("key expired")
	}

	return v.value, nil
}

func (store *RedisStore) IsEntryExpired(v RedisEntryValue) bool {
	if v.expiresAt.IsZero() {
		return false
	}
	return v.expiresAt.Before(time.Now())
}
