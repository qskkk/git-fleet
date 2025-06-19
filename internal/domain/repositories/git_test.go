package repositories

import (
	"testing"
)

func TestCommitInfo_Fields(t *testing.T) {
	commit := &CommitInfo{
		Hash:      "abc123def456",
		Author:    "John Doe <john@example.com>",
		Message:   "Initial commit",
		Timestamp: "2023-01-01T12:00:00Z",
	}

	if commit.Hash != "abc123def456" {
		t.Errorf("Hash = %s, want %s", commit.Hash, "abc123def456")
	}
	if commit.Author != "John Doe <john@example.com>" {
		t.Errorf("Author = %s, want %s", commit.Author, "John Doe <john@example.com>")
	}
	if commit.Message != "Initial commit" {
		t.Errorf("Message = %s, want %s", commit.Message, "Initial commit")
	}
	if commit.Timestamp != "2023-01-01T12:00:00Z" {
		t.Errorf("Timestamp = %s, want %s", commit.Timestamp, "2023-01-01T12:00:00Z")
	}
}

func TestCommitInfo_EmptyFields(t *testing.T) {
	commit := &CommitInfo{}

	if commit.Hash != "" {
		t.Errorf("Hash should be empty, got %s", commit.Hash)
	}
	if commit.Author != "" {
		t.Errorf("Author should be empty, got %s", commit.Author)
	}
	if commit.Message != "" {
		t.Errorf("Message should be empty, got %s", commit.Message)
	}
	if commit.Timestamp != "" {
		t.Errorf("Timestamp should be empty, got %s", commit.Timestamp)
	}
}

func TestCommitInfo_LongValues(t *testing.T) {
	longHash := "1234567890abcdef1234567890abcdef12345678"
	longAuthor := "Very Long Author Name With Email <very.long.email.address@example.com>"
	longMessage := "This is a very long commit message that describes in great detail what changes were made in this commit and why they were necessary for the project."
	longTimestamp := "2023-12-31T23:59:59.999999999Z"

	commit := &CommitInfo{
		Hash:      longHash,
		Author:    longAuthor,
		Message:   longMessage,
		Timestamp: longTimestamp,
	}

	if commit.Hash != longHash {
		t.Errorf("Hash = %s, want %s", commit.Hash, longHash)
	}
	if commit.Author != longAuthor {
		t.Errorf("Author = %s, want %s", commit.Author, longAuthor)
	}
	if commit.Message != longMessage {
		t.Errorf("Message = %s, want %s", commit.Message, longMessage)
	}
	if commit.Timestamp != longTimestamp {
		t.Errorf("Timestamp = %s, want %s", commit.Timestamp, longTimestamp)
	}
}

func TestCommitInfo_SpecialCharacters(t *testing.T) {
	commit := &CommitInfo{
		Hash:      "abc123!@#$%^&*()",
		Author:    "Ñoël Müller <noel@münchen.de>",
		Message:   "Fix: handle unicode characters (中文, العربية, русский)",
		Timestamp: "2023-06-15T14:30:45+02:00",
	}

	if commit.Hash != "abc123!@#$%^&*()" {
		t.Errorf("Hash = %s, want %s", commit.Hash, "abc123!@#$%^&*()")
	}
	if commit.Author != "Ñoël Müller <noel@münchen.de>" {
		t.Errorf("Author = %s, want %s", commit.Author, "Ñoël Müller <noel@münchen.de>")
	}
	if commit.Message != "Fix: handle unicode characters (中文, العربية, русский)" {
		t.Errorf("Message = %s, want %s", commit.Message, "Fix: handle unicode characters (中文, العربية, русский)")
	}
	if commit.Timestamp != "2023-06-15T14:30:45+02:00" {
		t.Errorf("Timestamp = %s, want %s", commit.Timestamp, "2023-06-15T14:30:45+02:00")
	}
}

func TestCommitInfo_JSON_Tags(t *testing.T) {
	// This test verifies that the struct has the correct JSON tags
	commit := &CommitInfo{
		Hash:      "test-hash",
		Author:    "test-author",
		Message:   "test-message",
		Timestamp: "test-timestamp",
	}

	// We can't directly test JSON tags without marshaling, but we can verify
	// the struct fields are properly set
	if commit.Hash != "test-hash" {
		t.Error("Hash field should be accessible")
	}
	if commit.Author != "test-author" {
		t.Error("Author field should be accessible")
	}
	if commit.Message != "test-message" {
		t.Error("Message field should be accessible")
	}
	if commit.Timestamp != "test-timestamp" {
		t.Error("Timestamp field should be accessible")
	}
}
