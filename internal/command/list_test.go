package command

import "testing"

func TestHandleLLen(t *testing.T) {
	d := newDispatcher(t)

	t.Run("with no key argument", func(t *testing.T) {
		want := wrongArgumentNumberErrString("llen")
		got := handle(t, d, "llen")

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("With non-existing key", func(t *testing.T) {
		resp := parseIntegerReply(t, handle(t, d, "LLEN", "name"))

		if resp != 0 {
			t.Errorf("got %d, want %d", resp, 0)
		}
	})

	t.Run("With existing key of non-list value", func(t *testing.T) {
		setValueWithDeleteCleanUp(t, d, "foo", "bar")

		want := "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
		got := handle(t, d, "LLEN", "foo")
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("With existing key of list value", func(t *testing.T) {
		res := parseIntegerReply(t, handle(t, d, "RPUSH", "name", "john"))
		if res != 1 {
			t.Fatalf("failed to rpush for llen command")
		}

		got := parseIntegerReply(t, handle(t, d, "LLEN", "name"))
		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}

		deleteKey(t, d, "name")
	})
}

func TestHandleRPush(t *testing.T) {
	d := newDispatcher(t)

	t.Run("With no values", func(t *testing.T) {
		want := wrongArgumentNumberErrString("rpush")
		got := handle(t, d, "rpush", "name")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("With non-existing key", func(t *testing.T) {
		got := parseIntegerReply(t, handle(t, d, "RPUSH", "name", "john"))

		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}

		deleteKey(t, d, "name")
	})

	t.Run("With existing key of non-list value", func(t *testing.T) {
		setValueWithDeleteCleanUp(t, d, "name", "john")

		want := "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
		got := handle(t, d, "RPUSH", "name", "rob")
		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("With existing key of list value", func(t *testing.T) {
		res := parseIntegerReply(t, handle(t, d, "RPUSH", "name", "john"))
		if res != 1 {
			t.Fatalf("failed to rpush for llen command")
		}

		got := parseIntegerReply(t, handle(t, d, "RPUSH", "name", "rob"))
		if got != 2 {
			t.Errorf("got %d, want %d", got, 2)
		}

		deleteKey(t, d, "name")
	})
}
