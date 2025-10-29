package validator

import (
	"fmt"
	"sync"
	"time"

	"linkchex/internal/fetcher"
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
	Results       []Result
	TotalLinks    int
	BrokenLinks   int
	WarningLinks  int
	SuccessLinks  int
	ExternalLinks int
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
}

// Validator validates links from pages
type Validator struct {
	client      *fetcher.Client
	concurrency int
	verbose     bool
}

// NewValidator creates a new link validator
func NewValidator(timeout, maxRetries, concurrency int, verbose bool) *Validator {
	return &Validator{
		client:      fetcher.NewClient(timeout, maxRetries),
		concurrency: concurrency,
		verbose:     verbose,
	}
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

	for i, link := range links {
		wg.Add(1)
		go func(idx int, l fetcher.Link) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			results[idx] = v.validateLink(sourceURL, l)
		}(i, link)
	}

	wg.Wait()
	return results
}

// validateLink validates a single link
func (v *Validator) validateLink(sourceURL string, link fetcher.Link) Result {
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

	return result
}

// ValidateMultiplePages validates links from multiple pages
func (v *Validator) ValidateMultiplePages(pageURLs []string, checkExternal bool) *ValidationReport {
	report := &ValidationReport{
		Results:   make([]Result, 0),
		StartTime: time.Now(),
	}

	if v.verbose {
		fmt.Printf("\nValidating %d pages...\n\n", len(pageURLs))
	}

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
		if result.IsBroken {
			report.BrokenLinks++
		} else if result.StatusCode >= 300 && result.StatusCode < 400 {
			report.WarningLinks++
		} else {
			report.SuccessLinks++
		}
		if result.IsExternal {
			report.ExternalLinks++
		}
	}

	return report
}
