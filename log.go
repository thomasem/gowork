package gowork

import (
	"fmt"
	"time"
)

type Level int

const (
	LogDebug Level = iota
	LogInfo
	LogWarn
	LogError
)

func (l Level) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Log struct {
	Level   Level
	Message string
}

type Logger interface {
	Log(Log)
}

type DefaultLogger struct {
	Level Level
}

func (l *DefaultLogger) Log(log Log) {
	if log.Level >= l.Level {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s: %s\n", log.Level.String(), timeStr, log.Message)
	}
}

func NewDefaultLogger(level Level) Logger {
	return &DefaultLogger{Level: level}
}

type NullLogger struct{}

func (n *NullLogger) Log(log Log) {}

func NewNullLogger() Logger {
	return &NullLogger{}
}
