// Package metrics implements lightweight, thread-safe counters for the
// logsnip pipeline.
//
// It tracks three categories of events:
//
//   - LinesRead    – every line consumed from the tailed file.
//   - LinesMatched – lines that passed the regex filter and were written
//     to the JSON output.
//   - LinesDropped – lines that were read but did not match the filter.
//
// Counters are updated via atomic operations so they are safe to call
// concurrently from the pipeline goroutine without additional locking.
//
// A Snapshot can be obtained at any time and serialised to JSON for
// inclusion in monitoring dashboards or health-check endpoints.
package metrics
