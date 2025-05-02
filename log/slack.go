package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// LogWriter represents a writer for logging messages to Slack.
// It contains the webhook URLs for different log levels and the log level itself.
type LogWriter struct {
	Log     string
	Error   string
	Warning string
	Info    string
	Debug   string
	Trace   string

	prefix string
	Level  LogLevel
}

// LogLevel represents the log level for the LogWriter, providing type safety.
type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

type Logger struct {
	Writer LogWriter

	err error
}

func (l *Logger) SetPrefix(p string) {
	l.Writer.prefix = p
}

// Err returns the error for the Logger.
func (l *Logger) Err() error {
	return l.err
}

var std = New("")

// Flags and prefix control
var (
	flags int
	mu    sync.RWMutex
)

// SetFlags sets the logging flags for the package.
// Currently unused but maintained for stdlib compatibility.
func SetFlags(f int) {
	mu.Lock()
	defer mu.Unlock()
	flags = f
}

// Flags returns the current logging flags.
// Currently unused but maintained for stdlib compatibility.
func Flags() int {
	mu.RLock()
	defer mu.RUnlock()
	return flags
}

// SetPrefix sets the prefix for all log messages.
// The prefix will be prepended to all messages sent to Slack.
func SetPrefix(p string) {
	std.SetPrefix(p)
}

// Prefix returns the current log message prefix.
func Prefix() string {
	return std.Writer.prefix
}

// Default returns the default logger instance.
func Default() *Logger {
	return std
}

// info writes an info level message to Slack.
// Returns the number of bytes written and any error encountered.
func (lw LogWriter) info(p []byte) (n int, err error) {
	if lw.Level < LevelInfo {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("INFO: %s", string(buf))
	return len(p), postSlack(lw.Log, strLine, lw.prefix)
}

// error writes an error level message to Slack.
// Returns the number of bytes written and any error encountered.
func (lw LogWriter) error(p []byte) (n int, err error) {
	if lw.Level < LevelError {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("ERRO: %s", string(buf))
	return len(p), postSlack(lw.Error, strLine, lw.prefix)
}

// warning writes a warning level message to Slack.
// Returns the number of bytes written and any error encountered.
func (lw LogWriter) warning(p []byte) (n int, err error) {
	if lw.Level < LevelWarning {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("WARN: %s", string(buf))
	return len(p), postSlack(lw.Warning, strLine, lw.prefix)
}

// debug writes a debug level message to Slack.
// Returns the number of bytes written and any error encountered.
func (lw LogWriter) debug(p []byte) (n int, err error) {
	if lw.Level < LevelDebug {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("DEBG: %s", string(buf))
	return len(p), postSlack(lw.Debug, strLine, lw.prefix)
}

// trace writes a trace level message to Slack.
// Returns the number of bytes written and any error encountered.
func (lw LogWriter) trace(p []byte) (n int, err error) {
	if lw.Level < LevelTrace {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("TRCE: %s", string(buf))
	return len(p), postSlack(lw.Trace, strLine, lw.prefix)
}

// log writes a message at the default info level to Slack.
// Returns the number of bytes written and any error encountered.
func (lw LogWriter) log(p []byte) (n int, err error) {
	return lw.info(p)
}

// Write implements the io.Writer interface for LogWriter.
// Writes the message to Slack at the default info level.
func (lw LogWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := string(buf)
	return len(p), postSlack(lw.Log, strLine, lw.prefix)
}

// postSlack sends a message to a Slack webhook.
// Returns any error encountered during the HTTP request.
func postSlack(webhook, text, prefix string) error {
	if prefix != "" {
		text = prefix + text
	}
	values := map[string]string{"text": text}
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonValue))
	return err
}

// New creates a new Logger with the specified webhook URL.
// The webhook URL will be used for all log levels.
func New(webhookLink string) *Logger {
	return &Logger{
		Writer: LogWriter{
			Log:     webhookLink,
			Error:   webhookLink,
			Warning: webhookLink,
			Info:    webhookLink,
			Debug:   webhookLink,
			Trace:   webhookLink,
			Level:   LevelTrace,
		},
	}
}

// WithLevel returns a new Logger with the specified log level.
func WithLevel(level LogLevel) Logger {
	return std.WithLevel(level)
}

// WithLevel sets the log level for the Logger.
func (l *Logger) WithLevel(level LogLevel) Logger {
	l.Writer.Level = level
	return *l
}

// WithWriter returns a new Logger with the specified LogWriter.
func WithWriter(w LogWriter) Logger {
	return std.WithWriter(w)
}

// WithWriter sets the LogWriter for the Logger.
func (l *Logger) WithWriter(w LogWriter) Logger {
	l.Writer = w
	return *l
}

// Log writes a message at the default info level.
func Log(msg string) {
	std.Log(msg)
}

// Log writes a message at the default info level.
func (l *Logger) Log(msg string) {
	l.Writer.log([]byte(msg))
}

// Logf writes a formatted message at the default info level.
func Logf(msg string, args ...interface{}) {
	std.Logf(msg, args...)
}

// Logf writes a formatted message at the default info level.
func (l *Logger) Logf(msg string, args ...interface{}) {
	l.Writer.log([]byte(fmt.Sprintf(msg, args...)))
}

// Logln writes a message at the default info level with a newline.
func Logln(args ...interface{}) {
	std.Logln(args...)
}

// Logln writes a message at the default info level with a newline.
func (l *Logger) Logln(args ...interface{}) {
	l.Writer.log([]byte(fmt.Sprintln(args...)))
}

// Error writes an error level message.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Error writes an error level message.
func (l *Logger) Error(args ...interface{}) {
	l.Writer.error([]byte(fmt.Sprintln(args...)))
}

// Errorf writes a formatted error level message.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Errorf writes a formatted error level message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Writer.error([]byte(fmt.Sprintf(format, args...)))
}

// Errorln writes an error level message with a newline.
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

// Errorln writes an error level message with a newline.
func (l *Logger) Errorln(args ...interface{}) {
	l.Writer.error([]byte(fmt.Sprintln(args...)))
}

// Warning writes a warning level message.
func Warning(warning string) {
	std.Warning(warning)
}

// Warning writes a warning level message.
func (l *Logger) Warning(warning string) {
	l.Writer.warning([]byte(warning))
}

// Warningf writes a formatted warning level message.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Warningf writes a formatted warning level message.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Writer.warning([]byte(fmt.Sprintf(format, args...)))
}

// Warningln writes a warning level message with a newline.
func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

// Warningln writes a warning level message with a newline.
func (l *Logger) Warningln(args ...interface{}) {
	l.Writer.warning([]byte(fmt.Sprintln(args...)))
}

// Info writes an info level message.
func Info(info string) {
	std.Info(info)
}

// Info writes an info level message.
func (l *Logger) Info(info string) {
	_, err := l.Writer.info([]byte(info))
	if err != nil {
		l.err = err
	}
}

// Infof writes a formatted info level message.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Infof writes a formatted info level message.
func (l *Logger) Infof(format string, args ...interface{}) {
	_, err := l.Writer.info([]byte(fmt.Sprintf(format, args...)))
	if err != nil {
		l.err = err
	}
}

// Infoln writes an info level message with a newline.
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

// Infoln writes an info level message with a newline.
func (l *Logger) Infoln(args ...interface{}) {
	_, err := l.Writer.info([]byte(fmt.Sprintln(args...)))
	if err != nil {
		l.err = err
	}
}

// Debug writes a debug level message.
func Debug(debug string) {
	std.Debug(debug)
}

// Debug writes a debug level message.
func (l *Logger) Debug(debug string) {
	_, err := l.Writer.debug([]byte(debug))
	if err != nil {
		l.err = err
	}
}

// Debugf writes a formatted debug level message.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	_, err := l.Writer.debug([]byte(fmt.Sprintf(format, args...)))
	if err != nil {
		l.err = err
	}
}

