package pipeline_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logsnip/internal/pipeline"
)

func writeTempLog(t *testing.T, lines ...string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logsnip-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	path := writeTempLog(t, "hello")
	_, err := pipeline.New(pipeline.Config{
		FilePath: path,
		Pattern:  "[", // invalid regex
		Writer:   &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern, got nil")
	}
}

func TestNew_InvalidPath_ReturnsError(t *testing.T) {
	_, err := pipeline.New(pipeline.Config{
		FilePath: "/nonexistent/path/to/file.log",
		Pattern:  "",
		Writer:   &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRun_FilteredLinesProduceJSON(t *testing.T) {
	path := writeTempLog(t,
		"INFO service started",
		"DEBUG health check ok",
		"ERROR connection refused",
	)

	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{
		FilePath: path,
		Pattern:  "ERROR",
		Writer:   &buf,
	})
	if err != nil {
		t.Fatalf("pipeline.New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx); err != nil && err != context.DeadlineExceeded {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	dec := json.NewDecoder(&buf)
	var entry map[string]string
	if err := dec.Decode(&entry); err != nil {
		t.Fatalf("decode JSON output: %v", err)
	}
	if entry["line"] != "ERROR connection refused" {
		t.Errorf("unexpected line value: %q", entry["line"])
	}
}
