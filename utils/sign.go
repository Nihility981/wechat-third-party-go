package utils

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

func SHA1(token, timestamp, nonce, encrypt string) string {
	array := []string{token, timestamp, nonce, ""}
	sort.Strings(array)
	str := strings.Join(array, "")

	hash := sha1.New()
	hash.Write([]byte(str))
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum
}
