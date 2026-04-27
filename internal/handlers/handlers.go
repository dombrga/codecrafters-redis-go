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

func NewRedisHandler(store *store.RedisStore) *RedisCommandHandlers {
	return &RedisCommandHandlers{
		Store: store,
	}
}

func (h *RedisCommandHandlers) HandleEchoCmd(arg string) string {
	// return in the format of Redis bulk string, that is, $<length>\r\n<data>\r\n.
	return fmt.Sprintf(REDIS_BULK_STRING, len(arg), arg)
}

func (h *RedisCommandHandlers) HandleSet(args []string) string {
	key, value := args[0], args[1]
	var ttl time.Duration
	fmt.Println("SET args", args)

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

/*
The RPUSH command is used to append elements to a list. If the list doesn't exist, it is created first.
The return value is the number of elements in the list after appending. This value is encoded as a RESP integer.
*/
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
