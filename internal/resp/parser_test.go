package resp

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr error
	}{
		{
			name:  "Valid Ping Inline/Simple String",
			input: "+PING\r\n",
			want: Value{
				Type: SimpleString,
				Str:  "PING",
			},
			wantErr: nil,
		},
		{
			name:  "Valid integer",
			input: ":33\r\n",
			want: Value{
				Type: Integer,
				Int:  33,
			},
			wantErr: nil,
		},
		{
			name:  "Valid MGET Array of Bulk Strings",
			input: "*3\r\n$4\r\nMGET\r\n$4\r\nkey1\r\n$4\r\nkey2\r\n",
			want: Value{
				Type: Array,
				Array: []Value{
					{Type: BulkString, Str: "MGET"},
					{Type: BulkString, Str: "key1"},
					{Type: BulkString, Str: "key2"},
				},
			},
			wantErr: nil,
		},
		{
			name:    "Invalid Protocol Byte",
			input:   "XInvalid\r\n",
			want:    Value{},
			wantErr: ErrInvalidProtocol,
		},
		{
			name:    "Unexpected EOF",
			input:   "*2\r\n$4\r\nGET\r\n",
			want:    Value{},
			wantErr: io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			parser := NewParser(reader, logger)

			res, err := parser.Parse()

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("Parse() got = %+v, want = %+v", res, tt.want)
			}
		})
	}
}
