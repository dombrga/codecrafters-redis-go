package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse the raw bytes from conn.Read into a slice of strings
func ParseRESP(data []byte) ([]string, error) {
	parts := strings.Split(string(data), "\r\n")
	// parts[0] is like "*3", telling you how many args follow
	count, _ := strconv.Atoi(parts[0][1:])
	fmt.Println("parts", count, parts)

	args := make([]string, 0, count)
	i := 1
	for len(args) < count {
		// fmt.Printf("i: %d, parts[i+1]: %s\n", i, parts[i+1])
		// parts[i] is like "$3" (length prefix)
		// parts[i+1] is the actual value
		args = append(args, parts[i+1])
		i += 2
	}
	return args, nil
}
