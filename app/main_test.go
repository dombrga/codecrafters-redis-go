package main

import (
	"testing"
)

func TestPingCommand(t *testing.T) {
	// var redisStore1 redis.RedisStore
	// redisHandler := handlers.NewRedisCommandHandlers(&redisStore1)

	arg := "*1\r\n$4\r\nPING\r\n"

	actual := handleIncomingCommand([]byte(arg))
	expected := "+PONG\r\n"

	// fmt.Println(actual, expected)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestEchoCommand(t *testing.T) {
	arg := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"

	actual := handleIncomingCommand([]byte(arg))
	expected := "$3\r\nhey\r\n"

	// fmt.Println(actual, expected)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestSetCommandResponse(t *testing.T) {
	// arg := "*2\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	arg := "*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nPX\r\n$3\r\n100"

	expected := "+OK\r\n"
	actual := handleIncomingCommand([]byte(arg))

	// fmt.Println(actual, expected)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestSetKeyValueCommand(t *testing.T) {
	arg := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	handleIncomingCommand([]byte(arg))

	expected := "bar"
	actual, err := redisStore1.Get("foo")
	if err != nil {
		t.Errorf("expected %s, got %s", "$-1\\r\\n", actual)
	}

	if string(actual) != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

// func TestSetWithPXValueCommand(t *testing.T) {
// 	arg := "*2\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nPX\r\n$3\r\n100"

// 	handleIncomingCommand([]byte(arg))
// 	actual := redisStore["foo"]
// 	expected := "bar"
// 	// expected := RedisValue{
// 	// 	Value: "bar",
// 	// 	// Expiry: time.Now(),
// 	// }
// 	// fmt.Println("redisStore", redisStore)
// 	// fmt.Println(actual, expected)

// 	if actual != expected {
// 		t.Errorf("expected %s, got %s", expected, actual)
// 	}
// }
