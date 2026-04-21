package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

const REDIS_BULK_STRING = "$%d\r\n%s\r\n"
const REDIS_RESP_INTEGER = ":%s\r\n"

type RedisCommandHandlers struct {
	Store *store.RedisStore
}

func NewRedisHandler(r *store.RedisStore) *RedisCommandHandlers {
	return &RedisCommandHandlers{
		Store: r,
	}
}

func (h *RedisCommandHandlers) HandleEchoCmd(arg string) string {
	// return in the format of Redis bulk string, that is, $<length>\r\n<data>\r\n.
	return fmt.Sprintf(REDIS_BULK_STRING, len(arg), arg)
}

func (h *RedisCommandHandlers) HandleSet(args []string) string {
	key, value := args[0], args[1]
	var ttl time.Duration

	// Walk remaining args as option pairs
	for i := 2; i < len(args)-1; i += 2 {
		switch strings.ToUpper(args[i]) {
		case "PX":
			ms, _ := strconv.ParseInt(args[i+1], 10, 64)
			ttl = time.Duration(ms) * time.Millisecond
		case "EX":
			sec, _ := strconv.ParseInt(args[i+1], 10, 64)
			ttl = time.Duration(sec) * time.Second
		}
	}

	h.Store.Set(key, value, ttl)

	return "+OK\r\n"
}

func (h *RedisCommandHandlers) HandleGet(args []string) string {
	val, err := h.Store.Get(args[0])
	fmt.Printf("val: %q\n", val)
	if err != nil {
		return string(val)
	}

	return fmt.Sprintf(REDIS_BULK_STRING, len(val), val)
}

func (h *RedisCommandHandlers) HandleRPush(args []string) string {
	key := args[0]
	rpushArgs := args[1:]
	// fmt.Println("RPUSH args:", args, key, rpushArgs)

	v, err := h.Store.RPush(key, rpushArgs)
	if err != nil {
		return string(v)
	}

	return fmt.Sprintf(REDIS_RESP_INTEGER, v)
}

// func (h *RedisCommandHandlers) HandleSetCommand(respInput []string) string {
// 	var currKey string

// 	// TODO: use new redisStore value types and fix.
// 	// key := RedisKey(respInput[4])
// 	// val := RedisValue{
// 	// 	Value: respInput[6],
// 	// }

// 	// isMSExpiry := respInput[8] == "PX"
// 	// if isMSExpiry {
// 	// 	s, err := strconv.ParseFloat(respInput[8], 64)
// 	// 	// TODO
// 	// 	if err != nil {
// 	// 	}

// 	// 	val.Expiry = time.Now().Add(time.Duration(float64(time.Millisecond) * s))
// 	// }
// 	// redisStore[key] = val

// 	// Start loop after the command index.
// 	for _, r := range respInput[3 : len(respInput)-1] {
// 		/**
// 			redis bulk string:
// 			1st - length of next elem
// 			2nd - SET
// 			3rd - length of next elem
// 			4th - key
// 			5th - length of next elem
// 			6th - value
// 		**/
// 		// combo := []any{}
// 		// fmt.Println("kv", currKey, redisStore, r)
// 		if string(r[0]) != "$" {
// 			if currKey == "" {
// 				currKey = r
// 				redisStore[r] = nil
// 			}

// 			v, _ := redisStore[currKey]
// 			// fmt.Println("setting", v, r, currKey)
// 			if v == nil && r != currKey && string(r[0]) != "$" {
// 				redisStore[currKey] = r
// 				currKey = ""
// 			}
// 		}

// 		// val = nil
// 		// fmt.Println("-------------")
// 	}

// 	fmt.Println("redisStore", redisStore)

// 	return "+OK\r\n"
// }
