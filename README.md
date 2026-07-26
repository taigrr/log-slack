# Log-Slack

[![Test](https://github.com/taigrr/log-slack/actions/workflows/test.yml/badge.svg)](https://github.com/taigrr/log-slack/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/taigrr/log-slack.svg)](https://pkg.go.dev/github.com/taigrr/log-slack)
[![Go Report Card](https://goreportcard.com/badge/github.com/taigrr/log-slack)](https://goreportcard.com/report/github.com/taigrr/log-slack)

A Go package for sending logs to Slack via webhooks with different log levels.

## Features

- Multiple log levels: Error, Warning, Info, Debug, Trace
- Configurable webhook URLs for each log level
- Formatted logging support (Printf-style)
- Default logger instance
- Thread-safe operations

## Log Levels

```go

import (
    "github.com/taigrr/log-slack/log"
)

const (
    LevelError LogLevel = iota
    LevelWarning
    LevelInfo
    LevelDebug
    LevelTrace
)
```

## Usage

### Basic Setup

```go
// Create a new logger with a webhook URL
logger := log.New("https://hooks.slack.com/services/...")

// Set log level
logger = logger.WithLevel(log.LevelInfo)

// Set prefix
logger.SetPrefix("[MyApp] ")

// Log messages
logger.Info("This is an info message")
logger.Error("This is an error message")
```

### Using Default Logger

```go
// Configure default logger
log.Default().WithLevel(log.LevelInfo)

// Log messages
log.Info("This is an info message")
log.Errorf("Error occurred: %v", err)
```

### Available Methods

For each log level (Error, Warning, Info, Debug, Trace), the following methods are available:

- `Level(message string)`
- `Levelf(format string, args ...any)`
- `Levelln(args ...any)`

Example:

```go
log.Info("Simple message")
log.Infof("Formatted message: %s", value)
log.Infoln("Message with newline")
```

### Stdlib Compatibility

The package provides all standard library log functions:

```go
log.Print("message")
log.Printf("format %s", "message")
log.Println("message")
log.Fatal("fatal message")  // Exits with status 1
log.Panic("panic message")  // Panics after logging
```

## Configuration

### LogWriter

The `LogWriter` struct allows configuration of different webhook URLs for each log level:

```go
type LogWriter struct {
    Log     string // used by Log/Print/Info-level sends
    Error   string // used by Error/Fatal/Panic sends
    Warning string // used by Warning sends
    Info    string
    Debug   string // used by Debug sends
    Trace   string // used by Trace sends
    Level   LogLevel
}
```

Each field is the webhook URL for its level. `Info`/`Log` messages post to the
`Log` field. When constructed via `New`, every field is set to the same
webhook URL.

### Custom Configuration

```go
writer := log.LogWriter{
    Log:   "info-webhook-url",  // Log/Info messages
    Error: "error-webhook-url", // Error messages
    Level: log.LevelInfo,
}
logger := log.Default().WithWriter(writer)
```

## Error Handling

By default, failed Slack posts do not interrupt your program. The most recent
send error is recorded on the logger and can be inspected via `Err()`:

```go
logger.Info("message")
if err := logger.Err(); err != nil {
    // Handle error
}
```

The stored error is **sticky**: once a send fails, `Err()` keeps returning that
error even after later successful sends, until you reset it with `ClearErr()`:

```go
if logger.Err() != nil {
    logger.ClearErr() // reset after handling
}
```

All loggers are safe for concurrent use by multiple goroutines.

## Notes

- Messages are automatically prefixed with their log level (ERRO, WARN, INFO, DEBG, TRCE)
- The package uses HTTP POST requests to send messages to Slack
- Log levels are hierarchical - setting a level will include all higher priority levels
- Global prefix is prepended to all messages

## License

Released under the [0BSD](LICENSE) license.
