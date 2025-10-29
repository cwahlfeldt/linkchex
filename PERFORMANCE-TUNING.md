# Linkchex Performance Tuning Guide

## For Large Sitemaps (500+ pages with many duplicate links)

### Recommended Settings

For your use case (581 URLs, ~80 links/page with many duplicates):

```bash
# Optimal configuration for speed
./linkchex --sitemap your-sitemap.xml \
  --concurrency 200 \
  --timeout 15 \
  --retries 1 \
  --progress

# If checking external links too
./linkchex --sitemap your-sitemap.xml \
  --concurrency 100 \
  --timeout 20 \
  --retries 1 \
  --progress \
  --check-external

# With rate limiting (be nice to servers)
./linkchex --sitemap your-sitemap.xml \
  --concurrency 100 \
  --rate-limit 50 \
  --timeout 15 \
  --retries 1 \
  --progress
```

### Concurrency Settings

| Scenario | Recommended Concurrency | Notes |
|----------|------------------------|-------|
| Internal links only | 200-300 | Your own server can handle it |
| External links included | 50-100 | Be nice to other servers |
| With rate limiting | 50-100 | Rate limit is the bottleneck |
| Slow network | 100-150 | More workers compensate for latency |

### Performance Factors

#### 1. **Concurrency (--concurrency)**
- **Default**: 50 workers
- **Fast**: 200-300 workers (internal links only)
- **Safe**: 100 workers (with external links)
- More workers = more parallel requests
- Diminishing returns after 200-300 workers

#### 2. **Caching (automatic)**
- Already enabled! No configuration needed
- Duplicate URLs validated only once
- Huge benefit for sites with shared navigation/footer links
- With 80% duplicate rate: 46,480 links → ~9,296 unique validations

#### 3. **Timeout (--timeout)**
- **Default**: 10 seconds
- **Fast**: 5-8 seconds (aggressive)
- **Safe**: 15-20 seconds (allows for slow servers)
- Lower timeout = faster failure detection

#### 4. **Retries (--retries)**
- **Default**: 2 retries
- **Fast**: 0-1 retries (fail fast)
- **Safe**: 2-3 retries (more reliable)
- Fewer retries = faster overall completion

#### 5. **Rate Limiting (--rate-limit)**
- **Default**: 0 (unlimited)
- Use when validating external links to avoid overwhelming servers
- 50 req/sec = fast but polite
- Trade-off between speed and server load

## Expected Performance

### Your Sitemap (581 pages, ~80 links/page)

Assumptions:
- 46,480 total link checks
- 80% are duplicates (common for navigation/footer)
- ~9,296 unique URLs to validate

**Without optimization (default settings: 50 workers):**
- Time: ~3-5 minutes

**With optimization (200 workers, internal only):**
- Time: ~30-90 seconds

**With optimization + external links (100 workers):**
- Time: ~1-3 minutes

**With rate limiting (50 req/sec):**
- Time: ~3-4 minutes (limited by rate)

## Quick Performance Test

Test different concurrency levels:

```bash
# Test 1: Default (50 workers)
time ./linkchex --sitemap your-sitemap.xml

# Test 2: High concurrency (200 workers)
time ./linkchex --sitemap your-sitemap.xml --concurrency 200

# Test 3: Very high concurrency (300 workers)
time ./linkchex --sitemap your-sitemap.xml --concurrency 300
```

## Real-World Example

```bash
# Run this on your actual sitemap
./linkchex --sitemap your-sitemap.xml \
  --concurrency 200 \
  --timeout 10 \
  --retries 1 \
  --progress \
  > report.txt

# The progress bar will show you:
# - Real-time validation speed (links/sec)
# - Completion percentage
# - ETA for completion
```

## Caching Statistics

After running, check the report for:
```
Total Links:       46,480   (all link checks)
Unique URLs:       9,296    (actually validated)
Cached Results:    9,296    (saved ~37,184 HTTP requests!)
```

The cache saves you from making ~80% of HTTP requests!

## Bottlenecks

1. **Network Latency**: Most common bottleneck
   - Solution: Increase concurrency (200-300)

2. **Target Server Rate Limits**: Server blocking/throttling you
   - Solution: Use --rate-limit or reduce concurrency

3. **DNS Resolution**: Many unique domains
   - Solution: System-level DNS caching (already happens)

4. **Memory**: Too many concurrent connections
   - Solution: Reduce concurrency if you see memory issues

## Tips for Maximum Speed

1. **Internal links only** (skip external): Add no flags, they're excluded by default
2. **High concurrency**: `--concurrency 200` or higher
3. **Fast fail**: `--retries 0` or `--retries 1`
4. **Aggressive timeout**: `--timeout 5`
5. **Use progress bar**: `--progress` to monitor speed

## If Things Go Wrong

**Too many open files error:**
```bash
# macOS/Linux: Increase file descriptor limit
ulimit -n 4096
```

**Memory issues:**
```bash
# Reduce concurrency
./linkchex --sitemap your-sitemap.xml --concurrency 50
```

**Rate limited by server:**
```bash
# Add rate limiting
./linkchex --sitemap your-sitemap.xml --rate-limit 20
```

## Monitoring

Watch the progress bar for these metrics:
- **it/s (iterations per second)**: How fast links are being validated
- **30-60 it/s**: Good speed for external links
- **100-300 it/s**: Great speed for internal links
- **<10 it/s**: Bottleneck detected (slow servers or network issues)
