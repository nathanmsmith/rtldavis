package processor

import (
	"fmt"
	"strings"
)

func bytesToSpacedHex(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	result := make([]string, len(data))
	for i, b := range data {
		result[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(result, " ")
}
