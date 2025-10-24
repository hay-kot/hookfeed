package adapters

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHeader(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		keys     []string
		expected string
	}{
		{
			name:     "single key found",
			headers:  map[string]string{"X-Title": "Test Title"},
			keys:     []string{"X-Title"},
			expected: "Test Title",
		},
		{
			name:     "multiple keys, first found",
			headers:  map[string]string{"X-Title": "Test Title"},
			keys:     []string{"X-Title", "Title"},
			expected: "Test Title",
		},
		{
			name:     "multiple keys, second found",
			headers:  map[string]string{"Title": "Test Title"},
			keys:     []string{"X-Title", "Title"},
			expected: "Test Title",
		},
		{
			name:     "no keys found",
			headers:  map[string]string{"Other": "Value"},
			keys:     []string{"X-Title", "Title"},
			expected: "",
		},
		{
			name:     "empty headers",
			headers:  map[string]string{},
			keys:     []string{"X-Title"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: make(http.Header),
			}
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			result := GetHeader(req, tt.keys...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetQueryParam(t *testing.T) {
	tests := []struct {
		name     string
		query    map[string][]string
		keys     []string
		expected string
	}{
		{
			name:     "single key found",
			query:    map[string][]string{"title": {"Test Title"}},
			keys:     []string{"title"},
			expected: "Test Title",
		},
		{
			name:     "multiple keys, first found",
			query:    map[string][]string{"title": {"Test Title"}},
			keys:     []string{"title", "t"},
			expected: "Test Title",
		},
		{
			name:     "multiple keys, second found",
			query:    map[string][]string{"t": {"Test Title"}},
			keys:     []string{"title", "t"},
			expected: "Test Title",
		},
		{
			name:     "no keys found",
			query:    map[string][]string{"other": {"Value"}},
			keys:     []string{"title", "t"},
			expected: "",
		},
		{
			name:     "empty query",
			query:    map[string][]string{},
			keys:     []string{"title"},
			expected: "",
		},
		{
			name:     "empty value",
			query:    map[string][]string{"title": {""}},
			keys:     []string{"title"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := url.Values(tt.query)
			result := GetQueryParam(query, tt.keys...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int32
		wantErr  bool
	}{
		{name: "named min", input: "min", expected: 1, wantErr: false},
		{name: "named low", input: "low", expected: 2, wantErr: false},
		{name: "named default", input: "default", expected: 3, wantErr: false},
		{name: "named high", input: "high", expected: 4, wantErr: false},
		{name: "named max", input: "max", expected: 5, wantErr: false},
		{name: "named urgent", input: "urgent", expected: 5, wantErr: false},
		{name: "numeric 1", input: "1", expected: 1, wantErr: false},
		{name: "numeric 2", input: "2", expected: 2, wantErr: false},
		{name: "numeric 3", input: "3", expected: 3, wantErr: false},
		{name: "numeric 4", input: "4", expected: 4, wantErr: false},
		{name: "numeric 5", input: "5", expected: 5, wantErr: false},
		{name: "empty string", input: "", expected: 3, wantErr: false},
		{name: "clamped low", input: "0", expected: 1, wantErr: false},
		{name: "clamped high", input: "10", expected: 5, wantErr: false},
		{name: "case insensitive MAX", input: "MAX", expected: 5, wantErr: false},
		{name: "case insensitive Low", input: "Low", expected: 2, wantErr: false},
		{name: "invalid string", input: "invalid", expected: 3, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePriority(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "true lowercase", input: "true", expected: true},
		{name: "true uppercase", input: "TRUE", expected: true},
		{name: "true mixed case", input: "True", expected: true},
		{name: "1", input: "1", expected: true},
		{name: "yes lowercase", input: "yes", expected: true},
		{name: "yes uppercase", input: "YES", expected: true},
		{name: "false", input: "false", expected: false},
		{name: "0", input: "0", expected: false},
		{name: "no", input: "no", expected: false},
		{name: "empty", input: "", expected: false},
		{name: "random string", input: "random", expected: false},
		{name: "with spaces", input: "  true  ", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseBool(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single value",
			input:    "tag1",
			expected: []string{"tag1"},
		},
		{
			name:     "multiple values",
			input:    "tag1,tag2,tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "values with spaces",
			input:    "tag1 , tag2 , tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single value with spaces",
			input:    "  tag1  ",
			expected: []string{"tag1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitAndTrim(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
