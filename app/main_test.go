package main

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func newTestStore() *store.RedisStore {
	return &store.RedisStore{Items: map[store.RedisEntryKey]store.RedisEntryValue{}}
}

func TestPingCommand(t *testing.T) {
	store := newTestStore()
	arg := "*1\r\n$4\r\nPING\r\n"

	actual := handleIncomingCommand([]byte(arg), store)
	expected := "+PONG\r\n"

	// fmt.Println(actual, expected)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestEchoCommand(t *testing.T) {
	store := newTestStore()
	arg := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"

	actual := handleIncomingCommand([]byte(arg), store)
	expected := "$3\r\nhey\r\n"

	// fmt.Println(actual, expected)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestSetCommandResponse(t *testing.T) {
	store := newTestStore()
	// arg := "*2\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	arg := "*4\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\n"

	expected := "+OK\r\n"
	actual := handleIncomingCommand([]byte(arg), store)

	// fmt.Println(actual, expected)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestSetKeyValueCommand(t *testing.T) {
	store := newTestStore()
	arg := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	handleIncomingCommand([]byte(arg), store)

	expected := "bar"
	actual, err := store.Get("foo")
	if err != nil {
		t.Errorf("expected %s, got %s", "$-1\\r\\n", actual)
	}
	fmt.Println("entry:", actual)

	if string(actual) != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestSetWithEXValueCommand(t *testing.T) {
	store := newTestStore()
	arg := "*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nEX\r\n$3\r\n100"
	handleIncomingCommand([]byte(arg), store)

	expected := "bar"
	actual, err := store.Get("foo")
	if err != nil {
		t.Errorf("expected %s, got %s", "$-1\\r\\n", actual)
	}
	// fmt.Println("entry:", actual)

	if string(actual) != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestRPushCommand(t *testing.T) {
	store := newTestStore()
	// RPUSH list_key "element"
	arg := "*3\r\n$5\r\nRPUSH\r\n$9\r\npineapple\r\n$5\r\ngrape\r\n"

	expected := ":1\r\n"
	actual := handleIncomingCommand([]byte(arg), store)

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestRPushToExistingList(t *testing.T) {
	store := newTestStore()
	arg1 := "*3\r\n$5\r\nRPUSH\r\n$9\r\npineapple\r\n$5\r\ngrape\r\n"
	arg2 := "*3\r\n$5\r\nRPUSH\r\n$9\r\npineapple\r\n$6\r\nbanana\r\n"
	arg3 := "*3\r\n$5\r\nRPUSH\r\n$9\r\npineapple\r\n$4\r\npear\r\n"

	expected := ":3\r\n"
	handleIncomingCommand([]byte(arg1), store)
	handleIncomingCommand([]byte(arg2), store)
	actual := handleIncomingCommand([]byte(arg3), store)

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
