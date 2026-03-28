package log

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// slackMessage represents the JSON payload sent to Slack webhooks.
type slackMessage struct {
	Text string `json:"text"`
}

// newTestServer creates an httptest server that records received messages.
// Returns the server and a function to retrieve received messages.
func newTestServer(t *testing.T) (*httptest.Server, func() []string) {
	t.Helper()
	var (
		mu       sync.Mutex
		messages []string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		var msg slackMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("unmarshaling message: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mu.Lock()
		messages = append(messages, msg.Text)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	get := func() []string {
		mu.Lock()
		defer mu.Unlock()
		cp := make([]string, len(messages))
		copy(cp, messages)
		return cp
	}
	return srv, get
}

func TestNew(t *testing.T) {
	logger := New("https://example.com/webhook")
	if logger == nil {
		t.Fatal("New returned nil")
	}
	if logger.Writer.Log != "https://example.com/webhook" {
		t.Errorf("expected Log webhook to be set, got %q", logger.Writer.Log)
	}
	if logger.Writer.Error != "https://example.com/webhook" {
		t.Errorf("expected Error webhook to be set, got %q", logger.Writer.Error)
	}
	if logger.Writer.Level != LevelTrace {
		t.Errorf("expected default level LevelTrace, got %d", logger.Writer.Level)
	}
}

func TestDefault(t *testing.T) {
	d := Default()
	if d == nil {
		t.Fatal("Default returned nil")
	}
	if d != std {
		t.Error("Default should return the package-level std logger")
	}
}

func TestWithLevel(t *testing.T) {
	logger := New("https://example.com/webhook")
	updated := logger.WithLevel(LevelWarning)
	if updated.Writer.Level != LevelWarning {
		t.Errorf("expected LevelWarning, got %d", updated.Writer.Level)
	}
}

func TestWithWriter(t *testing.T) {
	logger := New("https://example.com/webhook")
	w := LogWriter{
		Log:   "https://other.com/log",
		Error: "https://other.com/error",
		Level: LevelError,
	}
	updated := logger.WithWriter(w)
	if updated.Writer.Log != "https://other.com/log" {
		t.Errorf("expected updated Log webhook, got %q", updated.Writer.Log)
	}
	if updated.Writer.Level != LevelError {
		t.Errorf("expected LevelError, got %d", updated.Writer.Level)
	}
}

func TestSetPrefixAndPrefix(t *testing.T) {
	logger := New("https://example.com/webhook")
	logger.SetPrefix("[test] ")
	if logger.Writer.prefix != "[test] " {
		t.Errorf("expected prefix %q, got %q", "[test] ", logger.Writer.prefix)
	}
}

func TestFlagsAndSetFlags(t *testing.T) {
	original := Flags()
	SetFlags(42)
	if Flags() != 42 {
		t.Errorf("expected flags 42, got %d", Flags())
	}
	SetFlags(original)
}

func TestLoggerLog(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Log("hello world")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "INFO:") || !strings.Contains(msgs[0], "hello world") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerLogf(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Logf("count: %d", 42)
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "count: 42") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerLogln(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Logln("hello", "world")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "hello world") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerError(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Error("something broke")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "ERRO:") {
		t.Errorf("expected ERRO prefix: %q", msgs[0])
	}
}

func TestLoggerErrorf(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Errorf("error: %s", "test")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "error: test") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerErrorln(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Errorln("error", "line")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "ERRO:") {
		t.Errorf("expected ERRO prefix: %q", msgs[0])
	}
}

func TestLoggerWarning(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Warning("caution")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "WARN:") {
		t.Errorf("expected WARN prefix: %q", msgs[0])
	}
}

func TestLoggerWarningf(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Warningf("warn: %d", 1)
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "warn: 1") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerWarningln(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Warningln("warn", "line")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "WARN:") {
		t.Errorf("expected WARN prefix: %q", msgs[0])
	}
}

func TestLoggerInfo(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Info("info msg")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "INFO:") || !strings.Contains(msgs[0], "info msg") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
	if logger.Err() != nil {
		t.Errorf("unexpected error: %v", logger.Err())
	}
}

func TestLoggerInfof(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Infof("info: %s", "formatted")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "info: formatted") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerInfoln(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Infoln("info", "line")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "INFO:") {
		t.Errorf("expected INFO prefix: %q", msgs[0])
	}
}

func TestLoggerDebug(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Debug("debug msg")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "DEBG:") {
		t.Errorf("expected DEBG prefix: %q", msgs[0])
	}
}

func TestLoggerDebugf(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Debugf("debug: %v", true)
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "debug: true") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerDebugln(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Debugln("debug", "line")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "DEBG:") {
		t.Errorf("expected DEBG prefix: %q", msgs[0])
	}
}

func TestLoggerTrace(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Trace("trace msg")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "TRCE:") {
		t.Errorf("expected TRCE prefix: %q", msgs[0])
	}
}

func TestLoggerTracef(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Tracef("trace: %s", "detail")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "trace: detail") {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestLoggerTraceln(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.Traceln("trace", "line")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.Contains(msgs[0], "TRCE:") {
		t.Errorf("expected TRCE prefix: %q", msgs[0])
	}
}

func TestLogLevelFiltering(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	updated := logger.WithLevel(LevelWarning)

	// These should be filtered out (below warning level)
	updated.Debug("should not appear")
	updated.Trace("should not appear")
	updated.Info("should not appear")

	// These should pass through
	updated.Warning("warning msg")
	updated.Error("error msg")

	msgs := getMessages()
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages (warning + error), got %d: %v", len(msgs), msgs)
	}
}

