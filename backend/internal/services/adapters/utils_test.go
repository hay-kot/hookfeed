package adapters

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetHeader(t *testing.T) {
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

func Test_GetQueryParam(t *testing.T) {
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

func Test_ParsePriority(t *testing.T) {
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

func Test_ParseBool(t *testing.T) {
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

func Test_SplitAndTrim(t *testing.T) {
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

func Test_copyAndTransformValues(t *testing.T) {
	tests := []struct {
		name       string
		input      map[string][]string
		middleware []func(key, value string) string
		expected   string
		wantErr    bool
	}{
		{
			name:     "empty input",
			input:    map[string][]string{},
			expected: "{}",
			wantErr:  false,
		},
		{
			name: "single value unwrapped",
			input: map[string][]string{
				"key1": {"value1"},
			},
			expected: `{"key1":"value1"}`,
			wantErr:  false,
		},
		{
			name: "multiple values remain as array",
			input: map[string][]string{
				"key1": {"value1", "value2", "value3"},
			},
			expected: `{"key1":["value1","value2","value3"]}`,
			wantErr:  false,
		},
		{
			name: "mixed single and multiple values",
			input: map[string][]string{
				"single": {"one"},
				"multi":  {"first", "second"},
			},
			expected: `{"multi":["first","second"],"single":"one"}`,
			wantErr:  false,
		},
		{
			name: "empty value array omitted",
			input: map[string][]string{
				"key1":  {"value1"},
				"empty": {},
				"key2":  {"value2"},
			},
			expected: `{"key1":"value1","key2":"value2"}`,
			wantErr:  false,
		},
		{
			name: "with middleware single function",
			input: map[string][]string{
				"normal": {"test"},
				"secret": {"password123"},
			},
			middleware: []func(key, value string) string{
				func(key, value string) string {
					if key == "secret" {
						return "<redacted>"
					}
					return value
				},
			},
			expected: `{"normal":"test","secret":"<redacted>"}`,
			wantErr:  false,
		},
		{
			name: "with middleware multiple functions",
			input: map[string][]string{
				"text": {"hello"},
			},
			middleware: []func(key, value string) string{
				func(key, value string) string {
					return value + "-mw1"
				},
				func(key, value string) string {
					return value + "-mw2"
				},
			},
			expected: `{"text":"hello-mw1-mw2"}`,
			wantErr:  false,
		},
		{
			name: "middleware applied to multi-value arrays",
			input: map[string][]string{
				"values": {"a", "b", "c"},
			},
			middleware: []func(key, value string) string{
				func(key, value string) string {
					return value + "-modified"
				},
			},
			expected: `{"values":["a-modified","b-modified","c-modified"]}`,
			wantErr:  false,
		},
		{
			name: "sanitize secrets middleware",
			input: map[string][]string{
				"Authorization": {"Bearer abc123"},
				"X-Api-Key":     {"secret-key"},
				"Content-Type":  {"application/json"},
			},
			middleware: []func(key, value string) string{sanitizeSecrets},
			expected:   `{"Authorization":"Bearer <redacted>","Content-Type":"application/json","X-Api-Key":"<redacted>"}`,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := copyAndTransformValues(tt.input, tt.middleware...)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}
