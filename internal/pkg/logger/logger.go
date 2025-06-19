//go:generate go run go.uber.org/mock/mockgen -package=logger -destination=logger_mock.go github.com/qskkk/git-fleet/internal/pkg/logger Service
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Level represents the log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Color constants for styled logging (Catppuccin Mocha)
const (
	ColorDebug = "#9399b2" // Mocha Overlay 2 (gray)
	ColorInfo  = "#89dceb" // Mocha Sky (cyan)
	ColorWarn  = "#f9e2af" // Mocha Yellow
	ColorError = "#f38ba8" // Mocha Red
	ColorFatal = "#eba0ac" // Mocha Maroon
	ColorText  = "#a6adc8" // Mocha Subtext 0
	ColorDim   = "#6c7086" // Mocha Overlay 0
)

// Service defines the interface for logging
type Service interface {
	Debug(ctx context.Context, msg string, fields ...interface{})
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, err error, fields ...interface{})
	Fatal(ctx context.Context, msg string, err error, fields ...interface{})
	SetLevel(level Level)
}

// Logger implements the Service interface
type Logger struct {
	level      Level
	logger     *log.Logger
	styled     bool
	debugStyle lipgloss.Style
	infoStyle  lipgloss.Style
	warnStyle  lipgloss.Style
	errorStyle lipgloss.Style
	fatalStyle lipgloss.Style
	textStyle  lipgloss.Style
	dimStyle   lipgloss.Style
}

// New creates a new logger instance
func New() Service {
	return &Logger{
		level:  INFO,
		logger: log.New(os.Stderr, "", 0), // Remove timestamp, we'll add our own styling
		styled: true,
		debugStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorDebug)).
			Faint(true),
		infoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)),
		warnStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarn)).
			Bold(true),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)).
			Bold(true),
		fatalStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorFatal)).
			Bold(true).
			Background(lipgloss.Color("#1e1e2e")),
		textStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText)),
		dimStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorDim)).
			Faint(true),
	}
}

// NewWithLevel creates a new logger with a specific level
func NewWithLevel(level Level) Service {
	logger := New().(*Logger)
	logger.level = level
	return logger
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", msg, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, msg string, fields ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", msg, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", msg, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, msg string, err error, fields ...interface{}) {
	if l.level <= ERROR {
		// Add error to fields
		allFields := append([]interface{}{"error", err}, fields...)
		l.log("ERROR", msg, allFields...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(ctx context.Context, msg string, err error, fields ...interface{}) {
	// Add error to fields
	allFields := append([]interface{}{"error", err}, fields...)
	l.log("FATAL", msg, allFields...)
	os.Exit(1)
}

// log is the internal logging method
func (l *Logger) log(level, msg string, fields ...interface{}) {
	if !l.styled {
		// Fallback to basic logging if styling is disabled
		var fieldStr string
		if len(fields) > 0 {
			fieldStr = " "
			for i := 0; i < len(fields); i += 2 {
				if i+1 < len(fields) {
					fieldStr += fmt.Sprintf("%v=%v ", fields[i], fields[i+1])
				} else {
					fieldStr += fmt.Sprintf("%v=<missing_value> ", fields[i])
				}
			}
		}
		l.logger.Printf("[%s] %s%s", level, msg, fieldStr)
		return
	}

	// Styled logging
	var levelStyle lipgloss.Style
	var symbol string

	switch level {
	case "DEBUG":
		levelStyle = l.debugStyle
		symbol = "🔍"
	case "INFO":
		levelStyle = l.infoStyle
		symbol = "ℹ️"
	case "WARN":
		levelStyle = l.warnStyle
		symbol = "⚠️"
	case "ERROR":
		levelStyle = l.errorStyle
		symbol = "❌"
	case "FATAL":
		levelStyle = l.fatalStyle
		symbol = "💀"
	default:
		levelStyle = l.textStyle
		symbol = "•"
	}

	// Format message with styling
	styledLevel := levelStyle.Render(fmt.Sprintf("%s %s", symbol, level))
	styledMsg := l.textStyle.Render(msg)

	// Format fields more elegantly
	var fieldStr string
	if len(fields) > 0 {
		var parts []string
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				key := l.dimStyle.Render(fmt.Sprintf("%v", fields[i]))
				value := l.textStyle.Render(fmt.Sprintf("%v", fields[i+1]))
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			} else {
				key := l.dimStyle.Render(fmt.Sprintf("%v", fields[i]))
				value := l.errorStyle.Render("<missing_value>")
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			}
		}
		fieldStr = l.dimStyle.Render(" │ ") + strings.Join(parts, l.dimStyle.Render(" │ "))
	}

	// Print the styled log
	fmt.Fprintf(os.Stderr, "%s %s%s\n", styledLevel, styledMsg, fieldStr)
}
