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
				conn.Read(b)

				resp := handleCommandResponse(b)
				fmt.Println("resp", resp)
				conn.Write([]byte(resp))
			}
		}(conn)
	}

	// for {
	// 	conn, err := l.Accept()

	// 	cmdChan := make(chan []byte)
	// 	go func() {
	// 		for {
	// 			// var cmd string
	// 			buf := make([]byte, 1024)
	// 			conn.Read(buf)
	// 			// fmt.Scanln(&cmd)
	// 			cmdChan <- buf
	// 		}
	// 		for range cmdChan {
	// 			// fmt.Println("received:", cmd)
	// 			conn.Write([]byte("+PONG\r\n"))
	// 		}
	// 	}()

	// }

	// for {
	// 	buf := make([]byte, 1024)
	// 	conn.Read(buf)
	// 	conn.Write([]byte("+PONG\r\n"))
	// }

}

func handleCommandResponse(_input []byte) string {
	res := getRESPString(_input)
	// fmt.Println("res", res)

	// switch cmd {
	// case "ECHO":
	// 	// fmt.Println("echoing")
	// 	return res
	// }

	return res
}

func getRESPString(_input []byte) string {
	// The input is a Redis RESP Array of bulk string, that is, *<number-of-elements>\r\n<element-1>...<element-n>.
	input := string(_input)

	respInput := strings.Split(input, "\r\n")
	fmt.Printf("zxc %d %q\n", len(respInput), respInput)

	var cmd string
	// The command used is the 3rd element.
	if len(respInput) > 2 {
		cmd = respInput[2]
	}
	fmt.Printf("The command is %s\n", cmd)

	switch cmd {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		arg := respInput[4]

		// return in the format of Redis bulk string, that is, $<length>\r\n<data>\r\n.
		res := fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
		fmt.Printf("res is: %q\v", res)
		return res
	}

	return ""

}

func handleEchoCmd() {

}
