package validator

import (
	"fmt"
	"html"
	"os"
	"strings"
	"time"
)

// WriteHTMLReport generates an interactive HTML report with sortable/filterable table
func WriteHTMLReport(report *ValidationReport, filename string) error {
	htmlContent := generateHTMLReport(report)
	return os.WriteFile(filename, []byte(htmlContent), 0644)
}

func generateHTMLReport(report *ValidationReport) string {
	var sb strings.Builder

	// HTML Header
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Link Validation Report - Linkchex</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
        }

        .header h1 {
            font-size: 28px;
            margin-bottom: 10px;
        }

        .header p {
            opacity: 0.9;
            font-size: 14px;
        }

        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 30px;
            background: #f8f9fa;
            border-bottom: 1px solid #e0e0e0;
        }

        .stat-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }

        .stat-card.success {
            border-left-color: #10b981;
        }

        .stat-card.error {
            border-left-color: #ef4444;
        }

        .stat-card.warning {
            border-left-color: #f59e0b;
        }

        .stat-label {
            font-size: 12px;
            text-transform: uppercase;
            color: #6b7280;
            font-weight: 600;
            margin-bottom: 5px;
        }

        .stat-value {
            font-size: 28px;
            font-weight: bold;
            color: #1f2937;
        }

        .controls {
            padding: 20px 30px;
            background: white;
            border-bottom: 1px solid #e0e0e0;
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
            align-items: center;
        }

        .search-box {
            flex: 1;
            min-width: 300px;
        }

        .search-box input {
            width: 100%;
            padding: 10px 15px;
            border: 1px solid #d1d5db;
            border-radius: 6px;
            font-size: 14px;
        }

        .search-box input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .filter-group {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }

        .filter-btn {
            padding: 8px 16px;
            border: 1px solid #d1d5db;
            background: white;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.2s;
        }

        .filter-btn:hover {
            background: #f3f4f6;
        }

        .filter-btn.active {
            background: #667eea;
            color: white;
            border-color: #667eea;
        }

        .table-container {
            overflow-x: auto;
            padding: 30px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 14px;
        }

        thead {
            background: #f9fafb;
            position: sticky;
            top: 0;
        }

        th {
            padding: 12px;
            text-align: left;
            font-weight: 600;
            color: #374151;
            cursor: pointer;
            user-select: none;
            white-space: nowrap;
        }

        th:hover {
            background: #f3f4f6;
        }

        th::after {
            content: ' ↕';
            opacity: 0.3;
            font-size: 12px;
        }

        th.sort-asc::after {
            content: ' ↑';
            opacity: 1;
        }

        th.sort-desc::after {
            content: ' ↓';
            opacity: 1;
        }

        td {
            padding: 12px;
            border-bottom: 1px solid #e5e7eb;
        }

        tr:hover {
            background: #f9fafb;
        }

        .status-badge {
            display: inline-block;
            padding: 4px 10px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: 600;
        }

        .status-success {
            background: #d1fae5;
            color: #065f46;
        }

        .status-error {
            background: #fee2e2;
            color: #991b1b;
        }

        .status-warning {
            background: #fef3c7;
            color: #92400e;
        }

        .url-link {
            color: #667eea;
            text-decoration: none;
            word-break: break-all;
        }

        .url-link:hover {
            text-decoration: underline;
        }

        .tag-badge {
            background: #e0e7ff;
            color: #3730a3;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 600;
            font-family: monospace;
        }

        .no-results {
            text-align: center;
            padding: 60px 20px;
            color: #6b7280;
        }

        .footer {
            padding: 20px 30px;
            background: #f9fafb;
            text-align: center;
            color: #6b7280;
            font-size: 12px;
            border-top: 1px solid #e0e0e0;
        }

        .external-icon::after {
            content: ' ↗';
            font-size: 10px;
            opacity: 0.5;
        }

        .filter-section {
            margin-bottom: 15px;
        }

        .filter-label {
            font-size: 11px;
            font-weight: 600;
            color: #6b7280;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
`)

	// Header
	sb.WriteString(fmt.Sprintf(`
        <div class="header">
            <h1>🔗 Link Validation Report</h1>
            <p>Generated on %s | Duration: %s</p>
        </div>
`, report.StartTime.Format("January 2, 2006 at 3:04 PM"), report.Duration.Round(time.Millisecond)))

	// Summary Statistics
	sb.WriteString(fmt.Sprintf(`
        <div class="summary">
            <div class="stat-card">
                <div class="stat-label">Pages Processed</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Total Links</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Unique URLs</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card success">
                <div class="stat-label">Success</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card error">
                <div class="stat-label">Broken</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card warning">
                <div class="stat-label">Warnings</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Internal Links</div>
                <div class="stat-value">%d</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">External Links</div>
                <div class="stat-value">%d</div>
            </div>
        </div>
`,
		report.PagesProcessed,
		report.TotalLinks,
		report.UniqueURLs,
		report.SuccessLinks,
		report.BrokenLinks,
		report.WarningLinks,
		report.InternalLinks,
		report.ExternalLinks,
	))

	// Controls
	sb.WriteString(`
        <div class="controls">
            <div class="search-box">
                <input type="text" id="searchInput" placeholder="Search by URL, source, or error message...">
            </div>
            <div class="filter-section">
                <div class="filter-label">Filter by Status</div>
                <div class="filter-group">
                    <button class="filter-btn status-filter active" data-filter="all">All</button>
                    <button class="filter-btn status-filter" data-filter="broken">Broken</button>
                    <button class="filter-btn status-filter" data-filter="success">Success</button>
                    <button class="filter-btn status-filter" data-filter="warning">Warnings</button>
                    <button class="filter-btn status-filter" data-filter="external">External</button>
                    <button class="filter-btn status-filter" data-filter="internal">Internal</button>
                </div>
            </div>
            <div class="filter-section">
                <div class="filter-label">Filter by Status Code</div>
                <div class="filter-group">
                    <button class="filter-btn code-filter active" data-code-filter="all">All</button>
                    <button class="filter-btn code-filter" data-code-filter="2xx">2xx Success</button>
                    <button class="filter-btn code-filter" data-code-filter="3xx">3xx Redirects</button>
                    <button class="filter-btn code-filter" data-code-filter="4xx">4xx Client Errors</button>
                    <button class="filter-btn code-filter" data-code-filter="5xx">5xx Server Errors</button>
                    <button class="filter-btn code-filter" data-code-filter="200">200 OK</button>
                    <button class="filter-btn code-filter" data-code-filter="301">301</button>
                    <button class="filter-btn code-filter" data-code-filter="302">302</button>
                    <button class="filter-btn code-filter" data-code-filter="403">403</button>
                    <button class="filter-btn code-filter" data-code-filter="404">404</button>
                    <button class="filter-btn code-filter" data-code-filter="500">500</button>
                    <button class="filter-btn code-filter" data-code-filter="0">Connection Errors</button>
                </div>
            </div>
        </div>
`)

	// Table
	sb.WriteString(`
        <div class="table-container">
            <table id="resultsTable">
                <thead>
                    <tr>
                        <th data-sort="status">Status</th>
                        <th data-sort="target">Target URL</th>
                        <th data-sort="source">Source Page</th>
                        <th data-sort="tag">Tag</th>
                        <th data-sort="code">Code</th>
                        <th data-sort="type">Type</th>
                        <th data-sort="duration">Duration</th>
                    </tr>
                </thead>
                <tbody id="resultsBody">
`)

	// Table rows
	for _, result := range report.Results {
		statusClass := "success"
		statusText := "Success"
		if result.IsBroken {
			statusClass = "error"
			statusText = "Broken"
		} else if result.StatusCode >= 300 && result.StatusCode < 400 {
			statusClass = "warning"
			statusText = "Redirect"
		}

		linkType := "internal"
		if result.IsExternal {
			linkType = "external"
		}

		errorMsg := ""
		if result.Error != nil {
			errorMsg = html.EscapeString(result.Error.Error())
		}

		statusInfo := fmt.Sprintf("%d %s", result.StatusCode, html.EscapeString(result.Status))
		if result.StatusCode == 0 {
			statusInfo = errorMsg
		}

		externalClass := ""
		if result.IsExternal {
			externalClass = "external-icon"
		}

		sb.WriteString(fmt.Sprintf(`
                    <tr data-status="%s" data-type="%s" data-code="%d" data-search="%s">
                        <td><span class="status-badge status-%s">%s</span></td>
                        <td><a href="%s" class="url-link %s" target="_blank" rel="noopener">%s</a></td>
                        <td><a href="%s" class="url-link" target="_blank" rel="noopener">%s</a></td>
                        <td><span class="tag-badge">&lt;%s&gt;</span></td>
                        <td>%s</td>
                        <td>%s</td>
                        <td>%dms</td>
                    </tr>
`,
			statusText,
			linkType,
			result.StatusCode,
			strings.ToLower(result.TargetURL+" "+result.SourceURL+" "+errorMsg),
			statusClass,
			statusText,
			html.EscapeString(result.TargetURL),
			externalClass,
			html.EscapeString(truncate(result.TargetURL, 80)),
			html.EscapeString(result.SourceURL),
			html.EscapeString(truncate(result.SourceURL, 60)),
			html.EscapeString(result.Tag),
			statusInfo,
			linkType,
			result.Duration.Milliseconds(),
		))
	}

	sb.WriteString(`
                </tbody>
            </table>
            <div id="noResults" class="no-results" style="display: none;">
                <p>No results match your search or filter criteria.</p>
            </div>
        </div>
`)

	// Footer
	sb.WriteString(`
        <div class="footer">
            Generated by <a href="https://github.com/anthropics/claude-code" target="_blank" style="color: #667eea;">Linkchex</a> -
            A high-performance link validation tool
        </div>
    </div>
`)

	// JavaScript
	sb.WriteString(`
    <script>
        // Table sorting and filtering
        const table = document.getElementById('resultsTable');
        const tbody = document.getElementById('resultsBody');
        const searchInput = document.getElementById('searchInput');
        const statusFilterBtns = document.querySelectorAll('.status-filter');
        const codeFilterBtns = document.querySelectorAll('.code-filter');
        const noResults = document.getElementById('noResults');

        let currentSort = { column: null, direction: 'asc' };
        let currentFilter = 'all';
        let currentCodeFilter = 'all';
        let searchTerm = '';

        // Sorting
        table.querySelectorAll('th[data-sort]').forEach(th => {
            th.addEventListener('click', () => {
                const column = th.dataset.sort;

                if (currentSort.column === column) {
                    currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
                } else {
                    currentSort.column = column;
                    currentSort.direction = 'asc';
                }

                // Update header classes
                table.querySelectorAll('th').forEach(h => {
                    h.classList.remove('sort-asc', 'sort-desc');
                });
                th.classList.add('sort-' + currentSort.direction);

                sortTable();
            });
        });

        function sortTable() {
            const rows = Array.from(tbody.querySelectorAll('tr'));

            rows.sort((a, b) => {
                let aVal, bVal;

                switch(currentSort.column) {
                    case 'status':
                        aVal = a.querySelector('.status-badge').textContent;
                        bVal = b.querySelector('.status-badge').textContent;
                        break;
                    case 'target':
                        aVal = a.querySelectorAll('td')[1].textContent;
                        bVal = b.querySelectorAll('td')[1].textContent;
                        break;
                    case 'source':
                        aVal = a.querySelectorAll('td')[2].textContent;
                        bVal = b.querySelectorAll('td')[2].textContent;
                        break;
                    case 'tag':
                        aVal = a.querySelectorAll('td')[3].textContent;
                        bVal = b.querySelectorAll('td')[3].textContent;
                        break;
                    case 'code':
                        aVal = a.querySelectorAll('td')[4].textContent;
                        bVal = b.querySelectorAll('td')[4].textContent;
                        break;
                    case 'type':
                        aVal = a.querySelectorAll('td')[5].textContent;
                        bVal = b.querySelectorAll('td')[5].textContent;
                        break;
                    case 'duration':
                        aVal = parseInt(a.querySelectorAll('td')[6].textContent);
                        bVal = parseInt(b.querySelectorAll('td')[6].textContent);
                        break;
                }

                if (currentSort.direction === 'asc') {
                    return aVal > bVal ? 1 : -1;
                } else {
                    return aVal < bVal ? 1 : -1;
                }
            });

            rows.forEach(row => tbody.appendChild(row));
        }

        // Status Filtering
        statusFilterBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                statusFilterBtns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                currentFilter = btn.dataset.filter;
                applyFilters();
            });
        });

        // Status Code Filtering
        codeFilterBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                codeFilterBtns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                currentCodeFilter = btn.dataset.codeFilter;
                applyFilters();
            });
        });

        // Search
        searchInput.addEventListener('input', (e) => {
            searchTerm = e.target.value.toLowerCase();
            applyFilters();
        });

        function applyFilters() {
            const rows = tbody.querySelectorAll('tr');
            let visibleCount = 0;

            rows.forEach(row => {
                let show = true;

                // Filter by status/type
                if (currentFilter !== 'all') {
                    const status = row.dataset.status.toLowerCase();
                    const type = row.dataset.type.toLowerCase();

                    if (currentFilter === 'broken' && status !== 'broken') show = false;
                    if (currentFilter === 'success' && status !== 'success') show = false;
                    if (currentFilter === 'warning' && status !== 'redirect') show = false;
                    if (currentFilter === 'external' && type !== 'external') show = false;
                    if (currentFilter === 'internal' && type !== 'internal') show = false;
                }

                // Filter by status code
                if (currentCodeFilter !== 'all' && show) {
                    const code = parseInt(row.dataset.code);

                    // Handle range filters
                    if (currentCodeFilter === '2xx' && (code < 200 || code >= 300)) {
                        show = false;
                    } else if (currentCodeFilter === '3xx' && (code < 300 || code >= 400)) {
                        show = false;
                    } else if (currentCodeFilter === '4xx' && (code < 400 || code >= 500)) {
                        show = false;
                    } else if (currentCodeFilter === '5xx' && (code < 500 || code >= 600)) {
                        show = false;
                    } else if (currentCodeFilter !== '2xx' && currentCodeFilter !== '3xx' &&
                               currentCodeFilter !== '4xx' && currentCodeFilter !== '5xx') {
                        // Exact match for specific codes
                        if (code.toString() !== currentCodeFilter) {
                            show = false;
                        }
                    }
                }

                // Filter by search term
                if (searchTerm && show) {
                    const searchData = row.dataset.search;
                    if (!searchData.includes(searchTerm)) {
                        show = false;
                    }
                }

                row.style.display = show ? '' : 'none';
                if (show) visibleCount++;
            });

            // Show/hide no results message
            if (visibleCount === 0) {
                table.style.display = 'none';
                noResults.style.display = 'block';
            } else {
                table.style.display = 'table';
                noResults.style.display = 'none';
            }
        }
    </script>
</body>
</html>
`)

	return sb.String()
}
