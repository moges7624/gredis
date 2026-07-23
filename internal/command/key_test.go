package command

import "testing"

func TestHandleExists(t *testing.T) {
	d := newDispatcher(t)

	t.Run("with no arguments", func(t *testing.T) {
		want := "-ERR wrong number of arguments for 'exists' command\r\n"
		got := handle(t, d, "EXISTS")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("with non-existing key", func(t *testing.T) {
		want := int64(0)
		got := parseIntegerReply(t, handle(t, d, "EXISTS", "name"))

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("with existing key", func(t *testing.T) {
		if got := handle(t, d, "SET", "name", "john"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for exists: got %q, want %q", got, "+OK\r\n")
		}

		want := int64(1)
		got := parseIntegerReply(t, handle(t, d, "EXISTS", "name"))

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		t.Run("counting same key twice", func(t *testing.T) {
			want := int64(2)
			got := parseIntegerReply(t, handle(t, d, "EXISTS", "name", "name"))

			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	})

	t.Run("with deleted key", func(t *testing.T) {
		if got := handle(t, d, "SET", "foo", "bar"); got != "+OK\r\n" {
			t.Fatalf("Error setting key for del: got %q, want %q", got, "+OK\r\n")
		}

		if got := handle(t, d, "del", "foo"); got != ":1\r\n" {
			t.Errorf("Erro deleting key for exists: %q, want %q", got, ":1\r\n")
		}

		want := int64(0)
		got := parseIntegerReply(t, handle(t, d, "EXISTS", "foo"))

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
