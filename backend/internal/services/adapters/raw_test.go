package adapters

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawAdapter_UnmarshalRequest_ValidJSON(t *testing.T) {
	title := "Test Title"
	message := "Test Message"
	priority := int32(4)

	createDTO := dtos.FeedMessageCreate{
		FeedID:   "test-feed",
		Title:    title,
		Message:  message,
		Priority: priority,
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, "test-feed", adapter.CreateDTO.FeedID)
	assert.Equal(t, title, adapter.CreateDTO.Title)
	assert.Equal(t, message, adapter.CreateDTO.Message)
	assert.Equal(t, priority, adapter.CreateDTO.Priority)
}

func TestRawAdapter_UnmarshalRequest_WithRawRequest(t *testing.T) {
	rawRequest := json.RawMessage(`{"custom":"data"}`)
	createDTO := dtos.FeedMessageCreate{
		FeedID:     "test-feed",
		RawRequest: rawRequest,
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	// Should preserve the provided RawRequest
	assert.JSONEq(t, `{"custom":"data"}`, string(adapter.CreateDTO.RawRequest))
}

func TestRawAdapter_UnmarshalRequest_AutoPopulateRawRequest(t *testing.T) {
	createDTO := dtos.FeedMessageCreate{
		FeedID: "test-feed",
		// No RawRequest provided
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	// Should auto-populate RawRequest with the body
	assert.NotNil(t, adapter.CreateDTO.RawRequest)

	// Verify it's valid JSON
	assert.True(t, json.Valid(adapter.CreateDTO.RawRequest), "Auto-populated RawRequest should be valid JSON")
}

func TestRawAdapter_UnmarshalRequest_AutoPopulateRawHeaders(t *testing.T) {
	createDTO := dtos.FeedMessageCreate{
		FeedID: "test-feed",
		// No RawHeaders provided
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "custom-value")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	// Check that RawHeaders map was populated
	assert.NotEmpty(t, adapter.RawHeaders, "adapter.RawHeaders should be populated from request headers")

	// Should auto-populate RawHeaders
	assert.NotNil(t, adapter.CreateDTO.RawHeaders)

	// NewFeedMessageCreateFromHTTP stores headers as map[string][]string
	var headers map[string][]string
	err = json.Unmarshal(adapter.CreateDTO.RawHeaders, &headers)
	require.NoError(t, err)

	assert.Equal(t, []string{"application/json"}, headers["Content-Type"])
	assert.Equal(t, []string{"custom-value"}, headers["X-Custom-Header"])
}

func TestRawAdapter_UnmarshalRequest_WithRawHeaders(t *testing.T) {
	rawHeaders := json.RawMessage(`{"custom":"headers"}`)
	createDTO := dtos.FeedMessageCreate{
		FeedID:     "test-feed",
		RawHeaders: rawHeaders,
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	// Should preserve the provided RawHeaders
	assert.JSONEq(t, `{"custom":"headers"}`, string(adapter.CreateDTO.RawHeaders))
}

func TestRawAdapter_UnmarshalRequest_PlainTextBody(t *testing.T) {
	plainText := "This is plain text"
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(plainText))
	req.Header.Set("Content-Type", "text/plain")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)

	// Should fail to unmarshal plain text as JSON
	assert.Error(t, err)
}

func TestRawAdapter_UnmarshalRequest_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid json`
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	assert.Error(t, err)
}

func TestRawAdapter_UnmarshalRequest_CapturesHeaders(t *testing.T) {
	createDTO := dtos.FeedMessageCreate{
		FeedID: "test-feed",
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-1", "value1")
	req.Header.Set("X-Custom-2", "value2")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.NotEmpty(t, adapter.RawHeaders)
	assert.Equal(t, "application/json", adapter.RawHeaders["Content-Type"])
	assert.Equal(t, "value1", adapter.RawHeaders["X-Custom-1"])
	assert.Equal(t, "value2", adapter.RawHeaders["X-Custom-2"])
}

func TestRawAdapter_UnmarshalRequest_MultiValueHeaders(t *testing.T) {
	createDTO := dtos.FeedMessageCreate{
		FeedID: "test-feed",
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header["X-Multi"] = []string{"value1", "value2", "value3"}

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	// Should only capture first value for multi-value headers
	assert.Equal(t, "value1", adapter.RawHeaders["X-Multi"], "should capture first value for multi-value headers")
}

func TestRawAdapter_AsFeedMessage(t *testing.T) {
	title := "Test Title"
	message := "Test Message"
	priority := int32(4)

	adapter := &RawAdapter{
		CreateDTO: dtos.FeedMessageCreate{
			FeedID:   "test-feed",
			Title:    title,
			Message:  message,
			Priority: priority,
		},
	}

	dto := adapter.AsFeedMessage()

	assert.Equal(t, "test-feed", dto.FeedID)
	assert.Equal(t, title, dto.Title)
	assert.Equal(t, message, dto.Message)
	assert.Equal(t, priority, dto.Priority)
}

func TestRawAdapter_UnmarshalRequest_ReadBodyError(t *testing.T) {
	// Create a reader that always returns an error
	errReader := io.NopCloser(&errorReader{})
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", errReader)

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	assert.Error(t, err)
}

// errorReader is a helper that always returns an error when reading
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestRawAdapter_UnmarshalRequest_WithMetadata(t *testing.T) {
	metadata := json.RawMessage(`{"custom":"metadata","version":1}`)
	createDTO := dtos.FeedMessageCreate{
		FeedID:   "test-feed",
		Metadata: metadata,
	}

	bodyBytes, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/feed-messages", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")

	adapter := &RawAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	// Verify metadata is preserved
	assert.JSONEq(t, string(metadata), string(adapter.CreateDTO.Metadata))

	// Verify AsFeedMessage preserves metadata
	dto := adapter.AsFeedMessage()
	assert.JSONEq(t, string(metadata), string(dto.Metadata))
}
