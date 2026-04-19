package store

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
		return "$-1\r\n", errors.New("key does not exist")
	}

	// fmt.Println("redis entry:", v, v.expiresAt.Local().Format(time.RFC3339))
	// fmt.Println("time now:", time.Now().Local().Format(time.RFC3339))

	isExpired := store.IsEntryExpired(v)
	if isExpired {
		return "$-1\r\n", errors.New("key expired")
	}

	return v.value, nil
}

func (store *RedisStore) IsEntryExpired(v RedisEntryValue) bool {
	if v.expiresAt.IsZero() {
		return false
	}
	return v.expiresAt.Before(time.Now())
}
