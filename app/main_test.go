package main

import (
	"fmt"
	"testing"
)

func TestPingCommand(t *testing.T) {
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
	arg := "*4\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\n"

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
	fmt.Println("entry:", actual)

	if string(actual) != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestSetWithEXValueCommand(t *testing.T) {
	arg := "*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nEX\r\n$3\r\n100"
	handleIncomingCommand([]byte(arg))

	expected := "bar"
	actual, err := redisStore1.Get("foo")
	if err != nil {
		t.Errorf("expected %s, got %s", "$-1\\r\\n", actual)
	}
	// fmt.Println("entry:", actual)

	if string(actual) != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
