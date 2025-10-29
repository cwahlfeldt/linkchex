package validator

import (
	"regexp"
	"strings"
)

// URLMatcher matches URLs against patterns
type URLMatcher struct {
	excludePatterns []*regexp.Regexp
	includePatterns []*regexp.Regexp
}

// NewURLMatcher creates a new URL matcher
func NewURLMatcher(excludePatterns, includePatterns []string) (*URLMatcher, error) {
	matcher := &URLMatcher{
		excludePatterns: make([]*regexp.Regexp, 0),
		includePatterns: make([]*regexp.Regexp, 0),
	}

	// Compile exclude patterns
	for _, pattern := range excludePatterns {
		re, err := compilePattern(pattern)
		if err != nil {
			return nil, err
		}
		matcher.excludePatterns = append(matcher.excludePatterns, re)
	}

	// Compile include patterns
	for _, pattern := range includePatterns {
		re, err := compilePattern(pattern)
		if err != nil {
			return nil, err
		}
		matcher.includePatterns = append(matcher.includePatterns, re)
	}

	return matcher, nil
}

// compilePattern compiles a pattern string into a regex
// Supports glob-like patterns (* and ?) or full regex
func compilePattern(pattern string) (*regexp.Regexp, error) {
	// If pattern starts with ^, treat as regex
	if strings.HasPrefix(pattern, "^") {
		return regexp.Compile(pattern)
	}

	// Convert glob pattern to regex
	// Escape special regex characters except * and ?
	escaped := regexp.QuoteMeta(pattern)
	// Convert glob wildcards to regex
	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = strings.ReplaceAll(escaped, "\\?", ".")
	// Anchor the pattern
	escaped = "^" + escaped + "$"

	return regexp.Compile(escaped)
}

// ShouldCheck determines if a URL should be checked based on patterns
func (m *URLMatcher) ShouldCheck(url string) bool {
	// If include patterns are defined, URL must match at least one
	if len(m.includePatterns) > 0 {
		matched := false
		for _, pattern := range m.includePatterns {
			if pattern.MatchString(url) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check exclude patterns
	for _, pattern := range m.excludePatterns {
		if pattern.MatchString(url) {
			return false
		}
	}

	return true
}

// DefaultExcludePatterns returns common patterns to exclude
func DefaultExcludePatterns() []string {
	return []string{
		"*.pdf",
		"*.zip",
		"*.tar.gz",
		"*.exe",
		"*.dmg",
		"*/admin/*",
		"*/wp-admin/*",
		"*/wp-login.php",
		"*/login",
		"*/logout",
		"*/signin",
		"*/signout",
	}
}
