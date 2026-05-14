package tailer

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

const defaultPollInterval = 200 * time.Millisecond

// Tailer reads new lines appended to a file, similar to `tail -f`.
type Tailer struct {
	path         string
	pollInterval time.Duration
}

// New creates a Tailer for the given file path.
func New(path string) *Tailer {
	return &Tailer{path: path, pollInterval: defaultPollInterval}
}

// Tail opens the file, seeks to the end, and streams new lines to the returned
// channel until ctx is cancelled or an unrecoverable error occurs.
func (t *Tailer) Tail(ctx context.Context) (<-chan string, <-chan error) {
	lines := make(chan string, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		f, err := os.Open(t.path)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		if _, err := f.Seek(0, io.SeekEnd); err != nil {
			errs <- err
			return
		}

		reader := bufio.NewReader(f)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(t.pollInterval)
					continue
				}
				errs <- err
				return
			}

			if len(line) > 0 {
				lines <- line
			}
		}
	}()

	return lines, errs
}
