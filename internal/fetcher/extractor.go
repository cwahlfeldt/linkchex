package fetcher

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an extracted link from HTML
type Link struct {
	URL      string
	Tag      string // a, img, link, script
	Attr     string // href, src
	Text     string // Link text for <a> tags
	IsExternal bool
}

// ExtractLinks extracts all links from HTML content
func ExtractLinks(htmlContent []byte, baseURL string, skipResources bool) ([]Link, error) {
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	var links []Link
	var extract func(*html.Node)

	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var link Link

			switch n.Data {
			case "a":
				if href := getAttr(n, "href"); href != "" {
					link = Link{
						URL:  href,
						Tag:  "a",
						Attr: "href",
						Text: extractText(n),
					}
				}
			case "img":
				if src := getAttr(n, "src"); src != "" {
					link = Link{
						URL:  src,
						Tag:  "img",
						Attr: "src",
					}
				}
			case "link":
				if skipResources {
					break // Skip <link> tags if skipResources is true
				}
				if href := getAttr(n, "href"); href != "" {
					link = Link{
						URL:  href,
						Tag:  "link",
						Attr: "href",
					}
				}
			case "script":
				if skipResources {
					break // Skip <script> tags if skipResources is true
				}
				if src := getAttr(n, "src"); src != "" {
					link = Link{
						URL:  src,
						Tag:  "script",
						Attr: "src",
					}
				}
			}

			if link.URL != "" {
				// Resolve relative URLs
				parsedURL, err := url.Parse(link.URL)
				if err == nil {
					absoluteURL := base.ResolveReference(parsedURL)
					link.URL = absoluteURL.String()
					link.IsExternal = isExternalLink(base, absoluteURL)
				}
				links = append(links, link)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	return links, nil
}

// getAttr gets an attribute value from an HTML node
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return strings.TrimSpace(attr.Val)
		}
	}
	return ""
}

// extractText extracts text content from a node and its children
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c) + " "
	}
	return strings.TrimSpace(text)
}

// isExternalLink checks if a link points to an external domain
func isExternalLink(base, target *url.URL) bool {
	return base.Host != target.Host
}

// FilterLinks filters links based on criteria
func FilterLinks(links []Link, includeExternal bool) []Link {
	var filtered []Link
	seen := make(map[string]bool)

	for _, link := range links {
		// Skip duplicates
		if seen[link.URL] {
			continue
		}
		seen[link.URL] = true

		// Skip external links if not included
		if !includeExternal && link.IsExternal {
			continue
		}

		// Skip javascript:, mailto:, tel:, etc.
		if strings.HasPrefix(link.URL, "javascript:") ||
			strings.HasPrefix(link.URL, "mailto:") ||
			strings.HasPrefix(link.URL, "tel:") ||
			strings.HasPrefix(link.URL, "#") {
			continue
		}

		filtered = append(filtered, link)
	}

	return filtered
}