// Debugln writes a debug level message with a newline.
func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

// Debugln writes a debug level message with a newline.
func (l *Logger) Debugln(args ...interface{}) {
	_, err := l.Writer.debug([]byte(fmt.Sprintln(args...)))
	if err != nil {
		l.err = err
	}
}

// Trace writes a trace level message.
func Trace(trace string) {
	std.Trace(trace)
}

// Trace writes a trace level message.
func (l *Logger) Trace(trace string) {
	_, err := l.Writer.trace([]byte(trace))
	if err != nil {
		l.err = err
	}
}

// Tracef writes a formatted trace level message.
func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

// Tracef writes a formatted trace level message.
func (l *Logger) Tracef(format string, args ...interface{}) {
	_, err := l.Writer.trace([]byte(fmt.Sprintf(format, args...)))
	if err != nil {
		l.err = err
	}
}

// Traceln writes a trace level message with a newline.
func Traceln(args ...interface{}) {
	std.Traceln(args...)
}

// Traceln writes a trace level message with a newline.
func (l *Logger) Traceln(args ...interface{}) {
	_, err := l.Writer.trace([]byte(fmt.Sprintln(args...)))
	if err != nil {
		l.err = err
	}
}

// Basic logging functions
func Print(v ...interface{}) {
	std.Log(fmt.Sprint(v...))
}

// Printf writes a formatted message at the default info level.
func Printf(format string, v ...interface{}) {
	std.Logf(format, v...)
}

// Println writes a message at the default info level with a newline.
func Println(v ...interface{}) {
	std.Logln(v...)
}

// Fatal writes a message at the default error level.
// Subsequently, it calls os.Exit(1).
func Fatal(v ...interface{}) {
	std.Error(fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf writes a formatted message at the default error level.
func Fatalf(format string, v ...interface{}) {
	std.Errorf(format, v...)
	os.Exit(1)
}

// Fatalln writes a message at the default error level with a newline.
func Fatalln(v ...interface{}) {
	std.Errorln(v...)
	os.Exit(1)
}

// Panic writes a message at the default error level.
// Subsequently, it panics with the message.
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Error(s)
	panic(s)
}

// Panicf writes a formatted message at the default error level.
// Subsequently, it panics with the formatted message.
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Error(s)
	panic(s)
}

// Panicln writes a message at the default error level with a newline.
// Subsequently, it panics with the message.
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Error(s)
	panic(s)
}
