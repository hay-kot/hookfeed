package adapters

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
)

// RawAdapter adapts raw FeedMessageCreate requests
type RawAdapter struct {
	CreateDTO  dtos.FeedMessageCreate
	RawBody    []byte
	RawHeaders map[string]string
}

// UnmarshalRequest parses a raw HTTP request into a FeedMessageCreate
func (ra *RawAdapter) UnmarshalRequest(r *http.Request) error {
	// Read body once
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	ra.RawBody = bodyBytes

	// Capture raw headers
	ra.RawHeaders = make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			ra.RawHeaders[k] = v[0]
		}
	}

	// Unmarshal the body into FeedMessageCreate
	if err := json.Unmarshal(bodyBytes, &ra.CreateDTO); err != nil {
		return err
	}

	// If RawRequest is not provided, use the body as RawRequest
	if ra.CreateDTO.RawRequest == nil || len(ra.CreateDTO.RawRequest) == 0 || string(ra.CreateDTO.RawRequest) == "null" {
		if json.Valid(bodyBytes) {
			ra.CreateDTO.RawRequest = json.RawMessage(bodyBytes)
		} else {
			// Wrap plain text in JSON
			wrapped := map[string]string{"body": string(bodyBytes)}
			ra.CreateDTO.RawRequest, _ = json.Marshal(wrapped)
		}
	}

	// If RawHeaders is not provided, capture them
	if ra.CreateDTO.RawHeaders == nil || len(ra.CreateDTO.RawHeaders) == 0 || string(ra.CreateDTO.RawHeaders) == "null" {
		headersJSON, _ := json.Marshal(ra.RawHeaders)
		ra.CreateDTO.RawHeaders = json.RawMessage(headersJSON)
	}

	return nil
}

// AsFeedMessage returns the FeedMessageCreate DTO
func (ra *RawAdapter) AsFeedMessage() dtos.FeedMessageCreate {
	return ra.CreateDTO
}
