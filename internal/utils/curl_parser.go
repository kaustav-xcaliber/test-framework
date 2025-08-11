package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// CurlRequest represents a parsed curl command
type CurlRequest struct {
	Method        string            `json:"method"`
	URL           string            `json:"url"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	QueryParams   map[string]string `json:"queryParams"`
	PathVariables map[string]string `json:"pathVariables"`
	RequestType   string            `json:"requestType"`
	RawCommand    string            `json:"rawCommand"`
}

// ParseCurlCommand parses a curl command and returns a CurlRequest
func ParseCurlCommand(curlCmd string) (*CurlRequest, error) {
	if curlCmd == "" {
		return nil, fmt.Errorf("curl command cannot be empty")
	}

	// 1️⃣ Normalize multi-line commands
	command := strings.ReplaceAll(curlCmd, "\\\n", " ")
	command = strings.ReplaceAll(command, "\\", "")
	command = strings.Join(strings.Fields(command), " ") // collapse multiple spaces

	result := &CurlRequest{
		Method:        "GET",
		Headers:       make(map[string]string),
		QueryParams:   make(map[string]string),
		PathVariables: make(map[string]string),
		RawCommand:    curlCmd,
	}

	// 2️⃣ Tokenize, preserving quotes
	re := regexp.MustCompile(`'[^']*'|"[^"]*"|\S+`)
	tokens := re.FindAllString(command, -1)

	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid curl command")
	}

	for i := 0; i < len(tokens); i++ {
		token := stripQuotes(tokens[i])

		switch token {
		case "curl":
			continue
		case "-X", "--request":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing method after -X/--request")
			}
			i++
			result.Method = strings.ToUpper(stripQuotes(tokens[i]))
		case "-H", "--header":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing header value after -H/--header")
			}
			i++
			header := stripQuotes(tokens[i])
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				result.Headers[key] = value
			} else {
				return nil, fmt.Errorf("invalid header format: %s", header)
			}
		case "--location", "-L":
			// --location follows redirects, but doesn't change the request
			continue
		case "--data", "--data-raw", "--data-binary", "-d":
			if result.Method == "GET" {
				result.Method = "POST"
			}
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing data after --data/-d")
			}
			i++
			result.Body = stripQuotes(tokens[i])
		case "--form", "-F":
			if result.Method == "GET" {
				result.Method = "POST"
			}
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing form data after --form/-F")
			}
			i++
			formData := stripQuotes(tokens[i])
			parts := strings.SplitN(formData, "=", 2)
			if len(parts) == 2 {
				if result.Body == "" {
					result.Body = fmt.Sprintf("%s=%s", parts[0], parts[1])
				} else {
					result.Body += fmt.Sprintf("&%s=%s", parts[0], parts[1])
				}
			}
		case "--user", "-u":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing credentials after --user/-u")
			}
			i++
			credentials := stripQuotes(tokens[i])
			parts := strings.SplitN(credentials, ":", 2)
			if len(parts) == 2 {
				result.Headers["Authorization"] = fmt.Sprintf("Basic %s", credentials)
			}
		case "--cookie", "-b":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing cookie value after --cookie/-b")
			}
			i++
			cookie := stripQuotes(tokens[i])
			result.Headers["Cookie"] = cookie
		case "--compressed", "-C":
			// Handle compression flag
			continue
		case "--insecure", "-k":
			// Handle insecure flag
			continue
		case "--silent", "-s":
			// Handle silent flag
			continue
		case "--verbose", "-v":
			// Handle verbose flag
			continue
		default:
			// URL detection
			if strings.HasPrefix(token, "http://") || strings.HasPrefix(token, "https://") {
				result.URL = token
			} else if strings.HasPrefix(token, "localhost:") || strings.HasPrefix(token, "127.0.0.1:") {
				// Handle localhost URLs
				result.URL = "http://" + token
			}
		}
	}

	// 3️⃣ Validate URL
	if result.URL == "" {
		return nil, fmt.Errorf("no URL found in curl command")
	}

	// 4️⃣ Extract query params
	if strings.Contains(result.URL, "?") {
		parsedURL, err := url.Parse(result.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %v", err)
		}
		result.URL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
		for k, v := range parsedURL.Query() {
			if len(v) > 0 {
				result.QueryParams[k] = v[0]
			} else {
				result.QueryParams[k] = ""
			}
		}
	}

	// 5️⃣ Extract path variables like {id}
	pathVarRe := regexp.MustCompile(`\{([^}]+)\}`)
	matches := pathVarRe.FindAllStringSubmatch(result.URL, -1)
	for _, m := range matches {
		if len(m) > 1 {
			result.PathVariables[m[1]] = ""
		}
	}

	// 6️⃣ Classify request type
	result.RequestType = classifyRequest(result)

	return result, nil
}

// stripQuotes removes surrounding quotes from a string
func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) ||
		(strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) {
		return s[1 : len(s)-1]
	}
	return s
}

// classifyRequest determines the type of request based on URL and method
func classifyRequest(req *CurlRequest) string {
	urlLower := strings.ToLower(req.URL)
	method := strings.ToUpper(req.Method)

	switch {
	case strings.Contains(urlLower, "/$export"):
		return "FHIR Bulk Export"
	case strings.Contains(urlLower, "/fhir/") && method == "GET":
		return "FHIR Read"
	case strings.Contains(urlLower, "/fhir/") && method == "POST":
		return "FHIR Create/Action"
	case strings.Contains(urlLower, "/fhir/") && method == "PUT":
		return "FHIR Update"
	case strings.Contains(urlLower, "/fhir/") && method == "DELETE":
		return "FHIR Delete"
	case method == "GET":
		return "Fetch"
	case method == "POST":
		return "Submit"
	case method == "PUT":
		return "Update"
	case method == "DELETE":
		return "Delete"
	case method == "PATCH":
		return "Patch"
	default:
		return "Other"
	}
}

// ToTestSpec converts a CurlRequest to a TestSpec for the framework
func (c *CurlRequest) ToTestSpec(name, description string) map[string]interface{} {
	// Build the full URL with query parameters
	fullURL := c.URL
	if len(c.QueryParams) > 0 {
		queryParts := make([]string, 0, len(c.QueryParams))
		for k, v := range c.QueryParams {
			if v != "" {
				queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
			} else {
				queryParts = append(queryParts, k)
			}
		}
		fullURL += "?" + strings.Join(queryParts, "&")
	}

	// Convert body to interface{} for JSON handling
	var body interface{}
	if c.Body != "" {
		// Try to parse as JSON first
		if err := json.Unmarshal([]byte(c.Body), &body); err != nil {
			// If not JSON, treat as string
			body = c.Body
		}
	}

	return map[string]interface{}{
		"name":        name,
		"description": description,
		"service_name": "curl-service", // This will be overridden by the actual service
		"request": map[string]interface{}{
			"method":  c.Method,
			"url":     fullURL,
			"headers": c.Headers,
			"body":    body,
		},
		"assertions": []map[string]interface{}{
			{
				"type":     "status_code",
				"expected": 200,
			},
		},
	}
}

// String returns a string representation of the CurlRequest
func (c *CurlRequest) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}
