package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogWriter struct {
	Log     string
	Error   string
	Warning string
	Info    string
	Debug   string
	Trace   string

	Level LogLevel
}

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
}

var std = New("")

func Default() *Logger {
	return std
}

func (lw LogWriter) info(p []byte) (n int, err error) {
	if lw.Level < LevelInfo {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("INFO: %s", string(buf))
	return len(p), postSlack(lw.Log, strLine)
}

func (lw LogWriter) error(p []byte) (n int, err error) {
	if lw.Level < LevelError {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("ERRO: %s", string(buf))
	return len(p), postSlack(lw.Error, strLine)
}

func (lw LogWriter) warning(p []byte) (n int, err error) {
	if lw.Level < LevelWarning {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("WARN: %s", string(buf))
	return len(p), postSlack(lw.Warning, strLine)
}

func (lw LogWriter) debug(p []byte) (n int, err error) {
	if lw.Level < LevelDebug {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("DEBG: %s", string(buf))
	return len(p), postSlack(lw.Debug, strLine)
}

func (lw LogWriter) trace(p []byte) (n int, err error) {
	if lw.Level < LevelTrace {
		return
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := fmt.Sprintf("TRCE: %s", string(buf))
	return len(p), postSlack(lw.Trace, strLine)
}

func (lw LogWriter) log(p []byte) (n int, err error) {
	return lw.info(p)
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	strLine := string(buf)
	return len(p), postSlack(lw.Log, strLine)
}

func postSlack(webhook, text string) error {
	values := map[string]string{"text": text}
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonValue))
	return err
}

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

func WithLevel(level LogLevel) Logger {
	return std.WithLevel(level)
}

func (l *Logger) WithLevel(level LogLevel) Logger {
	l.Writer.Level = level
	return *l
}

func WithWriter(w LogWriter) Logger {
	return std.WithWriter(w)
}

func (l *Logger) WithWriter(w LogWriter) Logger {
	l.Writer = w
	return *l
}

func Log(msg string) {
	std.Log(msg)
}

func (l *Logger) Log(msg string) {
	l.Writer.log([]byte(msg))
}

func Logf(msg string, args ...interface{}) {
	std.Logf(msg, args...)
}

func (l *Logger) Logf(msg string, args ...interface{}) {
	l.Writer.log([]byte(fmt.Sprintf(msg, args...)))
}

func Logln(args ...interface{}) {
	std.Logln(args...)
}

func (l *Logger) Logln(args ...interface{}) {
	l.Writer.log([]byte(fmt.Sprintln(args...)))
}

func Error(args ...interface{}) {
	std.Error(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Writer.error([]byte(fmt.Sprintln(args...)))
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Writer.error([]byte(fmt.Sprintf(format, args...)))
}

func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Writer.error([]byte(fmt.Sprintln(args...)))
}

func Warning(warning string) {
	std.Warning(warning)
}

func (l *Logger) Warning(warning string) {
	l.Writer.warning([]byte(warning))
}

func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Writer.warning([]byte(fmt.Sprintf(format, args...)))
}

func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.Writer.warning([]byte(fmt.Sprintln(args...)))
}

func Info(info string) {
	std.Info(info)
}

func (l *Logger) Info(info string) {
	_, err := l.Writer.info([]byte(info))
	if err != nil {
		fmt.Println(err)
	}
}

func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Writer.info([]byte(fmt.Sprintf(format, args...)))
}

func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Writer.info([]byte(fmt.Sprintln(args...)))
}

func Debug(debug string) {
	std.Debug(debug)
}

func (l *Logger) Debug(debug string) {
	l.Writer.debug([]byte(debug))
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Writer.debug([]byte(fmt.Sprintf(format, args...)))
}

func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Writer.debug([]byte(fmt.Sprintln(args...)))
}

func Trace(trace string) {
	std.Trace(trace)
}

func (l *Logger) Trace(trace string) {
	l.Writer.trace([]byte(trace))
}

func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Writer.trace([]byte(fmt.Sprintf(format, args...)))
}

func Traceln(args ...interface{}) {
	std.Traceln(args...)
}

func (l *Logger) Traceln(args ...interface{}) {
	l.Writer.trace([]byte(fmt.Sprintln(args...)))
}
