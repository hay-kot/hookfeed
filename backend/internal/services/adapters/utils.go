package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
)

// GetHeader retrieves a header value, trying multiple possible keys
func GetHeader(r *http.Request, keys ...string) string {
	for _, key := range keys {
		if val := r.Header.Get(key); val != "" {
			return val
		}
	}
	return ""
}

// GetQueryParam retrieves a query parameter value, trying multiple possible keys
func GetQueryParam(query url.Values, keys ...string) string {
	for _, key := range keys {
		if vals, ok := query[key]; ok && len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}
	return ""
}

// ParsePriority converts priority string to int32 (1-5)
// Supports named priorities (min, low, default, high, max, urgent) and numeric values
func ParsePriority(s string) (int32, error) {
	// Handle named priorities (ntfy compatibility)
	switch strings.ToLower(s) {
	case "min", "1":
		return 1, nil
	case "low", "2":
		return 2, nil
	case "default", "3", "":
		return 3, nil
	case "high", "4":
		return 4, nil
	case "max", "urgent", "5":
		return 5, nil
	}

	// Try parsing as integer
	p, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 3, err
	}

	// Clamp to valid range
	if p < 1 {
		p = 1
	}
	if p > 5 {
		p = 5
	}

	return int32(p), nil
}

// SplitAndTrim splits a comma-separated string and trims whitespace from each element
func SplitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	return result
}

// sanitizeSecrets is middleware for copyAndTransformValues that redacts sensitive values
func sanitizeSecrets(key, value string) string {
	keyLower := strings.ToLower(key)

	// Redact common authentication headers
	switch keyLower {
	case "authorization", "proxy-authorization":
		// Handle Bearer, Basic, and other auth schemes
		parts := strings.SplitN(value, " ", 2)
		if len(parts) == 2 {
			return parts[0] + " <redacted>"
		}
		return "<redacted>"

	case "cookie", "set-cookie":
		return "<redacted>"

	case "x-api-key", "x-auth-token", "api-key", "apikey":
		return "<redacted>"
	}

	// Redact query parameters that commonly contain secrets
	if keyLower == "token" || keyLower == "api_key" || keyLower == "apikey" ||
		keyLower == "secret" || keyLower == "password" || keyLower == "key" {
		return "<redacted>"
	}

	return value
}

// feedMessageFromRequest copies all HTTP request data (headers, query params, body) into the FeedMessageCreate
func feedMessageFromRequest(r *http.Request) (dtos.FeedMessageCreate, error) {
	val := dtos.FeedMessageCreateNew()

	rawBody, err := copyBody(r)
	if err != nil {
		return dtos.FeedMessageCreate{}, err
	}

	rawHeaders, err := copyAndTransformValues(r.Header, sanitizeSecrets)
	if err != nil {
		return dtos.FeedMessageCreate{}, err
	}

	rawQueryParams, err := copyAndTransformValues(r.URL.Query(), sanitizeSecrets)
	if err != nil {
		return dtos.FeedMessageCreate{}, err
	}

	val.RawRequest = rawBody
	val.RawHeaders = rawHeaders
	val.RawQueryParams = rawQueryParams

	return val, nil
}

// copyAndTransformValues extracts HTTP key/value pairs (headers, query params, etc.) and returns them as JSON.
// Single-value arrays are unwrapped to strings, multi-value arrays remain as arrays.
// Keys with empty values are omitted. Optional middleware functions can transform/sanitize values.
func copyAndTransformValues(raw map[string][]string, mw ...func(key, value string) string) ([]byte, error) {
	if len(raw) == 0 {
		return []byte("{}"), nil
	}

	execmw := func(key, value string) string {
		for _, fn := range mw {
			value = fn(key, value)
		}

		return value
	}

	result := make(map[string]any)
	for k, v := range raw {
		switch {
		case len(v) == 1:
			result[k] = execmw(k, v[0])
		case len(v) > 1:
			lst := make([]string, len(v))
			for i, vv := range v {
				lst[i] = execmw(k, vv)
			}

			result[k] = lst
		default:
			// nop
		}
	}

	return json.Marshal(result)
}

// copyBody extracts the body from the request and returns it as JSON.
// The body is buffered so it can be read again by the caller.
// If the body is valid JSON, it's returned as-is.
// If the body is raw text (not valid JSON), it's wrapped in {"$body": "text content"}.
func copyBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	r.Body = io.NopCloser(bytes.NewReader(body))

	// If body is completely empty, return empty JSON object
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return []byte("{}"), nil
	}

	// Check if the body is valid JSON
	var js json.RawMessage
	if err := json.Unmarshal(body, &js); err == nil {
		// Body is valid JSON, return as-is
		return body, nil
	}

	// Body is not valid JSON, wrap it in a JSON object with "$body" key
	wrapped := map[string]string{
		"$body": string(body),
	}

	return json.Marshal(wrapped)
}

// isEmptyJSON checks if the given JSON byte slice represents an empty value
// Returns true for: '{}', ”, 'null', '[]', or whitespace-only strings
func isEmptyJSON(data []byte) bool {
	// Trim whitespace
	trimmed := bytes.TrimSpace(data)

	// Check for empty string
	if len(trimmed) == 0 {
		return true
	}

	// Check for common empty JSON representations
	s := string(trimmed)
	return s == "{}" || s == "null" || s == "[]"
}
