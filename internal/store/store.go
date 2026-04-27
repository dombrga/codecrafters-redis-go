package store

import (
	"errors"
	"strconv"
	"time"
)

const REDIS_NULL_BULK_STRING = "$-1\r\n"

type RedisEntryKey string
type RedisValueType string

// type RedisValueType interface { string | []string | int | float64 }
// type RedisValueType1

type RedisEntryValue struct {
	value     any
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
		return REDIS_NULL_BULK_STRING, errors.New("key does not exist")
	}

	isExpired := store.isEntryExpired(v)
	if isExpired {
		return REDIS_NULL_BULK_STRING, errors.New("key expired")
	}

	return v.value.(RedisValueType), nil
}

// The RPUSH command is used to append elements to a list. If the list doesn't exist, it is created first.
func (store *RedisStore) RPush(key string, values []string) (string, error) {
	list, ok := store.Items[RedisEntryKey(key)]
	if !ok {
		list = RedisEntryValue{
			value: values,
		}

		store.Items[RedisEntryKey(key)] = list
		length := len(list.value.([]string))
		return strconv.Itoa(length), nil
	}

	_list, ok := list.value.([]string)
	if !ok {
		return "", errors.New("not a list")
	}

	list.value = append(_list, values...)
	store.Items[RedisEntryKey(key)] = list

	return strconv.Itoa(len(list.value.([]string))), nil
}

func (store *RedisStore) isEntryExpired(v RedisEntryValue) bool {
	if v.expiresAt.IsZero() {
		return false
	}
	return v.expiresAt.Before(time.Now())
}
