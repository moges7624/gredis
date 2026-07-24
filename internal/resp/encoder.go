package resp

import (
	"strconv"
	"strings"
)

var CRLF = "\r\n"

func EncodeSimpleString(s string) []byte {
	return []byte("+" + s + CRLF)
}

func EncodeNullString() []byte {
	return []byte("$-1" + CRLF)
}

func EncodeBulkString(s string) []byte {
	return []byte("$" + strconv.Itoa(len(s)) + CRLF + s + CRLF)
}

func EncodeError(err string) []byte {
	if strings.HasPrefix(err, "WRONGTYPE") {
		return []byte("-" + err + CRLF)
	}

	return []byte("-ERR " + err + CRLF)
}

func EncodeInteger(n int64) []byte {
	return []byte(":" + strconv.FormatInt(n, 10) + CRLF)
}
