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

func Test_ParseNtfyMessage(t *testing.T) {
	setup := func(topic string, contentType string, body io.Reader) *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/"+topic, body)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		// Set up chi URL params
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("topic", topic)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		return req
	}

	t.Run("plain text body", func(t *testing.T) {
		body := "This is a plain text message"

		req := setup("test-topic", "text/plain", strings.NewReader(body))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, body, dto.Message)
		assert.Equal(t, "test-feed-id", dto.FeedID)
		assert.Equal(t, int32(3), dto.Priority)

		// Verify plain text is wrapped in $body key
		expectedRawRequest := `{"$body": "This is a plain text message"}`
		assert.JSONEq(t, expectedRawRequest, string(dto.RawRequest))
	})

	t.Run("JSON body", func(t *testing.T) {
		msg := ntfyMessage{
			Topic:    "test-topic",
			Message:  "JSON message",
			Title:    "JSON Title",
			Priority: 5,
		}
		bodyBytes, _ := json.Marshal(msg)

		req := setup("test-topic", "application/json", strings.NewReader(string(bodyBytes)))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, msg.Message, dto.Message)
		assert.Equal(t, msg.Title, dto.Title)
		assert.Equal(t, msg.Priority, dto.Priority)

		// Verify the raw request contains the original JSON
		assert.JSONEq(t, string(bodyBytes), string(dto.RawRequest))
	})

	t.Run("query params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test-topic?title=Query+Title&message=Query+Message&priority=4", strings.NewReader(""))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("topic", "test-topic")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, "Query Title", dto.Title)
		assert.Equal(t, "Query Message", dto.Message)
		assert.Equal(t, int32(4), dto.Priority)

		// Verify query params are captured
		expectedQueryParams := `{
			"title": "Query Title",
			"message": "Query Message",
			"priority": "4"
		}`
		assert.JSONEq(t, expectedQueryParams, string(dto.RawQueryParams))
	})

	t.Run("headers", func(t *testing.T) {
		req := setup("test-topic", "", strings.NewReader(""))
		req.Header.Set("X-Title", "Header Title")
		req.Header.Set("X-Message", "Header Message")
		req.Header.Set("X-Priority", "5")

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, "Header Title", dto.Title)
		assert.Equal(t, "Header Message", dto.Message)
		assert.Equal(t, int32(5), dto.Priority)

		// Verify headers are captured
		expectedHeaders := `{
			"X-Title": "Header Title",
			"X-Message": "Header Message",
			"X-Priority": "5"
		}`
		assert.JSONEq(t, expectedHeaders, string(dto.RawHeaders))
	})

	t.Run("headers override query params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test-topic?title=Query+Title&priority=1", strings.NewReader(""))
		req.Header.Set("X-Title", "Header Title")
		req.Header.Set("X-Priority", "5")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("topic", "test-topic")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, "Header Title", dto.Title, "headers should override query params")
		assert.Equal(t, int32(5), dto.Priority, "headers should override query params")

		// Verify both query params and headers are captured
		expectedQueryParams := `{
			"title": "Query Title",
			"priority": "1"
		}`
		assert.JSONEq(t, expectedQueryParams, string(dto.RawQueryParams))

		expectedHeaders := `{
			"X-Title": "Header Title",
			"X-Priority": "5"
		}`
		assert.JSONEq(t, expectedHeaders, string(dto.RawHeaders))
	})

	t.Run("query params override JSON", func(t *testing.T) {
		msg := ntfyMessage{
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

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, "Query Title", dto.Title, "query params should override JSON")
		assert.Equal(t, int32(5), dto.Priority, "query params should override JSON")
		assert.Equal(t, "JSON Message", dto.Message, "non-overridden should remain")

		// Verify raw request still contains original JSON
		assert.JSONEq(t, string(bodyBytes), string(dto.RawRequest))

		// Verify query params are captured
		expectedQueryParams := `{
			"title": "Query Title",
			"priority": "5"
		}`
		assert.JSONEq(t, expectedQueryParams, string(dto.RawQueryParams))
	})

	t.Run("short aliases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test-topic?t=Title&m=Message&p=4", strings.NewReader(""))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("topic", "test-topic")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, "Title", dto.Title)
		assert.Equal(t, "Message", dto.Message)
		assert.Equal(t, int32(4), dto.Priority)

		// Verify query params are captured with their short alias names
		expectedQueryParams := `{
			"t": "Title",
			"m": "Message",
			"p": "4"
		}`
		assert.JSONEq(t, expectedQueryParams, string(dto.RawQueryParams))
	})

	t.Run("default priority", func(t *testing.T) {
		req := setup("test-topic", "", strings.NewReader("test message"))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		assert.Equal(t, int32(3), dto.Priority, "default priority should be 3")

		// Verify plain text body is captured
		expectedRawRequest := `{"$body": "test message"}`
		assert.JSONEq(t, expectedRawRequest, string(dto.RawRequest))
	})

	t.Run("raw data capture", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test-topic?foo=bar", strings.NewReader("test body"))
		req.Header.Set("X-Custom-Header", "custom-value")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("topic", "test-topic")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		dto, err := ParseNtfyMessage(req, "test-feed-id")
		require.NoError(t, err)

		// Verify raw request wraps plain text
		expectedRawRequest := `{"$body": "test body"}`
		assert.JSONEq(t, expectedRawRequest, string(dto.RawRequest))

		// Verify query params are captured
		expectedQueryParams := `{"foo": "bar"}`
		assert.JSONEq(t, expectedQueryParams, string(dto.RawQueryParams))

		// Verify headers are captured (single values are unwrapped from arrays)
		expectedHeaders := `{"X-Custom-Header": "custom-value"}`
		assert.JSONEq(t, expectedHeaders, string(dto.RawHeaders))
	})
}
