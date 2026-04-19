package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/handlers"
	"github.com/codecrafters-io/redis-starter-go/internal/helpers"
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
	// args = ["SET", "foo", "bar", "PX", "5000"]
	args, _ := helpers.ParseRESP(_input)
	fmt.Println("args", args)
	// The input is a Redis RESP Array of bulk string,
	// that is, *<number-of-elements>\r\n<element-1>...<element-n>.
	// input := string(_input)
	// fmt.Printf("input here %q \n", input)

	// respInput := strings.Split(input, "\r\n")
	// fmt.Printf("respInput %d %q\n", len(respInput), respInput)

	cmd := strings.ToUpper(args[0])
	fmt.Printf("The command is %s\n", cmd)

	switch cmd {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		// ECHO argument is in 4th index.
		return handlers.HandleEchoCmd(args[1])
		// case "SET":
		// 	return handleSetCommand(args)
		// case "GET":
		// 	// ECHO argument is in 4th index.
		// 	return handleGetCommand(args)
	}

	return "not a valid command"
}

func handleEchoCmd(arg string) string {
	// return in the format of Redis bulk string, that is, $<length>\r\n<data>\r\n.
	return fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
}

func getCommand(input []string) string {
	// The command used is in 2nd index.
	if len(input) > 2 {
		return input[2]
	}
	return ""
}

func handleGetCommand(respInput []string) string {
	val, ok := redisStore[respInput[4]]
	// val, ok := redisStore[RedisKey(respInput[4])]
	if !ok {
		return "$-1\r\n"
	}

	// v, ok := val.Value.(string)
	// if !ok {
	// 	return "$-1\r\n"
	// }
	return fmt.Sprintf("$%d\r\n%s\r\n", len(val.(string)), val)
}

func handleSetCommand(respInput []string) string {
	var currKey string

	// TODO: use new redisStore value types and fix.
	// key := RedisKey(respInput[4])
	// val := RedisValue{
	// 	Value: respInput[6],
	// }

	// isMSExpiry := respInput[8] == "PX"
	// if isMSExpiry {
	// 	s, err := strconv.ParseFloat(respInput[8], 64)
	// 	// TODO
	// 	if err != nil {
	// 	}

	// 	val.Expiry = time.Now().Add(time.Duration(float64(time.Millisecond) * s))
	// }
	// redisStore[key] = val

	// Start loop after the command index.
	for _, r := range respInput[3 : len(respInput)-1] {
		/**
			redis bulk string:
			1st - length of next elem
			2nd - SET
			3rd - length of next elem
			4th - key
			5th - length of next elem
			6th - value
		**/
		// combo := []any{}
		// fmt.Println("kv", currKey, redisStore, r)
		if string(r[0]) != "$" {
			if currKey == "" {
				currKey = r
				redisStore[r] = nil
			}

			v, _ := redisStore[currKey]
			// fmt.Println("setting", v, r, currKey)
			if v == nil && r != currKey && string(r[0]) != "$" {
				redisStore[currKey] = r
				currKey = ""
			}
		}

		// val = nil
		// fmt.Println("-------------")
	}

	fmt.Println("redisStore", redisStore)

	return "+OK\r\n"
}
