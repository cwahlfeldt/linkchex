#!/bin/bash

# Benchmark script for Linkchex
# Tests different concurrency levels to find optimal setting

SITEMAP="$1"

if [ -z "$SITEMAP" ]; then
    echo "Usage: ./benchmark.sh <sitemap-file>"
    echo "Example: ./benchmark.sh your-sitemap.xml"
    exit 1
fi

echo "Benchmarking Linkchex Performance"
echo "=================================="
echo "Sitemap: $SITEMAP"
echo ""

# Test 1: Default settings (50 workers)
echo "Test 1: Default (50 workers)"
time ./linkchex --sitemap "$SITEMAP" > /dev/null 2>&1
echo ""

# Test 2: Medium concurrency (100 workers)
echo "Test 2: Medium (100 workers)"
time ./linkchex --sitemap "$SITEMAP" --concurrency 100 > /dev/null 2>&1
echo ""

# Test 3: High concurrency (200 workers)
echo "Test 3: High (200 workers)"
time ./linkchex --sitemap "$SITEMAP" --concurrency 200 > /dev/null 2>&1
echo ""

# Test 4: Very high concurrency (300 workers)
echo "Test 4: Very High (300 workers)"
time ./linkchex --sitemap "$SITEMAP" --concurrency 300 > /dev/null 2>&1
echo ""

echo "Benchmark complete! Check which concurrency level was fastest."
echo ""
echo "Recommended: Use the fastest setting with --progress for real runs:"
echo "./linkchex --sitemap $SITEMAP --concurrency <fastest> --progress"
