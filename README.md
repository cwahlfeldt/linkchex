# Linkchex

A high-performance CLI tool written in Go that discovers sitemaps, crawls pages, and validates all links with maximum speed and efficiency.

## Current Status: Phase 3 Complete ✓

**Phase 1: Project Setup & Sitemap Discovery** - ✓ Complete
**Phase 2: Basic Link Validation** - ✓ Complete
**Phase 3: Advanced Features & Optimization** - ✓ Complete

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

#### Phase 3
- ✓ Progress bar with real-time statistics
- ✓ Rate limiting to avoid overwhelming servers
- ✓ URL exclusion patterns (glob-style wildcards)
- ✓ URL caching for duplicate link detection
- ✓ Enhanced statistics (unique URLs, cache hits, links by type)
- ✓ Performance optimizations
- ✓ Link categorization by tag type and status code

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

# Show progress bar (Phase 3)
./linkchex --sitemap test-sitemap.xml --progress

# Rate limit to 5 requests per second (Phase 3)
./linkchex --sitemap test-sitemap.xml --rate-limit 5

# Exclude URLs matching pattern (Phase 3)
./linkchex --sitemap test-sitemap.xml --exclude "*/admin/*"
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
│   │   ├── client.go            # HTTP client with retries & rate limiting
│   │   ├── extractor.go         # HTML link extraction
│   │   └── ratelimiter.go       # Rate limiting implementation
│   └── validator/
│       ├── validator.go         # Link validation logic
│       ├── reporter.go          # Report formatting
│       └── patterns.go          # URL pattern matching
├── go.mod
├── PROJECT-PLAN.md
└── README.md
```

## Example Output

### Text Format (Phase 3 Enhanced)
```
Link Validation Report
======================

Pages Processed:   1
Total Links:       18
Unique URLs:       18
✓ Success:         18 (100.0%)
✗ Broken:          0 (0.0%)
⚠ Warnings:        0 (0.0%)
Internal Links:    11
External Links:    7
Cached Results:    18
Duration:          622ms

Links by Type:
  <a>: 17
  <img>: 1

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

## Technology Stack

- **Language**: Go 1.25.3
- **Standard Library**: net/http, encoding/xml, encoding/json, encoding/csv
- **Dependencies**:
  - `golang.org/x/net/html` - HTML parsing
  - `github.com/schollz/progressbar/v3` - Progress bar display

## Performance Features

- **Concurrency**: 50 workers by default (configurable)
- **Caching**: Duplicate URL detection and caching
- **Rate Limiting**: Configurable requests per second
- **Memory**: Efficient memory usage with streaming and caching
- **Speed**: Typical validation speed: 30-60 links/second (depends on network)

## Features by Phase

### Phase 1: Foundation
- Sitemap discovery and parsing
- Basic CLI framework

### Phase 2: Core Validation
- HTTP client with retries
- Link extraction and validation
- Multiple output formats
- Concurrent processing

### Phase 3: Advanced Features
- Progress bar with real-time stats
- Rate limiting
- URL exclusion patterns
- Enhanced statistics
- Performance optimizations

## Future Enhancements

Potential features for future development:
- Recursive crawling (follow links beyond sitemap)
- Database storage for results
- Web dashboard
- Scheduled validation runs
- Email notifications
- Custom report templates
- Browser-based validation for JavaScript-heavy sites

See [PROJECT-PLAN.md](PROJECT-PLAN.md) for the complete implementation roadmap.

## License

MIT
