// Package pipeline wires together the tailer, filter, and output
// components into a single cohesive processing pipeline.
package pipeline

import (
	"context"
	"io"

	"github.com/yourorg/logsnip/internal/filter"
	"github.com/yourorg/logsnip/internal/output"
	"github.com/yourorg/logsnip/internal/tailer"
)

// Config holds the configuration required to build and run a Pipeline.
type Config struct {
	// FilePath is the path to the log file to tail.
	FilePath string

	// Pattern is the regex pattern used to filter log lines.
	// An empty string disables filtering (all lines pass through).
	Pattern string

	// Writer is the destination for structured JSON output.
	Writer io.Writer
}

// Pipeline coordinates reading, filtering, and writing log lines.
type Pipeline struct {
	tailer *tailer.Tailer
	filter *filter.Filter
	output *output.Writer
}

// New constructs a Pipeline from the provided Config.
// It returns an error if the tailer or filter cannot be initialised.
func New(cfg Config) (*Pipeline, error) {
	t, err := tailer.New(cfg.FilePath)
	if err != nil {
		return nil, err
	}

	f, err := filter.New(cfg.Pattern)
	if err != nil {
		return nil, err
	}

	w := output.New(cfg.Writer)

	return &Pipeline{
		tailer: t,
		filter: f,
		output: w,
	}, nil
}

// Run starts the pipeline, reading lines from the tailer, applying the
// filter, and writing matching lines as JSON until ctx is cancelled or
// the tailer signals completion via the lines channel being closed.
func (p *Pipeline) Run(ctx context.Context) error {
	lines, err := p.tailer.Tail(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-lines:
			if !ok {
				return nil
			}
			if p.filter.Match(line) {
				if werr := p.output.Write(line); werr != nil {
					return werr
				}
			}
		}
	}
}
