package main

import (
	"flag"
	"fmt"
	"os"

	"linkchex/internal/sitemap"
	"linkchex/internal/validator"
)

const version = "0.1.1"

func main() {
	// Define CLI flags
	url := flag.String("url", "", "Base URL to discover sitemap from")
	sitemapURL := flag.String("sitemap", "", "Direct URL or path to sitemap file")
	concurrency := flag.Int("concurrency", 200, "Number of concurrent workers")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	versionFlag := flag.Bool("version", false, "Show version information")
	timeout := flag.Int("timeout", 10, "Request timeout in seconds")
	format := flag.String("format", "text", "Output format (text, json, csv)")
	output := flag.String("output", "", "Output file path (default: stdout)")
	maxRetries := flag.Int("retries", 1, "Maximum number of retries for failed requests")
	checkExternal := flag.Bool("check-external", false, "Check external links (default: internal only)")
	listOnly := flag.Bool("list-only", false, "Only list URLs from sitemap without validating links")
	rateLimit := flag.Float64("rate-limit", 0, "Rate limit in requests per second (0 = unlimited)")
	excludePattern := flag.String("exclude", "", "Exclude URLs matching pattern (supports * and ? wildcards)")
	showProgress := flag.Bool("progress", false, "Show progress bar (auto-disabled with --verbose)")

	flag.Parse()

	// Show version
	if *versionFlag {
		fmt.Printf("linkchex version %s\n", version)
		os.Exit(0)
	}

	// Validate input
	if *url == "" && *sitemapURL == "" {
		fmt.Fprintln(os.Stderr, "Error: Either --url or --sitemap must be provided")
		flag.Usage()
		os.Exit(1)
	}

	if *url != "" && *sitemapURL != "" {
		fmt.Fprintln(os.Stderr, "Error: Cannot specify both --url and --sitemap")
		flag.Usage()
		os.Exit(1)
	}

	// Configuration
	config := &Config{
		URL:            *url,
		SitemapURL:     *sitemapURL,
		Concurrency:    *concurrency,
		Verbose:        *verbose,
		Timeout:        *timeout,
		Format:         *format,
		Output:         *output,
		MaxRetries:     *maxRetries,
		CheckExternal:  *checkExternal,
		ListOnly:       *listOnly,
		RateLimit:      *rateLimit,
		ExcludePattern: *excludePattern,
		ShowProgress:   *showProgress,
	}

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type Config struct {
	URL            string
	SitemapURL     string
	Concurrency    int
	Verbose        bool
	Timeout        int
	Format         string
	Output         string
	MaxRetries     int
	CheckExternal  bool
	ListOnly       bool
	RateLimit      float64
	ExcludePattern string
	ShowProgress   bool
}

func run(config *Config) error {
	if config.Verbose {
		fmt.Println("Starting linkchex...")
		fmt.Printf("Configuration: %+v\n\n", config)
	}

	// Discover or use provided sitemap
	var sitemapURLs []string
	var err error

	if config.URL != "" {
		if config.Verbose {
			fmt.Printf("Discovering sitemap from base URL: %s\n", config.URL)
		}
		sitemapURLs, err = sitemap.Discover(config.URL)
		if err != nil {
			return fmt.Errorf("sitemap discovery failed: %w", err)
		}
	} else {
		if config.Verbose {
			fmt.Printf("Using provided sitemap: %s\n", config.SitemapURL)
		}
		sitemapURLs = []string{config.SitemapURL}
	}

	if config.Verbose {
		fmt.Printf("Found %d sitemap(s)\n", len(sitemapURLs))
		for i, url := range sitemapURLs {
			fmt.Printf("  %d. %s\n", i+1, url)
		}
		fmt.Println()
	}

	// Parse sitemaps and extract URLs
	var allURLs []string
	for _, sitemapURL := range sitemapURLs {
		if config.Verbose {
			fmt.Printf("Parsing sitemap: %s\n", sitemapURL)
		}
		urls, err := sitemap.Parse(sitemapURL)
		if err != nil {
			return fmt.Errorf("failed to parse sitemap %s: %w", sitemapURL, err)
		}
		allURLs = append(allURLs, urls...)
	}

	if config.Verbose {
		fmt.Printf("\nDiscovered %d URLs from sitemap(s)\n\n", len(allURLs))
	}

	// If list-only mode, just display URLs and exit
	if config.ListOnly {
		if config.Format == "text" {
			fmt.Println("URLs discovered:")
			fmt.Println("================")
			for i, url := range allURLs {
				fmt.Printf("%d. %s\n", i+1, url)
			}
			fmt.Printf("\nTotal: %d URLs\n", len(allURLs))
		}
		return nil
	}

	// Validate links on all pages
	if config.Verbose {
		fmt.Println("Starting link validation...")
	}

	v := validator.NewValidator(config.Timeout, config.MaxRetries, config.Concurrency, config.Verbose)

	// Set rate limiting if specified
	if config.RateLimit > 0 {
		if config.Verbose {
			fmt.Printf("Rate limiting enabled: %.2f requests/second\n", config.RateLimit)
		}
		v.SetRateLimit(config.RateLimit)
	}

	// Set exclude patterns if specified
	if config.ExcludePattern != "" {
		patterns := []string{config.ExcludePattern}
		if err := v.SetExcludePatterns(patterns); err != nil {
			return fmt.Errorf("invalid exclude pattern: %w", err)
		}
		if config.Verbose {
			fmt.Printf("Excluding URLs matching: %s\n", config.ExcludePattern)
		}
	}

	// Set progress bar visibility
	if config.ShowProgress && !config.Verbose {
		v.SetShowProgress(true)
	}

	report := v.ValidateMultiplePages(allURLs, config.CheckExternal)

	// Format and output report
	reportText, err := validator.FormatReport(report, config.Format)
	if err != nil {
		return fmt.Errorf("failed to format report: %w", err)
	}

	if config.Output != "" {
		// Write to file
		if err := validator.WriteReportToFile(report, config.Format, config.Output); err != nil {
			return fmt.Errorf("failed to write report to file: %w", err)
		}
		if config.Verbose {
			fmt.Printf("\nReport written to: %s\n", config.Output)
		}
	} else {
		// Write to stdout
		fmt.Println(reportText)
	}

	// Exit with error code if broken links found
	if report.BrokenLinks > 0 {
		os.Exit(1)
	}

	return nil
}
