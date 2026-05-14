package metrics_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/logsnip/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	snap := c.Snapshot()
	if snap.LinesRead != 0 || snap.LinesMatched != 0 || snap.LinesDropped != 0 {
		t.Fatalf("expected all zero counters, got %+v", snap)
	}
}

func TestRecordRead_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordRead()
	c.RecordRead()
	if got := c.Snapshot().LinesRead; got != 2 {
		t.Fatalf("expected LinesRead=2, got %d", got)
	}
}

func TestRecordMatch_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordMatch()
	if got := c.Snapshot().LinesMatched; got != 1 {
		t.Fatalf("expected LinesMatched=1, got %d", got)
	}
}

func TestRecordDrop_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordDrop()
	c.RecordDrop()
	c.RecordDrop()
	if got := c.Snapshot().LinesDropped; got != 3 {
		t.Fatalf("expected LinesDropped=3, got %d", got)
	}
}

func TestSnapshot_IsImmutable(t *testing.T) {
	c := metrics.New()
	c.RecordRead()
	snap1 := c.Snapshot()
	c.RecordRead()
	snap2 := c.Snapshot()
	if snap1.LinesRead == snap2.LinesRead {
		t.Fatal("expected snapshots to differ after additional RecordRead")
	}
}

func TestSnapshot_WriteTo_ValidJSON(t *testing.T) {
	c := metrics.New()
	c.RecordRead()
	c.RecordRead()
	c.RecordMatch()
	c.RecordDrop()

	var buf bytes.Buffer
	if err := c.Snapshot().WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo returned error: %v", err)
	}

	var out map[string]int64
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if out["lines_read"] != 2 {
		t.Errorf("expected lines_read=2, got %d", out["lines_read"])
	}
	if out["lines_matched"] != 1 {
		t.Errorf("expected lines_matched=1, got %d", out["lines_matched"])
	}
	if out["lines_dropped"] != 1 {
		t.Errorf("expected lines_dropped=1, got %d", out["lines_dropped"])
	}
}
