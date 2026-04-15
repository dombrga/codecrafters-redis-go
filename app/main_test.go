package main

import (
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
