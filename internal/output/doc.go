// Package output provides structured JSON output writing for logsnip.
//
// Each log line that passes through the filter pipeline is serialised as a
// newline-delimited JSON object (JSON-L / NDJSON), making the stream easy to
// pipe into downstream monitoring or aggregation tools such as jq, Elasticsearch
// Logstash, or Fluentd.
//
// Example output line:
//
//	{"timestamp":"2024-01-15T10:23:45.123456789Z","source":"app.log","line":"ERROR: disk full","matched":"ERROR"}
//
// The "matched" field contains the first regex match found in the line and is
// omitted entirely when no pattern is active or no match was captured.
package output
