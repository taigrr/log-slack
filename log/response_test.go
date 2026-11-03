package log

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type trackedResponseBody struct {
	reader    io.Reader
	bytesRead int
	closed    bool
}

func (body *trackedResponseBody) Read(buffer []byte) (int, error) {
	count, err := body.reader.Read(buffer)
	body.bytesRead += count
	return count, err
}

func (body *trackedResponseBody) Close() error {
	body.closed = true
	return nil
}

type failingResponseReader struct{ err error }

func (reader failingResponseReader) Read([]byte) (int, error) { return 0, reader.err }

func TestWebhookResponseHandling(t *testing.T) {
	readErr := errors.New("response read failed")
	tests := []struct {
		name      string
		status    int
		reader    io.Reader
		wantRead  int
		wantError string
		wantCause error
	}{
		{name: "success drains body", status: http.StatusOK, reader: strings.NewReader("ok"), wantRead: 2},
		{name: "success drains beyond error limit", status: http.StatusOK, reader: strings.NewReader(strings.Repeat("x", slackResponseBodyLimit+100)), wantRead: slackResponseBodyLimit + 100},
		{name: "last successful status", status: http.StatusMultipleChoices - 1, reader: strings.NewReader("accepted"), wantRead: 8},
		{name: "redirect rejected", status: http.StatusMultipleChoices, reader: strings.NewReader("redirect"), wantRead: 8, wantError: "slack webhook returned status 300 Multiple Choices: redirect"},
		{name: "error text trimmed", status: http.StatusBadRequest, reader: strings.NewReader(" invalid_payload\n"), wantRead: 17, wantError: "slack webhook returned status 400 Bad Request: invalid_payload"},
		{name: "empty error body", status: http.StatusBadRequest, reader: strings.NewReader(""), wantError: "slack webhook returned status 400 Bad Request"},
		{name: "whitespace error body", status: http.StatusBadRequest, reader: strings.NewReader(" \n"), wantRead: 2, wantError: "slack webhook returned status 400 Bad Request"},
		{name: "error body bounded", status: http.StatusInternalServerError, reader: strings.NewReader(strings.Repeat("x", slackResponseBodyLimit+100)), wantRead: slackResponseBodyLimit, wantError: "slack webhook returned status 500 Internal Server Error: " + strings.Repeat("x", slackResponseBodyLimit)},
		{name: "read error preserves cause", status: http.StatusBadGateway, reader: failingResponseReader{err: readErr}, wantError: "slack webhook returned status 502 Bad Gateway and the response body could not be read: response read failed", wantCause: readErr},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := &trackedResponseBody{reader: test.reader}
			client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: test.status, Status: fmt.Sprintf("%d %s", test.status, http.StatusText(test.status)), Body: body, Header: make(http.Header)}, nil
			})}
			logger := New("https://example.com/webhook")
			writer := logger.Writer
			writer.Client = client
			logger.WithWriter(writer).Info("test message")
			err := logger.Err()
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("unexpected send error: %v", err)
				}
			} else if err == nil || err.Error() != test.wantError {
				t.Fatalf("send error = %v, want %q", err, test.wantError)
			}
			if test.wantCause != nil && !errors.Is(err, test.wantCause) {
				t.Errorf("send error does not wrap %v: %v", test.wantCause, err)
			}
			if !body.closed {
				t.Error("response body was not closed")
			}
			if body.bytesRead != test.wantRead {
				t.Errorf("read %d response bytes, want %d", body.bytesRead, test.wantRead)
			}
		})
	}
}
