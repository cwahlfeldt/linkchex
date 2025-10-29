# Sitemap Link Checker - Implementation Phases aka "Linkchex"

## Project Overview
A high-performance CLI tool written in Go that discovers sitemaps, crawls pages, and validates all links with maximum speed and efficiency.

**Target Scale**: 500-1000 URLs per sitemap  
**Technology**: Go with native concurrency  
**Output**: Single compiled binary

---

## Phase 1: Project Setup & Sitemap Discovery
**Impact**: Small  
**Goal**: Bootstrap project and implement sitemap discovery functionality

### Tasks
1. **Project Initialization**
   - Initialize Go module
   - Set up project structure (cmd/, internal/ directories)
   - Create main.go entry point
   - Set up basic CLI framework (cobra optional, or use flag package)

2. **Sitemap Discovery**
   - Implement automatic sitemap discovery from base URL
   - Check common sitemap locations:
     - `/sitemap.xml`
     - `/sitemap_index.xml`
     - `/sitemap/`
     - Parse `/robots.txt` for Sitemap: directives
   - Fallback to user-provided sitemap URL if discovery fails
   - Handle both sitemap files and sitemap indexes

3. **Sitemap Parser**
   - Parse XML sitemap format
   - Extract all `<loc>` URLs
   - Handle sitemap index files (nested sitemaps)
   - Support local sitemap files and remote URLs
   - Validate sitemap structure

### Deliverables
- Basic CLI that accepts `--url` or `--sitemap` flag
- Sitemap discovery from base URL
- Parsed list of URLs from sitemap
- Simple text output showing discovered URLs

### Testing
- Test with various sitemap formats
- Test robots.txt parsing
- Test sitemap index files
- Test error handling for missing sitemaps

---

## Phase 2: Basic Link Validation
**Impact**: Medium  
**Goal**: Implement core link checking functionality

### Tasks
1. **HTTP Client Setup**
   - Configure HTTP client with timeouts
   - Set up user agent string
   - Implement request retry logic
   - Handle redirects (follow up to 10 hops)

2. **Page Fetcher**
   - Fetch individual pages from sitemap
   - Implement HEAD request strategy (fallback to GET)
   - Basic error handling for network failures
   - Status code validation (200-299 = success)

3. **Link Extractor**
   - Parse HTML using goquery
   - Extract links from:
     - `<a href>` tags
     - `<img src>` tags
     - `<link href>` tags (CSS, favicons)
     - `<script src>` tags
   - URL normalization (resolve relative paths)
   - Remove duplicates and fragments

4. **Link Validator**
   - Check each extracted link
   - Validate status codes
   - Categorize results (success, broken, timeout)
   - Track response times

### Deliverables
- Single-threaded link checker (no concurrency yet)
- Text output showing broken links
- Basic error messages for failed checks

### Testing
- Test with pages containing various link types
- Test relative vs absolute URLs
- Test different status codes (200, 404, 500, etc.)
- Test redirect following

---

## Phase 3: Concurrency & Performance
**Impact**: Large  
**Goal**: Add concurrent processing for maximum speed

### Tasks
1. **Worker Pool Implementation**
   - Create configurable worker pool (default 50 workers)
   - Channel-based job distribution
   - Use goroutines for concurrent page fetching
   - Implement semaphore pattern for rate limiting

2. **Concurrent Link Validation**
   - Parallel link checking across workers
   - Thread-safe URL deduplication cache
   - Use sync.Map or mutex-protected map
   - Prevent duplicate checks across pages

3. **Connection Pooling**
   - Configure HTTP transport for connection reuse
   - Set MaxIdleConns and MaxIdleConnsPerHost
   - Enable keep-alive connections
   - Optimize for high-throughput scenarios

4. **Memory Management**
   - Stream processing where possible
   - Efficient data structures for URL storage
   - Garbage collection considerations
   - Monitor memory usage during large runs

### Deliverables
- Fully concurrent link checker
- Configurable concurrency via `--concurrency` flag
- Significant performance improvement (10-30 seconds for 500 pages)
- Thread-safe operations

### Testing
- Benchmark against single-threaded version
- Test with various concurrency levels (10, 50, 100, 200)
- Memory profiling
- Race condition testing (`go test -race`)

---

## Phase 4: Progress Reporting & Output Formats
**Impact**: Medium  
**Goal**: Add user-friendly progress tracking and multiple output formats

### Tasks
1. **Progress Tracking**
   - Real-time progress bar (if verbose mode)
   - Statistics tracking:
     - Pages crawled / total
     - Links validated / total
     - Broken links found
     - Average response time
   - Time elapsed and ETA
   - Current page being processed

2. **Output Formatters**
   - **Text format** (default):
     - Colored output (green/red for pass/fail)
     - Grouped by page
     - Summary statistics
   - **JSON format**:
     - Structured data for programmatic use
     - Include all metadata (status codes, response times)
   - **CSV format**:
     - Spreadsheet-compatible
     - Columns: source_page, link_url, status_code, error
   - **HTML report** (optional):
     - Styled report with filtering
     - Interactive dashboard

