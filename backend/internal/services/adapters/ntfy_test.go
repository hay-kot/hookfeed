package adapters

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNtfyAdapter_UnmarshalRequest_PlainText(t *testing.T) {
	body := "This is a plain text message"
	req := httptest.NewRequest(http.MethodPost, "/test-topic", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")

	// Set up chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, body, adapter.Message.Message)
	assert.Equal(t, "test-topic", adapter.Message.Topic)
	assert.Equal(t, int32(3), adapter.Message.Priority)
}

func TestNtfyAdapter_UnmarshalRequest_JSON(t *testing.T) {
	msg := NtfyMessage{
		Topic:    "test-topic",
		Message:  "JSON message",
		Title:    "JSON Title",
		Priority: 5,
		Tags:     []string{"tag1", "tag2"},
		Click:    "https://example.com",
		Icon:     "https://example.com/icon.png",
		Markdown: true,
	}
	bodyBytes, _ := json.Marshal(msg)

	req := httptest.NewRequest(http.MethodPost, "/test-topic", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, msg.Message, adapter.Message.Message)
	assert.Equal(t, msg.Title, adapter.Message.Title)
	assert.Equal(t, msg.Priority, adapter.Message.Priority)
	assert.Equal(t, msg.Tags, adapter.Message.Tags)
	assert.Equal(t, msg.Click, adapter.Message.Click)
	assert.Equal(t, msg.Icon, adapter.Message.Icon)
	assert.Equal(t, msg.Markdown, adapter.Message.Markdown)
}

func TestNtfyAdapter_UnmarshalRequest_QueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test-topic?title=Query+Title&message=Query+Message&priority=4&tags=tag1,tag2&click=https://example.com&icon=https://example.com/icon.png&markdown=true", strings.NewReader(""))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, "Query Title", adapter.Message.Title)
	assert.Equal(t, "Query Message", adapter.Message.Message)
	assert.Equal(t, int32(4), adapter.Message.Priority)
	assert.Equal(t, []string{"tag1", "tag2"}, adapter.Message.Tags)
	assert.Equal(t, "https://example.com", adapter.Message.Click)
	assert.Equal(t, "https://example.com/icon.png", adapter.Message.Icon)
	assert.True(t, adapter.Message.Markdown)
}

func TestNtfyAdapter_UnmarshalRequest_Headers(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test-topic", strings.NewReader(""))
	req.Header.Set("X-Title", "Header Title")
	req.Header.Set("X-Message", "Header Message")
	req.Header.Set("X-Priority", "5")
	req.Header.Set("X-Tags", "tag1,tag2")
	req.Header.Set("X-Click", "https://example.com")
	req.Header.Set("X-Icon", "https://example.com/icon.png")
	req.Header.Set("X-Markdown", "true")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, "Header Title", adapter.Message.Title)
	assert.Equal(t, "Header Message", adapter.Message.Message)
	assert.Equal(t, int32(5), adapter.Message.Priority)
	assert.Equal(t, []string{"tag1", "tag2"}, adapter.Message.Tags)
	assert.Equal(t, "https://example.com", adapter.Message.Click)
	assert.Equal(t, "https://example.com/icon.png", adapter.Message.Icon)
	assert.True(t, adapter.Message.Markdown)
}

func TestNtfyAdapter_UnmarshalRequest_HeadersOverrideQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test-topic?title=Query+Title&priority=1", strings.NewReader(""))
	req.Header.Set("X-Title", "Header Title")
	req.Header.Set("X-Priority", "5")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, "Header Title", adapter.Message.Title, "headers should override query params")
	assert.Equal(t, int32(5), adapter.Message.Priority, "headers should override query params")
}

func TestNtfyAdapter_UnmarshalRequest_QueryParamsOverrideJSON(t *testing.T) {
	msg := NtfyMessage{
		Title:    "JSON Title",
		Message:  "JSON Message",
		Priority: 1,
	}
	bodyBytes, _ := json.Marshal(msg)

	req := httptest.NewRequest(http.MethodPost, "/test-topic?title=Query+Title&priority=5", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, "Query Title", adapter.Message.Title, "query params should override JSON")
	assert.Equal(t, int32(5), adapter.Message.Priority, "query params should override JSON")
	assert.Equal(t, "JSON Message", adapter.Message.Message, "non-overridden should remain")
}

func TestNtfyAdapter_UnmarshalRequest_ShortAliases(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test-topic?t=Title&m=Message&p=4&ta=tag1,tag2&md=1", strings.NewReader(""))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	require.NoError(t, err)

	assert.Equal(t, "Title", adapter.Message.Title)
	assert.Equal(t, "Message", adapter.Message.Message)
	assert.Equal(t, int32(4), adapter.Message.Priority)
	assert.Len(t, adapter.Message.Tags, 2)
	assert.True(t, adapter.Message.Markdown)
}

func TestNtfyAdapter_AsFeedMessage(t *testing.T) {
	adapter := &NtfyAdapter{
		Message: NtfyMessage{
			Topic:    "test-topic",
			Message:  "Test Message",
			Title:    "Test Title",
			Priority: 4,
			Tags:     []string{"tag1", "tag2"},
			Click:    "https://example.com",
			Icon:     "https://example.com/icon.png",
			Markdown: true,
		},
		RawBody:    []byte(`{"message":"Test Message"}`),
		RawHeaders: map[string]string{"Content-Type": "application/json"},
	}

	dto := adapter.AsFeedMessage()

	assert.Equal(t, "test-topic", dto.FeedID)
	assert.Equal(t, "Test Title", *dto.Title)
	assert.Equal(t, "Test Message", *dto.Message)
	assert.Equal(t, int32(4), *dto.Priority)

	// Check metadata contains expected fields
	var metadata map[string]interface{}
	err := json.Unmarshal(dto.Metadata, &metadata)
	require.NoError(t, err)

	tags, ok := metadata["tags"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, tags, 2)
	assert.Equal(t, "https://example.com", metadata["click"])
	assert.Equal(t, "https://example.com/icon.png", metadata["icon"])
	assert.Equal(t, true, metadata["markdown"])
}

func TestNtfyAdapter_AsFeedMessage_PlainTextBody(t *testing.T) {
	adapter := &NtfyAdapter{
		Message: NtfyMessage{
			Topic:   "test-topic",
			Message: "Plain text message",
		},
		RawBody:    []byte("Plain text message"),
		RawHeaders: map[string]string{"Content-Type": "text/plain"},
	}

	dto := adapter.AsFeedMessage()

	// Should wrap plain text in JSON
	var rawRequest map[string]string
	err := json.Unmarshal(dto.RawRequest, &rawRequest)
	require.NoError(t, err)

	assert.Equal(t, "Plain text message", rawRequest["body"])
}

func TestNtfyAdapter_ReadBodyError(t *testing.T) {
	// Create a reader that always returns an error
	errReader := io.NopCloser(&errorReader{})
	req := httptest.NewRequest(http.MethodPost, "/test-topic", errReader)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("topic", "test-topic")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	adapter := &NtfyAdapter{}
	err := adapter.UnmarshalRequest(req)
	assert.Error(t, err)
}

// errorReader is a helper that always returns an error when reading
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
