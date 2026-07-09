// Package resp implements RESP parsing
package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Value struct {
	Type  ValueType
	Str   string
	Int   int64
	Array []Value
}

// ValueType defines RESP data types
type ValueType int

const (
	SimpleString ValueType = iota
	Error
	Integer
	BulkString
	Array
	Null
)

func Parse(r *bufio.Reader, logger *slog.Logger) (Value, error) {
	typeByte, err := r.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte {
	case '+':
		return parseSimpleString(r, logger)
	case '-':
		return parseError(r)
	case ':':
		return parseInteger(r)
	case '$':
		return parseBulkString(r, logger)
	case '*':
		return parseArray(r, logger)
	default:
		return Value{}, fmt.Errorf("invalid type")
	}
}

func parseSimpleString(r *bufio.Reader, logger *slog.Logger) (Value, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	str := strings.TrimSpace(line)
	return Value{Type: SimpleString, Str: str}, err
}

func parseError(r *bufio.Reader) (Value, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	str := strings.TrimSpace(line)
	return Value{Type: Error, Str: str}, nil
}

func parseInteger(r *bufio.Reader) (Value, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
	if err != nil {
		return Value{}, err
	}
	return Value{Type: Integer, Int: i}, nil
}

func parseArray(r *bufio.Reader, logger *slog.Logger) (Value, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return Value{}, err
	}

	array := make([]Value, count)
	for i := range count {
		elem, err := parseBulkString(r, logger)
		if err != nil {
			return Value{}, nil
		}

		array[i] = elem
	}

	return Value{Type: Array, Array: array}, nil
}

func parseBulkString(r *bufio.Reader, logger *slog.Logger) (Value, error) {
	_, err := r.ReadByte()
	if err != nil {
		return Value{}, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return Value{}, err
	}

	data := make([]byte, length+2)
	_, err = r.Read(data)
	if err != nil {
		return Value{}, err
	}

	if !bytes.HasSuffix(data, []byte("\r\n")) {
		return Value{}, fmt.Errorf("invalid protocol")
	}

	return Value{Type: BulkString, Str: string(data[:length])}, nil
}
