package tailer_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logsnip/internal/tailer"
)

func TestTail_ReceivesAppendedLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logsnip-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	tl := tailer.New(f.Name())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lines, errs := tl.Tail(ctx)

	// Give the goroutine time to seek to EOF before we write.
	time.Sleep(50 * time.Millisecond)

	want := "hello logsnip\n"
	if _, err := f.WriteString(want); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	select {
	case got := <-lines:
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	case err := <-errs:
		t.Fatalf("unexpected error: %v", err)
	case <-ctx.Done():
		t.Fatal("timed out waiting for line")
	}
}

func TestTail_InvalidPath_ReturnsError(t *testing.T) {
	tl := tailer.New("/nonexistent/path/to/file.log")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, errs := tl.Tail(ctx)

	select {
	case err := <-errs:
		if err == nil {
			t.Fatal("expected an error for invalid path")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for error")
	}
}
