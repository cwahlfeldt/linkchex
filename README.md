# Linkchex

A high-performance CLI tool written in Go that discovers sitemaps, crawls pages, and validates all links with maximum speed and efficiency.

## Current Status: Phase 2 Complete ✓

**Phase 1: Project Setup & Sitemap Discovery** - ✓ Complete
**Phase 2: Basic Link Validation** - ✓ Complete

### Features Implemented

#### Phase 1
- ✓ Project initialization with Go modules
- ✓ Basic CLI framework with flags
- ✓ Automatic sitemap discovery from base URLs
- ✓ robots.txt parsing for sitemap directives
- ✓ XML sitemap parsing
- ✓ Sitemap index file support (nested sitemaps)
- ✓ Local and remote sitemap support

#### Phase 2
- ✓ HTTP client with timeout and retry logic
- ✓ Page fetching functionality
- ✓ HTML link extraction (a, img, link, script tags)
- ✓ Link validation with status code checking
- ✓ Concurrent link validation
- ✓ Broken link detection
- ✓ Multiple output formats (text, JSON, CSV)
- ✓ Internal/external link filtering
- ✓ Detailed validation reports
- ✓ File output support

## Installation

### Build from source
```bash
go build -o linkchex ./cmd/linkchex
```

## Usage

### Basic Usage

```bash
# Validate all links from a sitemap (internal links only)
./linkchex --sitemap test-sitemap.xml

# Include external links in validation
./linkchex --sitemap test-sitemap.xml --check-external

# Discover sitemap from base URL
./linkchex --url https://example.com

# Just list URLs without validating (Phase 1 behavior)
./linkchex --sitemap test-sitemap.xml --list-only

# Verbose output with detailed progress
./linkchex --sitemap test-sitemap.xml --verbose

# Output to JSON file
./linkchex --sitemap test-sitemap.xml --format json --output report.json

# Output to CSV file
./linkchex --sitemap test-sitemap.xml --format csv --output report.csv
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
  -retries int
      Maximum number of retries for failed requests (default 2)
  -check-external
      Check external links (default: internal only)
  -list-only
      Only list URLs from sitemap without validating links
  -format string
      Output format (text, json, csv) (default "text")
  -output string
      Output file path (default: stdout)
  -verbose
      Enable verbose output
  -version
      Show version information
```

### Exit Codes

- `0` - Success, all links are valid
- `1` - Failure, broken links were found or an error occurred

## Testing

Test files are included:
- `test-sitemap.xml` - Sample sitemap with 4 URLs
- `test-sitemap-index.xml` - Sample sitemap index file
- `test-google-sitemap.xml` - Real-world test with Google homepage

```bash
# Test with local sitemap (list only)
./linkchex --sitemap test-sitemap.xml --list-only

# Test with sitemap index
./linkchex --sitemap test-sitemap-index.xml --verbose

# Test link validation with a real website
./linkchex --sitemap test-google-sitemap.xml --verbose --check-external

# Test JSON output
./linkchex --sitemap test-google-sitemap.xml --format json

# Test CSV output to file
./linkchex --sitemap test-google-sitemap.xml --format csv --output results.csv
```

## Project Structure

```
linkchex/
├── cmd/
│   └── linkchex/
│       └── main.go              # CLI entry point
├── internal/
│   ├── sitemap/
│   │   ├── discover.go          # Sitemap discovery logic
│   │   └── parser.go            # XML parsing logic
│   ├── fetcher/
│   │   ├── client.go            # HTTP client with retries
│   │   └── extractor.go         # HTML link extraction
│   └── validator/
│       ├── validator.go         # Link validation logic
│       └── reporter.go          # Report formatting
├── go.mod
├── PROJECT-PLAN.md
└── README.md
```

## Example Output

### Text Format
```
Link Validation Report
======================

Total Links:    18
✓ Success:      18 (100.0%)
✗ Broken:       0 (0.0%)
⚠ Warnings:     0 (0.0%)
External Links: 7
Duration:       636ms

✓ All links are valid!
```

### When Broken Links Are Found
```
Broken Links:
-------------

✗ https://example.com/nonexistent
  Source: https://example.com/page
  Tag:    <a>
  Status: 404 Not Found
```

## Next Steps: Phase 3

**Phase 3: Advanced Features & Optimization** will implement:
- Performance optimization and tuning
- Progress bars and better UX
- Rate limiting
- Custom user agents
- Link categorization
- Statistics and analytics

See [PROJECT-PLAN.md](PROJECT-PLAN.md) for the complete implementation roadmap.

## Technology Stack

- **Language**: Go 1.25.3
- **Standard Library**: net/http, encoding/xml, encoding/json, encoding/csv
- **Dependencies**:
  - `golang.org/x/net/html` - HTML parsing

## Performance Targets

- 500 pages: 10-30 seconds (target for Phase 3)
- 1000 pages: 20-60 seconds (target for Phase 3)
- Concurrency: 50-200 workers (Phase 3)
- Memory: <500MB for typical workloads

## License

MIT
