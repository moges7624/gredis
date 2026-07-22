package resp

import (
	"strconv"
)

var CRLF = "\r\n"

func EncodeSimpleString(s string) []byte {
	return []byte("+" + s + CRLF)
}

func EncodeBulkString(s string) []byte {
	if s == "-1" {
		return []byte("$-1" + CRLF)
	}

	return []byte("$" + strconv.Itoa(len(s)) + CRLF + s + CRLF)
}

func EncodeError(err string) []byte {
	return []byte("-ERR " + err + CRLF)
}

func EncodeInteger(n int) []byte {
	return []byte(":" + strconv.Itoa(n) + CRLF)
}
