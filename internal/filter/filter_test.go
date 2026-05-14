package filter_test

import (
	"testing"

	"github.com/yourorg/logsnip/internal/filter"
)

func TestNew_ValidPattern(t *testing.T) {
	f, err := filter.New(`ERROR|WARN`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := filter.New(`[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	f, err := filter.New("")
	if err != nil {
		t.Fatalf("expected no error for empty pattern, got %v", err)
	}
	if !f.Match("any line should match") {
		t.Error("empty pattern should match all lines")
	}
}

func TestMatch_Positive(t *testing.T) {
	f, _ := filter.New(`ERROR`)
	lines := []string{
		"2024-01-01 ERROR something went wrong",
		"ERROR: disk full",
	}
	for _, line := range lines {
		if !f.Match(line) {
			t.Errorf("expected line %q to match", line)
		}
	}
}

func TestMatch_Negative(t *testing.T) {
	f, _ := filter.New(`ERROR`)
	lines := []string{
		"2024-01-01 INFO service started",
		"DEBUG: cache hit",
	}
	for _, line := range lines {
		if f.Match(line) {
			t.Errorf("expected line %q not to match", line)
		}
	}
}

func TestPattern_ReturnsString(t *testing.T) {
	pattern := `WARN|ERROR`
	f, _ := filter.New(pattern)
	if f.Pattern() != pattern {
		t.Errorf("expected pattern %q, got %q", pattern, f.Pattern())
	}
}

func TestPattern_EmptyWhenNoPattern(t *testing.T) {
	f, _ := filter.New("")
	if f.Pattern() != "" {
		t.Errorf("expected empty pattern string, got %q", f.Pattern())
	}
}
