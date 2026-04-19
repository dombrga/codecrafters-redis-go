package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/redis"
)

type RedisCommandHandlers struct {
	Store *redis.RedisStore
}

func NewRedisCommandHandlers(r *redis.RedisStore) *RedisCommandHandlers {
	return &RedisCommandHandlers{
		Store: r,
	}
}

func (h *RedisCommandHandlers) HandleEchoCmd(arg string) string {
	// return in the format of Redis bulk string, that is, $<length>\r\n<data>\r\n.
	return fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
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

func (h *RedisCommandHandlers) HandleGet(args []string) redis.RedisValueType {
	v, err := h.Store.Get(args[0])
	if err != nil {
		return "$-1\r\n"
	}

	return v
}
