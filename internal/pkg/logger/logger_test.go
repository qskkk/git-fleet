package logger

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Check if it implements the interface
	if _, ok := logger.(*Logger); !ok {
		t.Error("New should return a *Logger")
	}
}

func TestNewWithLevel(t *testing.T) {
	logger := NewWithLevel(DEBUG)
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Check if it implements the interface
	if _, ok := logger.(*Logger); !ok {
		t.Error("NewWithLevel should return a *Logger")
	}
}

func TestLogger_SetLevel(t *testing.T) {
	logger := &Logger{level: INFO}

	logger.SetLevel(DEBUG)
	if logger.level != DEBUG {
		t.Errorf("Level = %d, want %d", logger.level, DEBUG)
	}

	logger.SetLevel(ERROR)
	if logger.level != ERROR {
		t.Errorf("Level = %d, want %d", logger.level, ERROR)
	}
}

func TestLogger_Debug(t *testing.T) {
	logger := New()
	ctx := context.Background()

	// Should not panic
	logger.Debug(ctx, "debug message")
	logger.Debug(ctx, "debug message with fields", "key", "value")
}

func TestLogger_Info(t *testing.T) {
	logger := New()
	ctx := context.Background()

	// Should not panic
	logger.Info(ctx, "info message")
	logger.Info(ctx, "info message with fields", "key", "value")
}

func TestLogger_Warn(t *testing.T) {
	logger := New()
	ctx := context.Background()

	// Should not panic
	logger.Warn(ctx, "warn message")
	logger.Warn(ctx, "warn message with fields", "key", "value")
}

func TestLogger_Error(t *testing.T) {
	logger := New()
	ctx := context.Background()

	// Should not panic
	logger.Error(ctx, "error message", nil)
	logger.Error(ctx, "error message with fields", nil, "key", "value")
}

func TestLevel_Constants(t *testing.T) {
	if DEBUG != 0 {
		t.Errorf("DEBUG = %d, want %d", DEBUG, 0)
	}
	if INFO != 1 {
		t.Errorf("INFO = %d, want %d", INFO, 1)
	}
	if WARN != 2 {
		t.Errorf("WARN = %d, want %d", WARN, 2)
	}
	if ERROR != 3 {
		t.Errorf("ERROR = %d, want %d", ERROR, 3)
	}
}

func TestLogger_Fields(t *testing.T) {
	logger := &Logger{
		level:  INFO,
		styled: true,
	}

	if logger.level != INFO {
		t.Errorf("level = %d, want %d", logger.level, INFO)
	}
	if !logger.styled {
		t.Error("styled should be true")
	}
}
