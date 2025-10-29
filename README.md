# Linkchex

A high-performance CLI tool written in Go that discovers sitemaps, crawls pages, and validates all links with maximum speed and efficiency.

## Current Status: Phase 1 Complete ✓

**Phase 1: Project Setup & Sitemap Discovery** - All tasks completed!

### Features Implemented
- ✓ Project initialization with Go modules
- ✓ Basic CLI framework with flags
- ✓ Automatic sitemap discovery from base URLs
- ✓ robots.txt parsing for sitemap directives
- ✓ XML sitemap parsing
- ✓ Sitemap index file support (nested sitemaps)
- ✓ Local and remote sitemap support
- ✓ Text output with discovered URLs

## Installation

### Build from source
```bash
go build -o linkchex ./cmd/linkchex
```

## Usage

### Basic Usage

```bash
# Use a local sitemap file
./linkchex --sitemap test-sitemap.xml

# Discover sitemap from base URL
./linkchex --url https://example.com

# Verbose output
./linkchex --sitemap test-sitemap.xml --verbose
```

### CLI Flags

```
  -url string
      Base URL to discover sitemap from
  -sitemap string
      Direct URL or path to sitemap file
  -concurrency int
      Number of concurrent workers (default 50)
  -timeout int
      Request timeout in seconds (default 10)
  -format string
      Output format (text, json, csv) (default "text")
  -output string
      Output file path (default: stdout)
  -verbose
      Enable verbose output
  -version
      Show version information
```

## Testing

Test files are included:
- `test-sitemap.xml` - Sample sitemap with 4 URLs
- `test-sitemap-index.xml` - Sample sitemap index file

```bash
# Test with local sitemap
./linkchex --sitemap test-sitemap.xml --verbose

# Test with sitemap index
./linkchex --sitemap test-sitemap-index.xml --verbose
```

## Project Structure

```
linkchex/
├── cmd/
│   └── linkchex/
│       └── main.go           # CLI entry point
├── internal/
│   ├── sitemap/
│   │   ├── discover.go       # Sitemap discovery logic
│   │   └── parser.go         # XML parsing logic
│   ├── validator/            # (Phase 2)
│   └── fetcher/              # (Phase 2)
├── go.mod
├── PROJECT-PLAN.md
└── README.md
```

## Next Steps: Phase 2

**Phase 2: Basic Link Validation** will implement:
- HTTP client setup with timeouts and retries
- Page fetching functionality
- HTML link extraction (a, img, link, script tags)
- Basic link validation
- Status code checking

See [PROJECT-PLAN.md](PROJECT-PLAN.md) for the complete implementation roadmap.

## Technology Stack

- **Language**: Go 1.25.3
- **Standard Library**: net/http, encoding/xml
- **Dependencies**: None yet (minimal external dependencies)

## Performance Targets

- 500 pages: 10-30 seconds (target for Phase 3)
- 1000 pages: 20-60 seconds (target for Phase 3)
- Concurrency: 50-200 workers (Phase 3)
- Memory: <500MB for typical workloads

## License

MIT
