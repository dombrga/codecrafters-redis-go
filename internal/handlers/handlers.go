package handlers

import "fmt"

func HandleEchoCmd(arg string) string {
	// return in the format of Redis bulk string, that is, $<length>\r\n<data>\r\n.
	return fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
}
