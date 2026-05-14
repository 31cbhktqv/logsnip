// Package filter provides regex-based line filtering for logsnip.
//
// It exposes a Filter type that wraps a compiled regular expression
// and can be used to selectively pass log lines to downstream processors.
//
// Usage:
//
//	f, err := filter.New(`ERROR|WARN`)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	scanner := bufio.NewScanner(os.Stdin)
//	for scanner.Scan() {
//		line := scanner.Text()
//		if f.Match(line) {
//			// process matched line
//			fmt.Println(line)
//		}
//	}
//
// An empty pattern string causes every line to match, making the filter
// effectively a no-op pass-through. This is the default behaviour when
// no --filter flag is provided by the user.
package filter
