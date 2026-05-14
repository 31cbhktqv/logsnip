package filter

import (
	"fmt"
	"regexp"
)

// Filter holds a compiled regex pattern used to match log lines.
type Filter struct {
	pattern *regexp.Regexp
}

// New creates a new Filter from the given regex pattern string.
// Returns an error if the pattern is invalid.
func New(pattern string) (*Filter, error) {
	if pattern == "" {
		return &Filter{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return &Filter{pattern: re}, nil
}

// Match reports whether the given line matches the filter pattern.
// If no pattern is set, all lines match.
func (f *Filter) Match(line string) bool {
	if f.pattern == nil {
		return true
	}
	return f.pattern.MatchString(line)
}

// Pattern returns the string representation of the compiled pattern,
// or an empty string if no pattern is set.
func (f *Filter) Pattern() string {
	if f.pattern == nil {
		return ""
	}
	return f.pattern.String()
}