3. **Logging & Verbosity**
   - Quiet mode (errors only)
   - Normal mode (summary + broken links)
   - Verbose mode (all details + progress)
   - Debug mode (HTTP requests/responses)

### Deliverables
- Multiple output formats via `--format` flag
- Progress bar in verbose mode
- File output via `--output` flag
- Summary statistics

### Testing
- Verify output format correctness
- Test file writing permissions
- Validate JSON/CSV structure
- Test progress bar rendering

---

## Phase 5: Advanced Features & Refinement
**Impact**: Medium  
**Goal**: Add configuration options and edge case handling

### Tasks
1. **Configuration Options**
   - `--timeout`: Request timeout (default: 10s)
   - `--max-retries`: Retry attempts (default: 2)
   - `--user-agent`: Custom user agent
   - `--external`: Check external links (default: internal only)
   - `--ignore-patterns`: Regex patterns to skip
   - `--follow-redirects`: Enable/disable redirect following
   - `--max-redirects`: Maximum redirect hops

2. **Link Filtering**
   - Internal vs external link detection
   - Ignore specific domains or patterns
   - Option to check only internal links
   - Option to check only broken links
   - Skip specific status codes (e.g., ignore 403s)

3. **Error Handling & Recovery**
   - Graceful handling of malformed HTML
   - Handle SSL/TLS errors
   - Timeout handling for slow sites
   - DNS resolution failures
   - Connection refused errors
   - Resume capability (checkpoint progress)

4. **Robots.txt Respect**
   - Optional robots.txt checking
   - Respect crawl-delay directives
   - Skip disallowed paths
   - Configurable user-agent for robots.txt

### Deliverables
- Comprehensive CLI flags
- Robust error handling
- Configuration file support (optional)
- Edge case coverage

### Testing
- Test with malformed HTML
- Test with slow/timeout scenarios
- Test with various network errors
- Test robots.txt compliance
- Integration tests with real websites

---

## Phase 6: Optimization & Polish
**Impact**: Small  
**Goal**: Final optimizations and production readiness

### Tasks
1. **Performance Tuning**
   - Profile CPU and memory usage
   - Optimize hot paths
   - Reduce allocations
   - Benchmark improvements

2. **Documentation**
   - Comprehensive README with examples
   - CLI help text and usage examples
   - Installation instructions
   - Performance tuning guide
   - Troubleshooting section

3. **Distribution**
   - Build scripts for multiple platforms
   - Cross-compilation (Linux, macOS, Windows)
   - Release binaries
   - Version management
   - Optional: Docker container

4. **Testing & Quality**
   - Increase test coverage (aim for >80%)
   - Add integration tests
   - CI/CD pipeline setup
   - Static analysis (golint, go vet)
   - Security scanning

### Deliverables
- Production-ready binary
- Complete documentation
- Multi-platform releases
- Comprehensive test suite

### Testing
- End-to-end testing
- Performance benchmarks
- Cross-platform testing
- Security audit

---

## CLI Usage Examples

### Basic Usage
```bash
# Discover and check sitemap from base URL
./sitemap-checker --url https://example.com

# Use specific sitemap
./sitemap-checker --sitemap https://example.com/sitemap.xml

# Local sitemap file
./sitemap-checker --sitemap ./sitemap.xml
```

### Advanced Usage
```bash
# High concurrency with verbose output
./sitemap-checker --url https://example.com --concurrency 100 --verbose

# Check external links with custom timeout
./sitemap-checker --url https://example.com --external --timeout 30

# Export to JSON file
./sitemap-checker --url https://example.com --format json --output results.json

# Ignore specific patterns
./sitemap-checker --url https://example.com --ignore-patterns ".*\.pdf$,.*linkedin\.com.*"

# Custom user agent with retries
./sitemap-checker --url https://example.com --user-agent "MyBot/1.0" --max-retries 3
```

---

## Success Metrics

### Performance Targets
- **500 pages**: Complete in 10-30 seconds
- **1000 pages**: Complete in 20-60 seconds
- **Concurrency**: Handle 50-200 concurrent connections
- **Memory**: Stay under 500MB for typical workloads

### Quality Targets
- **Accuracy**: 100% accurate link validation
- **Reliability**: Handle network failures gracefully
- **Usability**: Clear, actionable output
- **Maintainability**: Well-structured, documented code

---

## Dependencies

### Required
- Go 1.21+ (for standard library features)
- `golang.org/x/net/html` (HTML parsing)
- `github.com/PuerkitoBio/goquery` (jQuery-like selectors)

### Optional
- `github.com/spf13/cobra` (CLI framework)
- `github.com/spf13/viper` (configuration)
- `github.com/schollz/progressbar/v3` (progress bars)
- `github.com/fatih/color` (colored output)

---

## Implementation Order

1. **Start Here**: Phase 1 (Sitemap Discovery)
2. **Core Functionality**: Phase 2 (Basic Link Validation)
3. **Performance**: Phase 3 (Concurrency)
4. **User Experience**: Phase 4 (Progress & Output)
5. **Robustness**: Phase 5 (Advanced Features)
6. **Polish**: Phase 6 (Optimization)

**Estimated Timeline**: 2-3 weeks for full implementation (can complete core functionality in Phase 1-3 within 1 week)
