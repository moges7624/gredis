package command

import (
	"testing"
)

func TestUnknownCommand(t *testing.T) {
	d := newDispatcher(t)

	got := handle(t, d, "write")
	want := "-ERR unknown command 'WRITE'\r\n"
	if got != want {
		t.Errorf("Unknown command: got %q, want %q", got, want)
	}
}

func TestPing(t *testing.T) {
	d := newDispatcher(t)

	t.Run("PING with no arg", func(t *testing.T) {
		if got := handle(t, d, "PING"); got != "+PONG\r\n" {
			t.Errorf("got %q, want %q", got, "+PONG\r\n")
		}
	})

	t.Run("PING with 1 arg", func(t *testing.T) {
		if got := handle(t, d, "PING", "hello"); got != "$5\r\nhello\r\n" {
			t.Errorf("got %q, want %q", got, "$5\r\nhello\r\n")
		}
	})

	t.Run("PING with more than 1 args", func(t *testing.T) {
		got := handle(t, d, "PING", "hello", "there")
		want := "-ERR wrong number of arguments for 'ping' command\r\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestEcho(t *testing.T) {
	d := newDispatcher(t)
	t.Run("ECHO with no arg", func(t *testing.T) {
		got := handle(t, d, "echo")
		want := "-ERR wrong number of arguments for 'echo' command\r\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ECHO with 1 arg", func(t *testing.T) {
		if got := handle(t, d, "echo", "hello"); got != "$5\r\nhello\r\n" {
			t.Errorf("got %q, want %q", got, "$5\r\nhello\r\n")
		}
	})

	t.Run("ECHO with more than 1 args", func(t *testing.T) {
		got := handle(t, d, "echo", "hello", "there")
		want := "-ERR wrong number of arguments for 'echo' command\r\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestSetGet(t *testing.T) {
	d := newDispatcher(t)

	t.Run("SET new key", func(t *testing.T) {
		if got := handle(t, d, "SET", "foo", "bar"); got != "+OK\r\n" {
			t.Fatalf("got %q, want %q", got, "+OK\r\n")
		}
	})

	t.Run("GET existing key", func(t *testing.T) {
		if got := handle(t, d, "GET", "foo"); got != "$3\r\nbar\r\n" {
			t.Errorf("got %q, want %q", got, "$3\r\nbar\r\n")
		}
	})

	t.Run("GET non-existing key", func(t *testing.T) {
		if got := handle(t, d, "GET", "age"); got != "$-1\r\n" {
			t.Errorf("got %q, want %q", got, "$-\r\n")
		}
	})

	t.Run("GET with more than 1 args", func(t *testing.T) {
		want := "-ERR wrong number of arguments for 'get' command\r\n"
		got := handle(t, d, "GET", "foo", "name")
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDel(t *testing.T) {
	d := newDispatcher(t)

	t.Run("DEL with no args", func(t *testing.T) {
		want := "-ERR wrong number of arguments for 'del' command\r\n"
		got := handle(t, d, "del")
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("DEL with non-existing key", func(t *testing.T) {
		if got := handle(t, d, "del", "name"); got != ":0\r\n" {
			t.Errorf("got %q, want %q", got, ":0\r\n")
		}
	})

	t.Run("DEL with 1 existing key", func(t *testing.T) {
		if got := handle(t, d, "SET", "foo", "bar"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for del: got %q, want %q", got, "+OK\r\n")
		}

		if got := handle(t, d, "del", "foo"); got != ":1\r\n" {
			t.Errorf("got %q, want %q", got, ":1\r\n")
		}
	})

	t.Run("DEL with 2 existing key", func(t *testing.T) {
		if got := handle(t, d, "SET", "foo", "bar"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for del: got %q, want %q", got, "+OK\r\n")
		}

		if got := handle(t, d, "SET", "name", "john"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for del: got %q, want %q", got, "+OK\r\n")
		}

		if got := handle(t, d, "del", "foo", "name"); got != ":2\r\n" {
			t.Errorf("got %q, want %q", got, ":2\r\n")
		}
	})
}

func TestSet_EX(t *testing.T) {
	d := newDispatcher(t)

	t.Run("EX with no expire time arg", func(t *testing.T) {
		want := "-ERR syntax error\r\n"
		got := handle(t, d, "SET", "name", "john", "EX")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("EX with invalid expire time", func(t *testing.T) {
		want := "-ERR invalid expire time in 'set' command\r\n"
		got := handle(t, d, "SET", "name", "john", "EX", "0")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("EX with valid expire time", func(t *testing.T) {
		want := "+OK\r\n"
		got := handle(t, d, "SET", "name", "john", "EX", "120")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestIncr(t *testing.T) {
	d := newDispatcher(t)

	t.Run("With more than 1 arg", func(t *testing.T) {
		got := handle(t, d, "INCR", "age", "count")
		want := "-ERR wrong number of arguments for 'incr' command\r\n"
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("With non-existing key", func(t *testing.T) {
		if got := handle(t, d, "INCR", "age"); got != ":1\r\n" {
			t.Errorf("got %q, want %q", got, ":1\r\n")
		}
	})

	t.Run("With existing key", func(t *testing.T) {
		if got := handle(t, d, "SET", "num", "4"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for incr: got %q, want %q", got, "+OK\r\n")
		}

		if got := handle(t, d, "INCR", "num"); got != ":5\r\n" {
			t.Errorf("got %q, want %q", got, ":5\r\n")
		}
	})

	t.Run("With a key holding non-integer string value", func(t *testing.T) {
		if got := handle(t, d, "SET", "foo", "bar"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for incr: got %q, want %q", got, "+OK\r\n")
		}

		got := handle(t, d, "INCR", "foo")
		want := "-ERR value is not an integer or out of range\r\n"
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDecr(t *testing.T) {
	d := newDispatcher(t)

	t.Run("DECR with more than 1 arg", func(t *testing.T) {
		got := handle(t, d, "DECR", "age", "count")
		want := "-ERR wrong number of arguments for 'decr' command\r\n"
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("DECR non-existing key", func(t *testing.T) {
		if got := handle(t, d, "DECR", "age"); got != ":-1\r\n" {
			t.Errorf("got %q, want %q", got, ":-1\r\n")
		}
	})

	t.Run("DECR existing key", func(t *testing.T) {
		if got := handle(t, d, "SET", "num", "4"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for decr: got %q, want %q", got, "+OK\r\n")
		}

		if got := handle(t, d, "decr", "num"); got != ":3\r\n" {
			t.Errorf("got %q, want %q", got, ":3\r\n")
		}
	})

	t.Run("With a key holding non-integer string value", func(t *testing.T) {
		if got := handle(t, d, "SET", "foo", "bar"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for decr: got %q, want %q", got, "+OK\r\n")
		}

		got := handle(t, d, "DECR", "foo")
		want := "-ERR value is not an integer or out of range\r\n"
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
