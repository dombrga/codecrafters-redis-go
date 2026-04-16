package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var rcache = map[string]any{}

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
	// The input is a Redis RESP Array of bulk string,
	// that is, *<number-of-elements>\r\n<element-1>...<element-n>.
	input := string(_input)
	fmt.Printf("input here %q \n", input)

	respInput := strings.Split(input, "\r\n")
	// fmt.Printf("zxc %d %q\n", len(respInput), respInput)

	cmd := getCommand(respInput)
	fmt.Printf("The command is %s\n", cmd)

	switch cmd {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		// ECHO argument is in 4th index.
		return handleEchoCmd(respInput[4])
	case "SET":
		return handleSetCommand(respInput)
	case "GET":
		// ECHO argument is in 4th index.
		return handleGetCommand(respInput)
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

func handleSetCommand(respInput []string) string {
	var currKey string
	// var val any

	// Start loop after the command index.
	for _, r := range respInput[3 : len(respInput)-1] {
		// fmt.Println("kv", currKey, rcache, r)
		if string(r[0]) != "$" {
			if currKey == "" {
				currKey = r
				rcache[r] = nil
			}

			v, _ := rcache[currKey]
			// fmt.Println("setting", v, r, currKey)
			if v == nil && r != currKey && string(r[0]) != "$" {
				rcache[currKey] = r
				currKey = ""
			}
		}

		// val = nil
		// fmt.Println("-------------")
	}

	fmt.Println("rcache", rcache)

	return "+OK\r\n"
}

func handleGetCommand(respInput []string) string {
	val, ok := rcache[respInput[4]]
	if !ok {
		return "$-1\r\n"
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", len(val.(string)), val)
}
