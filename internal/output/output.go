package output

import (
	"encoding/json"
	"io"
	"time"
)

// Entry represents a structured log entry for JSON output.
type Entry struct {
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Line      string `json:"line"`
	Matched   string `json:"matched,omitempty"`
}

// Writer wraps an io.Writer and writes structured JSON log entries.
type Writer struct {
	w       io.Writer
	encoder *json.Encoder
}

// New creates a new JSON Writer that writes to w.
func New(w io.Writer) *Writer {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Writer{w: w, encoder: enc}
}

// Write encodes a single log entry as a JSON line to the underlying writer.
func (w *Writer) Write(source, line, matched string) error {
	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Source:    source,
		Line:      line,
		Matched:   matched,
	}
	return w.encoder.Encode(entry)
}
