# Log Package

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

// Set global prefix
log.SetPrefix("[MyApp]")

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
- `Levelf(format string, args ...interface{})`
- `Levelln(args ...interface{})`

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
    Log     string
    Error   string
    Warning string
    Info    string
    Debug   string
    Trace   string
    Level   LogLevel
}
```

### Custom Configuration

```go
writer := log.LogWriter{
    Error: "error-webhook-url",
    Info: "info-webhook-url",
    Level: log.LevelInfo,
}
logger := log.Default().WithWriter(writer)
```

## Error Handling

By default, failed Slack posts are silently ignored. To handle errors, you can check the return value of logging methods:

        ```go
        logger.Info("message")
        if err := logger.Err(); err != nil {
            // Handle error
        }
        ```

## Notes

- Messages are automatically prefixed with their log level (ERRO, WARN, INFO, DEBG, TRCE)
- The package uses HTTP POST requests to send messages to Slack
- Log levels are hierarchical - setting a level will include all higher priority levels
- Global prefix is prepended to all messages
