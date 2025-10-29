package sitemap

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Discover attempts to find sitemaps from a base URL
// It checks common locations and robots.txt
func Discover(baseURL string) ([]string, error) {
	// Ensure baseURL has a scheme
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	// Parse base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	baseURL = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	var sitemaps []string

	// Check robots.txt first
	robotsSitemaps, err := checkRobotsTxt(baseURL)
	if err == nil && len(robotsSitemaps) > 0 {
		sitemaps = append(sitemaps, robotsSitemaps...)
	}

	// If no sitemaps found in robots.txt, check common locations
	if len(sitemaps) == 0 {
		commonPaths := []string{
			"/sitemap.xml",
			"/sitemap_index.xml",
			"/sitemap/sitemap.xml",
			"/sitemap/index.xml",
		}

		var lastErr error
		for _, path := range commonPaths {
			sitemapURL := baseURL + path
			exists, err := urlExists(sitemapURL)
			if err != nil {
				lastErr = err
				continue
			}
			if exists {
				sitemaps = append(sitemaps, sitemapURL)
				break // Found one, stop searching
			}
		}

		if len(sitemaps) == 0 && lastErr != nil {
			return nil, fmt.Errorf("no sitemap found at %s (last error: %v)", baseURL, lastErr)
		}
	}

	if len(sitemaps) == 0 {
		return nil, fmt.Errorf("no sitemap found at %s", baseURL)
	}

	return sitemaps, nil
}

// checkRobotsTxt parses robots.txt and extracts Sitemap directives
func checkRobotsTxt(baseURL string) ([]string, error) {
	robotsURL := baseURL + "/robots.txt"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(robotsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("robots.txt returned status %d", resp.StatusCode)
	}

	var sitemaps []string
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			// Extract sitemap URL
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				sitemapURL := strings.TrimSpace(parts[1])
				sitemaps = append(sitemaps, sitemapURL)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sitemaps, nil
}

// urlExists checks if a URL returns a successful status code
func urlExists(urlStr string) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Use HEAD request first (faster)
	resp, err := client.Head(urlStr)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Some servers don't support HEAD, fallback to GET if needed
	if resp.StatusCode == http.StatusMethodNotAllowed {
		resp, err = client.Get(urlStr)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
	}

	// Return error for rate limiting so caller knows what happened
	if resp.StatusCode == http.StatusTooManyRequests {
		return false, fmt.Errorf("rate limited (429)")
	}

	// Accept 2xx status codes as success
	return resp.StatusCode >= 200 && resp.StatusCode < 300, nil
}