func TestLogLevelFilteringErrorOnly(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	updated := logger.WithLevel(LevelError)

	updated.Warning("should not appear")
	updated.Info("should not appear")
	updated.Debug("should not appear")
	updated.Trace("should not appear")
	updated.Error("error only")

	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d: %v", len(msgs), msgs)
	}
	if !strings.Contains(msgs[0], "ERRO:") {
		t.Errorf("expected error message, got: %q", msgs[0])
	}
}

func TestWriteInterface(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	n, err := logger.Writer.Write([]byte("io.Writer test"))
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != len("io.Writer test") {
		t.Errorf("expected %d bytes written, got %d", len("io.Writer test"), n)
	}
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0] != "io.Writer test" {
		t.Errorf("unexpected message: %q", msgs[0])
	}
}

func TestPrefix(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()
	logger := New(srv.URL)
	logger.SetPrefix("[APP] ")
	logger.Log("prefixed msg")
	msgs := getMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if !strings.HasPrefix(msgs[0], "[APP] ") {
		t.Errorf("expected prefix [APP], got: %q", msgs[0])
	}
}

func TestPrintFunctions(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()

	// Override the std logger temporarily
	oldStd := std
	std = New(srv.URL)
	defer func() { std = oldStd }()

	Print("print test")
	Printf("printf %d", 99)
	Println("println test")

	msgs := getMessages()
	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d: %v", len(msgs), msgs)
	}
}

func TestPackageLevelFunctions(t *testing.T) {
	srv, getMessages := newTestServer(t)
	defer srv.Close()

	oldStd := std
	std = New(srv.URL)
	defer func() { std = oldStd }()

	Log("log msg")
	Logf("logf %s", "msg")
	Logln("logln msg")
	Error("error msg")
	Errorf("errorf %s", "msg")
	Errorln("errorln msg")
	Warning("warning msg")
	Warningf("warningf %s", "msg")
	Warningln("warningln msg")
	Info("info msg")
	Infof("infof %s", "msg")
	Infoln("infoln msg")
	Debug("debug msg")
	Debugf("debugf %s", "msg")
	Debugln("debugln msg")
	Trace("trace msg")
	Tracef("tracef %s", "msg")
	Traceln("traceln msg")

	msgs := getMessages()
	if len(msgs) != 18 {
		t.Fatalf("expected 18 messages, got %d", len(msgs))
	}
}

func TestPackageLevelSetPrefix(t *testing.T) {
	oldStd := std
	std = New("")
	defer func() { std = oldStd }()

	SetPrefix("[PKG] ")
	if Prefix() != "[PKG] " {
		t.Errorf("expected prefix %q, got %q", "[PKG] ", Prefix())
	}
}

func TestPackageLevelWithLevel(t *testing.T) {
	oldStd := std
	std = New("")
	defer func() { std = oldStd }()

	updated := WithLevel(LevelDebug)
	if updated.Writer.Level != LevelDebug {
		t.Errorf("expected LevelDebug, got %d", updated.Writer.Level)
	}
}

func TestPackageLevelWithWriter(t *testing.T) {
	oldStd := std
	std = New("")
	defer func() { std = oldStd }()

	w := LogWriter{Log: "https://example.com", Level: LevelInfo}
	updated := WithWriter(w)
	if updated.Writer.Log != "https://example.com" {
		t.Errorf("expected updated writer, got %q", updated.Writer.Log)
	}
}

func TestPanicFunctions(t *testing.T) {
	srv, _ := newTestServer(t)
	defer srv.Close()

	oldStd := std
	std = New(srv.URL)
	defer func() { std = oldStd }()

	t.Run("Panic", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}
			if s, ok := r.(string); !ok || s != "panic test" {
				t.Errorf("unexpected panic value: %v", r)
			}
		}()
		Panic("panic test")
	})

	t.Run("Panicf", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}
		}()
		Panicf("panic %d", 42)
	})

	t.Run("Panicln", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}
		}()
		Panicln("panic", "line")
	})
}

func TestLogLevelConstants(t *testing.T) {
	if LevelError != 0 {
		t.Errorf("LevelError should be 0, got %d", LevelError)
	}
	if LevelWarning != 1 {
		t.Errorf("LevelWarning should be 1, got %d", LevelWarning)
	}
	if LevelInfo != 2 {
		t.Errorf("LevelInfo should be 2, got %d", LevelInfo)
	}
	if LevelDebug != 3 {
		t.Errorf("LevelDebug should be 3, got %d", LevelDebug)
	}
	if LevelTrace != 4 {
		t.Errorf("LevelTrace should be 4, got %d", LevelTrace)
	}
}

func TestErrTrackingOnBadWebhook(t *testing.T) {
	// Use an invalid URL that will fail to connect
	logger := New("http://127.0.0.1:1")
	logger.Info("should fail")
	if logger.Err() == nil {
		t.Error("expected error from bad webhook URL")
	}
}
