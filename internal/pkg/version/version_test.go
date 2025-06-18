package version

import (
	"testing"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()

	// Test that all fields are set
	if info.Version == "" {
		t.Error("Version should not be empty")
	}
	// GitCommit and BuildDate can be empty, so we just check they're strings
	if info.GitCommit != "" && len(info.GitCommit) == 0 {
		t.Error("GitCommit should be a string")
	}
	if info.BuildDate != "" && len(info.BuildDate) == 0 {
		t.Error("BuildDate should be a string")
	}
	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
}

func TestInfo_Fields(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2023-01-01",
		GoVersion: "go1.20",
	}

	if info.Version != "1.0.0" {
		t.Errorf("Version = %s, want %s", info.Version, "1.0.0")
	}
	if info.GitCommit != "abc123" {
		t.Errorf("GitCommit = %s, want %s", info.GitCommit, "abc123")
	}
	if info.BuildDate != "2023-01-01" {
		t.Errorf("BuildDate = %s, want %s", info.BuildDate, "2023-01-01")
	}
	if info.GoVersion != "go1.20" {
		t.Errorf("GoVersion = %s, want %s", info.GoVersion, "go1.20")
	}
}

func TestGetVersion(t *testing.T) {
	result := GetVersion()
	if result == "" {
		t.Error("GetVersion() should not return empty string")
	}
	// Check that it starts with "GitFleet"
	if !contains(result, "GitFleet") {
		t.Error("GetVersion() should contain 'GitFleet'")
	}
}

func TestGetVersionLong(t *testing.T) {
	result := GetVersionLong()
	if result == "" {
		t.Error("GetVersionLong() should not return empty string")
	}
	// Check that it contains "GitFleet"
	if !contains(result, "GitFleet") {
		t.Error("GetVersionLong() should contain 'GitFleet'")
	}
}

func TestDefaultValues(t *testing.T) {
	// Test default values set by global variables
	if Version == "" {
		t.Error("Version global variable should have a default value")
	}
	// GitCommit and BuildDate can be empty by default
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
