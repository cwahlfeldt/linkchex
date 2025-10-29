package validator

import (
	"fmt"
	"sync"
	"time"

	"linkchex/internal/fetcher"
	"github.com/schollz/progressbar/v3"
)

// Result represents the validation result for a single URL
type Result struct {
	SourceURL  string        // The page where the link was found
	TargetURL  string        // The link being validated
	StatusCode int           // HTTP status code
	Status     string        // Status text
	Error      error         // Error if validation failed
	IsExternal bool          // Whether the link is external
	Tag        string        // HTML tag (a, img, link, script)
	LinkText   string        // Text content of the link (for <a> tags)
	Duration   time.Duration // Time taken to validate
	IsBroken   bool          // Whether the link is broken
}

// ValidationReport contains all validation results
type ValidationReport struct {
	Results        []Result
	TotalLinks     int
	BrokenLinks    int
	WarningLinks   int
	SuccessLinks   int
	ExternalLinks  int
	InternalLinks  int
	CachedLinks    int
	UniqueURLs     int
	PagesProcessed int
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	LinksByTag     map[string]int // Count of links by tag type
	LinksByStatus  map[int]int    // Count of links by status code
}

// Validator validates links from pages
type Validator struct {
	client       *fetcher.Client
	concurrency  int
	verbose      bool
	showProgress bool
	urlCache     map[string]*Result
	cacheMutex   sync.RWMutex
	urlMatcher   *URLMatcher
}

// NewValidator creates a new link validator
func NewValidator(timeout, maxRetries, concurrency int, verbose bool) *Validator {
	return &Validator{
		client:       fetcher.NewClient(timeout, maxRetries),
		concurrency:  concurrency,
		verbose:      verbose,
		showProgress: !verbose, // Show progress bar only when not verbose
		urlCache:     make(map[string]*Result),
	}
}

// SetShowProgress controls whether to show progress bar
func (v *Validator) SetShowProgress(show bool) {
	v.showProgress = show
}

// SetExcludePatterns sets URL patterns to exclude from validation
func (v *Validator) SetExcludePatterns(patterns []string) error {
	matcher, err := NewURLMatcher(patterns, nil)
	if err != nil {
		return err
	}
	v.urlMatcher = matcher
	return nil
}

// SetURLMatcher sets a custom URL matcher
func (v *Validator) SetURLMatcher(matcher *URLMatcher) {
	v.urlMatcher = matcher
}

// SetRateLimit sets the rate limit for HTTP requests (requests per second)
func (v *Validator) SetRateLimit(requestsPerSecond float64) {
	v.client.SetRateLimit(requestsPerSecond)
}

// ValidatePage fetches a page and validates all links on it
func (v *Validator) ValidatePage(pageURL string, checkExternal bool) ([]Result, error) {
	if v.verbose {
		fmt.Printf("Fetching page: %s\n", pageURL)
	}

	// Fetch the page
	resp := v.client.Get(pageURL)
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", resp.Error)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("page returned status %d", resp.StatusCode)
	}

	// Extract links
	links, err := fetcher.ExtractLinks(resp.Body, pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract links: %w", err)
	}

	// Filter links
	links = fetcher.FilterLinks(links, checkExternal)

	if v.verbose {
		fmt.Printf("Found %d links to validate\n", len(links))
	}

	// Validate links concurrently
	return v.validateLinks(pageURL, links), nil
}

// validateLinks validates multiple links concurrently
func (v *Validator) validateLinks(sourceURL string, links []fetcher.Link) []Result {
	results := make([]Result, len(links))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, v.concurrency)

	// Create progress bar if enabled
	var bar *progressbar.ProgressBar
	if v.showProgress && len(links) > 0 {
		bar = progressbar.NewOptions(len(links),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetDescription("[cyan]Validating links...[reset]"),
			progressbar.OptionSetWidth(50),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
	}

	for i, link := range links {
		wg.Add(1)
		go func(idx int, l fetcher.Link) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			results[idx] = v.validateLink(sourceURL, l)

			if bar != nil {
				bar.Add(1)
			}
		}(i, link)
	}

	wg.Wait()

	if bar != nil {
		bar.Finish()
		fmt.Println() // Add newline after progress bar
	}

	return results
}

