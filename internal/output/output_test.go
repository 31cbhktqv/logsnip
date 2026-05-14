package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logsnip/internal/output"
)

func TestNew_ReturnsWriter(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)
	if w == nil {
		t.Fatal("expected non-nil Writer")
	}
}

func TestWrite_ProducesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	if err := w.Write("app.log", "error occurred", "error"); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	var entry output.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if entry.Source != "app.log" {
		t.Errorf("expected source 'app.log', got %q", entry.Source)
	}
	if entry.Line != "error occurred" {
		t.Errorf("expected line 'error occurred', got %q", entry.Line)
	}
	if entry.Matched != "error" {
		t.Errorf("expected matched 'error', got %q", entry.Matched)
	}
	if entry.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestWrite_EmptyMatchedOmitted(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	if err := w.Write("app.log", "info line", ""); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	if strings.Contains(buf.String(), "matched") {
		t.Error("expected 'matched' field to be omitted when empty")
	}
}

func TestWrite_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := w.Write("test.log", l, ""); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	}

	parts := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(parts) != 3 {
		t.Errorf("expected 3 JSON lines, got %d", len(parts))
	}
}
