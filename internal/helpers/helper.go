package helpers

import (
	"fmt"
	"strings"
)

// "*4\r\n$5\r\nRPUSH\r\n$9\r\npineapple\r\n$5\r\ngrape\r\n$4\r\npear\r\n"
func EncodeAsBulkString(args string) string {
	args = `RPUSH list_key "a" "b" "c" "d" "e"`
	str := strings.Split(args, " ")
	fmt.Println("strs", str)

	ret := fmt.Sprintf("*%d", len(str))
	for _, s := range str {
		_s := strings.ReplaceAll(s, "\"", "")
		ret = fmt.Sprintf("%s\r\n$%d\r\n%s", ret, len(_s), _s)
	}

	fmt.Printf("ret: %q\n", ret+"\r\n")
	return ret + "\r\n"
}
