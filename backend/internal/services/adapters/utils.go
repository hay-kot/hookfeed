package adapters

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// ParseBool converts a string to bool, accepting various formats
func ParseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
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