// validateLink validates a single link
func (v *Validator) validateLink(sourceURL string, link fetcher.Link) Result {
	// Check if URL should be validated
	if v.urlMatcher != nil && !v.urlMatcher.ShouldCheck(link.URL) {
		return Result{
			SourceURL:  sourceURL,
			TargetURL:  link.URL,
			StatusCode: 0,
			Status:     "Skipped (excluded by pattern)",
			Error:      nil,
			IsExternal: link.IsExternal,
			Tag:        link.Tag,
			LinkText:   link.Text,
			Duration:   0,
			IsBroken:   false,
		}
	}

	// Check cache first
	v.cacheMutex.RLock()
	if cached, found := v.urlCache[link.URL]; found {
		v.cacheMutex.RUnlock()
		// Return cached result with updated source
		cachedCopy := *cached
		cachedCopy.SourceURL = sourceURL
		cachedCopy.Tag = link.Tag
		cachedCopy.LinkText = link.Text
		return cachedCopy
	}
	v.cacheMutex.RUnlock()

	// Use HEAD request for efficiency
	resp := v.client.Head(link.URL)

	result := Result{
		SourceURL:  sourceURL,
		TargetURL:  link.URL,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Error:      resp.Error,
		IsExternal: link.IsExternal,
		Tag:        link.Tag,
		LinkText:   link.Text,
		Duration:   resp.Duration,
	}

	// Determine if link is broken
	if resp.Error != nil {
		result.IsBroken = true
	} else if resp.StatusCode >= 400 {
		result.IsBroken = true
	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		// Redirects are warnings, not broken
		result.IsBroken = false
	}

	// Cache the result
	v.cacheMutex.Lock()
	v.urlCache[link.URL] = &result
	v.cacheMutex.Unlock()

	return result
}

// ValidateMultiplePages validates links from multiple pages
func (v *Validator) ValidateMultiplePages(pageURLs []string, checkExternal bool) *ValidationReport {
	report := &ValidationReport{
		Results:      make([]Result, 0),
		StartTime:    time.Now(),
		LinksByTag:   make(map[string]int),
		LinksByStatus: make(map[int]int),
	}

	if v.verbose {
		fmt.Printf("\nValidating %d pages...\n\n", len(pageURLs))
	}

	report.PagesProcessed = len(pageURLs)
	uniqueURLs := make(map[string]bool)

	for i, pageURL := range pageURLs {
		if v.verbose {
			fmt.Printf("[%d/%d] Validating: %s\n", i+1, len(pageURLs), pageURL)
		}

		results, err := v.ValidatePage(pageURL, checkExternal)
		if err != nil {
			if v.verbose {
				fmt.Printf("  ⚠ Error validating page: %v\n", err)
			}
			// Create a result for the page itself
			report.Results = append(report.Results, Result{
				SourceURL:  "sitemap",
				TargetURL:  pageURL,
				StatusCode: 0,
				Status:     "Failed",
				Error:      err,
				IsBroken:   true,
			})
			continue
		}

		report.Results = append(report.Results, results...)

		if v.verbose {
			fmt.Printf("  Found %d links\n", len(results))
		}
	}

	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime)

	// Calculate statistics
	for _, result := range report.Results {
		report.TotalLinks++

		// Track unique URLs
		uniqueURLs[result.TargetURL] = true

		// Categorize by status
		if result.IsBroken {
			report.BrokenLinks++
		} else if result.StatusCode >= 300 && result.StatusCode < 400 {
			report.WarningLinks++
		} else {
			report.SuccessLinks++
		}

		// Count internal vs external
		if result.IsExternal {
			report.ExternalLinks++
		} else {
			report.InternalLinks++
		}

		// Count by tag type
		if result.Tag != "" {
			report.LinksByTag[result.Tag]++
		}

		// Count by status code
		if result.StatusCode > 0 {
			report.LinksByStatus[result.StatusCode]++
		}
	}

	// Track unique URLs count
	report.UniqueURLs = len(uniqueURLs)

	// Track cached links
	v.cacheMutex.RLock()
	report.CachedLinks = len(v.urlCache)
	v.cacheMutex.RUnlock()

	return report
}
