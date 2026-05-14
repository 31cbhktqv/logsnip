# logsnip

Lightweight log tail utility with regex filtering and structured JSON output for piping into monitoring tools.

---

## Installation

```bash
go install github.com/yourusername/logsnip@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logsnip.git && cd logsnip && go build -o logsnip .
```

---

## Usage

Tail a log file and filter lines matching a regex pattern:

```bash
logsnip -f /var/log/app.log -p "ERROR|WARN"
```

Output structured JSON for piping into monitoring tools:

```bash
logsnip -f /var/log/app.log -p "ERROR" --json | jq '.message'
```

Pipe from stdin:

```bash
tail -f /var/log/syslog | logsnip -p "segfault" --json
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-f, --file` | Log file to tail | stdin |
| `-p, --pattern` | Regex filter pattern | `.*` |
| `--json` | Output as structured JSON | false |
| `-n, --lines` | Number of historical lines to show on start | `10` |

### JSON Output Format

```json
{
  "timestamp": "2024-05-01T12:34:56Z",
  "level": "ERROR",
  "message": "connection refused on port 5432",
  "source": "/var/log/app.log"
}
```

---

## License

MIT © 2024 yourusername