package sitemap

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// URLSet represents the root element of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single URL entry in a sitemap
type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

// SitemapIndex represents a sitemap index file
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// Sitemap represents a sitemap reference in an index file
type Sitemap struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// Parse parses a sitemap URL or local file and returns all URLs
// Handles both regular sitemaps and sitemap index files
func Parse(sitemapURL string) ([]string, error) {
	var reader io.ReadCloser
	var err error

	// Check if it's a local file or remote URL
	if strings.HasPrefix(sitemapURL, "http://") || strings.HasPrefix(sitemapURL, "https://") {
		reader, err = fetchRemoteSitemap(sitemapURL)
	} else {
		reader, err = openLocalSitemap(sitemapURL)
	}

	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read sitemap: %w", err)
	}

	// Try parsing as sitemap index first
	var sitemapIndex SitemapIndex
	if err := xml.Unmarshal(content, &sitemapIndex); err == nil && len(sitemapIndex.Sitemaps) > 0 {
		// It's a sitemap index, recursively parse each sitemap
		return parseSitemapIndex(&sitemapIndex)
	}

	// Parse as regular sitemap
	var urlSet URLSet
	if err := xml.Unmarshal(content, &urlSet); err != nil {
		return nil, fmt.Errorf("failed to parse sitemap XML: %w", err)
	}

	// Extract URLs
	var urls []string
	for _, url := range urlSet.URLs {
		if url.Loc != "" {
			urls = append(urls, url.Loc)
		}
	}

	return urls, nil
}

// parseSitemapIndex recursively parses a sitemap index file
func parseSitemapIndex(index *SitemapIndex) ([]string, error) {
	var allURLs []string

	for _, sitemap := range index.Sitemaps {
		urls, err := Parse(sitemap.Loc)
		if err != nil {
			// Log error but continue with other sitemaps
			fmt.Fprintf(os.Stderr, "Warning: Failed to parse sitemap %s: %v\n", sitemap.Loc, err)
			continue
		}
		allURLs = append(allURLs, urls...)
	}

	return allURLs, nil
}

// fetchRemoteSitemap fetches a sitemap from a remote URL
func fetchRemoteSitemap(url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("sitemap returned status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// openLocalSitemap opens a local sitemap file
func openLocalSitemap(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sitemap file: %w", err)
	}
	return file, nil
}
