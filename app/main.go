package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/handlers"
	"github.com/codecrafters-io/redis-starter-go/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type RedisKey string
type RedisValue struct {
	Value  any
	Expiry time.Time
}

type RCache map[RedisKey]RedisValue

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var redisStore = map[string]any{}

var redisStore1 = store.RedisStore{
	Items: map[store.RedisEntryKey]store.RedisEntryValue{},
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		// Accept method is blocking.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		// Once connection is accepted, open a goroutine to handle the
		// connection to prevent blocking incoming new connnections.
		go func(conn net.Conn) {
			defer conn.Close()
			for {
				b := make([]byte, 1024)
				n, err := conn.Read(b)
				if err != nil {
					return
				}

				resp := handleIncomingCommand(b[:n])
				fmt.Println("resp", resp)
				conn.Write([]byte(resp))
			}
		}(conn)
	}
}

func handleIncomingCommand(_input []byte) string {
	res := getResponse(_input)
	// fmt.Println("res", res)

	return res
}

func getResponse(_input []byte) string {
	redisHandler := handlers.NewRedisHandler(&redisStore1)

	// args = ["SET", "foo", "bar", "PX", "5000"]
	args, _ := parser.ParseRESP(_input)
	fmt.Println("args", args)

	cmd := strings.ToUpper(args[0])
	// fmt.Printf("The command is %s\n", cmd)

	switch cmd {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		return redisHandler.HandleEchoCmd(args[1])
	case "SET":
		return redisHandler.HandleSet(args[1:])
	case "GET":
		return redisHandler.HandleGet(args[1:])
	}

	return "not a valid command"
}
