package log

import "testing"

func TestPaniclnSingleNewline(t *testing.T) {
	server, messages := newTestServer(t)
	defer server.Close()
	previous := Default()
	std.Store(New(server.URL))
	defer std.Store(previous)
	func() {
		defer func() {
			if got := recover(); got != "panic 100% message\n" {
				t.Errorf("panic value = %q, want one trailing newline", got)
			}
		}()
		Panicln("panic 100%", "message")
	}()
	got := messages()
	if len(got) != 1 || got[0] != "ERRO: panic 100% message\n" {
		t.Fatalf("Panicln messages = %q, want one message with one trailing newline", got)
	}
}
