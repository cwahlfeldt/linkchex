# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Linkchex is a high-performance CLI tool written in Go that discovers sitemaps, crawls pages, and validates all links with maximum speed and efficiency. It's optimized for large sitemaps (500+ pages) with heavy caching to avoid duplicate URL validation.

## Build and Run Commands

```bash
# Build the binary
go build -o linkchex ./cmd/linkchex

# Run with local sitemap (internal links only)
./linkchex --sitemap test-sitemap.xml

# Run with high concurrency for large sitemaps
./linkchex --sitemap test-sitemap.xml --concurrency 200 --progress

# Run with external link checking
./linkchex --sitemap test-sitemap.xml --check-external

# Skip <link> and <script> tag validation (faster, focuses on content links)
./linkchex --sitemap test-sitemap.xml --skip-resources

# Run performance benchmark
./benchmark.sh test-sitemap.xml

# Install dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Architecture

### Core Components

**Entry Point**: [cmd/linkchex/main.go](cmd/linkchex/main.go)
- CLI flag parsing and configuration
- Orchestrates the sitemap → validation → reporting workflow
- Default concurrency: 400 workers, timeout: 10s, retries: 1

**Sitemap Discovery**: [internal/sitemap/discover.go](internal/sitemap/discover.go)
- Discovers sitemaps from base URLs via robots.txt parsing
- Falls back to common sitemap locations (/sitemap.xml, etc.)
- Handles both local files and remote URLs

**XML Parsing**: [internal/sitemap/parser.go](internal/sitemap/parser.go)
- Parses sitemap XML and sitemap index files
- Recursively handles nested sitemaps
- Supports both local and remote sitemap files

**HTTP Client**: [internal/fetcher/client.go](internal/fetcher/client.go)
- Configurable timeout, retries, and rate limiting
- Uses HEAD requests for link validation (efficiency)
- Uses GET requests for page fetching (extracts HTML)
- Custom user agent: "Linkchex/0.1.0 (Link Validator)"

**Link Extraction**: [internal/fetcher/extractor.go](internal/fetcher/extractor.go)
- Extracts links from HTML: `<a>`, `<img>`, `<link>`, `<script>` tags
- Resolves relative URLs to absolute URLs
- Classifies links as internal/external based on domain

**Link Validation**: [internal/validator/validator.go](internal/validator/validator.go)
- Validates links with configurable concurrency (semaphore pattern)
- **Critical feature**: URL caching to avoid duplicate validation
- Progress bar with real-time statistics (using progressbar/v3)
- Rate limiting to avoid overwhelming servers
- URL pattern exclusion (glob-style wildcards)

**Reporting**: [internal/validator/reporter.go](internal/validator/reporter.go)
- Formats output as text, JSON, or CSV
- Tracks statistics: success/broken/warning links, cache hits, unique URLs
- Groups links by tag type (`<a>`, `<img>`, etc.) and status code

### Data Flow

1. **Discovery**: Base URL → robots.txt → sitemap URLs
2. **Parsing**: Sitemap XML → list of page URLs
3. **Fetching**: Page URL → HTML content
4. **Extraction**: HTML → list of links (with metadata)
5. **Validation**: Links → HTTP HEAD requests (with caching)
6. **Reporting**: Results → formatted output

### Concurrency Model

- **Worker pool**: Semaphore with configurable concurrency (default: 400)
- **Caching**: Thread-safe map with RWMutex for duplicate URL detection
- **Rate limiting**: Token bucket algorithm in [internal/fetcher/ratelimiter.go](internal/fetcher/ratelimiter.go)
- **Parallel processing**: Validates all pages' links concurrently

### Performance Optimizations

1. **URL Caching**: Each unique URL validated only once, stored in memory
   - Critical for large sitemaps with shared navigation/footer links
   - Typical cache hit rate: 80%+ on real-world sites

2. **HEAD Requests**: Uses HEAD instead of GET for validation (faster)

3. **High Concurrency**: Default 400 workers for internal links (configurable via --concurrency)

4. **Fast Fail**: Default 1 retry (configurable via --retries)

5. **Progress Bar**: Shows real-time speed (links/sec) and ETA

See [PERFORMANCE-TUNING.md](PERFORMANCE-TUNING.md) for detailed tuning guidance.

## Key Design Decisions

### Why caching is critical
Large sites often have 80%+ duplicate links across pages (navigation, footer, common content). Without caching, a 500-page site with 80 links each = 40,000 HTTP requests. With caching and 80% duplicates, only ~8,000 unique validations needed.

### Why high default concurrency (400)
Internal links are to your own server, which can handle high concurrency. External link checking should use lower concurrency (50-100) to be polite to other servers.

### Why HEAD requests
HEAD requests are lighter than GET requests since they only fetch headers, not the full response body. This speeds up validation significantly.

### Exit codes
- `0` = Success (all links valid)
- `1` = Failure (broken links found or error occurred)

## Testing

Test files included:
- [test-sitemap.xml](test-sitemap.xml) - Basic 4-URL sitemap
- [test-sitemap-index.xml](test-sitemap-index.xml) - Sitemap index (nested sitemaps)
- [test-google-sitemap.xml](test-google-sitemap.xml) - Real-world Google homepage
- [test-duplicate-links.xml](test-duplicate-links.xml) - Tests caching behavior
- [test-many-duplicates.xml](test-many-duplicates.xml) - Heavy duplication scenario

## Common Development Patterns

### Adding new CLI flags
1. Define flag in [cmd/linkchex/main.go](cmd/linkchex/main.go) `main()` function
2. Add to `Config` struct
3. Pass through to relevant component (validator, fetcher, etc.)

### Adding new output formats
1. Add format type to [internal/validator/reporter.go](internal/validator/reporter.go)
2. Implement formatting logic in `FormatReport()` function
3. Update `WriteReportToFile()` if needed

### Modifying validation logic
Primary location: [internal/validator/validator.go](internal/validator/validator.go)
- `validateLink()`: Single link validation (with caching)
- `validateLinks()`: Concurrent validation with progress bar
- `ValidateMultiplePages()`: Orchestrates multi-page validation

### Adding new link types
1. Update HTML parser in [internal/fetcher/extractor.go](internal/fetcher/extractor.go)
2. Add tag type to extraction logic
3. Statistics automatically tracked by tag type

## Dependencies

- `golang.org/x/net/html` - HTML parsing
- `github.com/schollz/progressbar/v3` - Progress bar UI
- Standard library: net/http, encoding/xml, encoding/json, encoding/csv

## Performance Characteristics

For a typical large sitemap (581 pages, ~80 links/page):
- **Total link checks**: 46,480
- **With 80% duplicates**: ~9,300 unique validations
- **Default settings (50 workers)**: 3-5 minutes
- **Optimized (200 workers)**: 30-90 seconds
- **Validation speed**: 30-60 links/sec (network dependent)
