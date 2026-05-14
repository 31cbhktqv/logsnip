// Package metrics provides runtime counters for tracking pipeline activity
// such as lines read, lines matched, and lines dropped during filtering.
package metrics

import (
	"encoding/json"
	"io"
	"sync/atomic"
)

// Counters holds atomic counters for pipeline telemetry.
type Counters struct {
	LinesRead    atomic.Int64
	LinesMatched atomic.Int64
	LinesDropped atomic.Int64
}

// New returns a zeroed Counters instance ready for use.
func New() *Counters {
	return &Counters{}
}

// RecordRead increments the lines-read counter by 1.
func (c *Counters) RecordRead() {
	c.LinesRead.Add(1)
}

// RecordMatch increments the lines-matched counter by 1.
func (c *Counters) RecordMatch() {
	c.LinesMatched.Add(1)
}

// RecordDrop increments the lines-dropped counter by 1.
func (c *Counters) RecordDrop() {
	c.LinesDropped.Add(1)
}

// Snapshot returns a point-in-time copy of the current counter values.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:    c.LinesRead.Load(),
		LinesMatched: c.LinesMatched.Load(),
		LinesDropped: c.LinesDropped.Load(),
	}
}

// Snapshot is an immutable point-in-time view of pipeline counters.
type Snapshot struct {
	LinesRead    int64 `json:"lines_read"`
	LinesMatched int64 `json:"lines_matched"`
	LinesDropped int64 `json:"lines_dropped"`
}

// WriteTo serialises the snapshot as a single JSON object to w.
func (s Snapshot) WriteTo(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s)
}
