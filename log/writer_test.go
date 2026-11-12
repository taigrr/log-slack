package log

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestWriteDeliveryCounts(t *testing.T) {
	transportErr := errors.New("transport unavailable")
	for _, test := range []struct {
		name         string
		status       int
		transportErr error
		wantError    bool
	}{
		{name: "accepted", status: http.StatusOK},
		{name: "rejected", status: http.StatusBadRequest, wantError: true},
		{name: "transport failure", transportErr: transportErr, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			writer := LogWriter{Log: "https://example.com/webhook", Client: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if err := request.Body.Close(); err != nil {
					t.Errorf("close request: %v", err)
				}
				if test.transportErr != nil {
					return nil, test.transportErr
				}
				return &http.Response{StatusCode: test.status, Status: http.StatusText(test.status), Header: make(http.Header), Body: io.NopCloser(strings.NewReader(""))}, nil
			})}}
			payload := []byte("delivery count")
			count, err := writer.Write(payload)
			if (err != nil) != test.wantError {
				t.Fatalf("Write error = %v, want error %v", err, test.wantError)
			}
			wantCount := len(payload)
			if test.wantError {
				wantCount = 0
			}
			if count != wantCount {
				t.Errorf("Write count = %d, want %d", count, wantCount)
			}
			if test.transportErr != nil && !errors.Is(err, test.transportErr) {
				t.Errorf("lost transport cause: %v", err)
			}
		})
	}
}
