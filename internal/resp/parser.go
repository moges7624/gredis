// Package resp implements RESP parsing
package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

var ErrInvalidProtocol = errors.New("protocol error")

type Value struct {
	Type  ValueType
	Str   string
	Int   int64
	Array []Value
}

type ValueType int

const (
	SimpleString ValueType = iota
	Error
	Integer
	BulkString
	Array
	Null
)

type Parser struct {
	reader *bufio.Reader
	logger *slog.Logger
}

func NewParser(r *bufio.Reader, l *slog.Logger) *Parser {
	return &Parser{
		reader: r,
		logger: l,
	}
}

func (p *Parser) Parse() (Value, error) {
	typeByte, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte {
	case '+':
		return p.parseSimpleString()
	case '-':
		return p.parseError()
	case ':':
		return p.parseInteger()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArray()
	default:
		return Value{}, ErrInvalidProtocol
	}
}

func (p *Parser) parseSimpleString() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	str := strings.TrimSpace(line)
	return Value{Type: SimpleString, Str: str}, err
}

func (p *Parser) parseError() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	str := strings.TrimSpace(line)
	return Value{Type: Error, Str: str}, nil
}

func (p *Parser) parseInteger() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
	if err != nil {
		return Value{}, err
	}
	return Value{Type: Integer, Int: i}, nil
}

func (p *Parser) parseArray() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return Value{}, err
	}

	array := make([]Value, count)
	for i := range count {
		elem, err := p.parseBulkString()
		if err != nil {
			return Value{}, err
		}

		array[i] = elem
	}

	return Value{Type: Array, Array: array}, nil
}

func (p *Parser) parseBulkString() (Value, error) {
	_, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return Value{}, err
	}

	data := make([]byte, length+2)
	_, err = p.reader.Read(data)
	if err != nil {
		return Value{}, err
	}

	if !bytes.HasSuffix(data, []byte("\r\n")) {
		return Value{}, io.ErrUnexpectedEOF
	}

	return Value{Type: BulkString, Str: string(data[:length])}, nil
}
