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

	// Capture raw headers for reference
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

	// Helper to check if a field is empty or null
	isEmpty := func(data json.RawMessage) bool {
		return len(data) == 0 || string(data) == "null" || string(data) == "{}"
	}

	// If any raw fields are missing, ensure they're properly initialized
	needsInit := isEmpty(ra.CreateDTO.RawRequest) ||
		isEmpty(ra.CreateDTO.RawHeaders) ||
		isEmpty(ra.CreateDTO.RawQueryParams)

	if needsInit {
		// Determine the body to use for RawRequest
		var bodyData any
		if isEmpty(ra.CreateDTO.RawRequest) {
			if json.Valid(bodyBytes) {
				json.Unmarshal(bodyBytes, &bodyData)
			} else {
				bodyData = map[string]string{"body": string(bodyBytes)}
			}
		} else {
			// Use existing RawRequest
			json.Unmarshal(ra.CreateDTO.RawRequest, &bodyData)
		}

		// Convert headers to http.Header format
		headers := make(map[string][]string)
		if isEmpty(ra.CreateDTO.RawHeaders) {
			for k, v := range ra.RawHeaders {
				headers[k] = []string{v}
			}
		} else {
			// Keep existing headers
			json.Unmarshal(ra.CreateDTO.RawHeaders, &headers)
		}

		// Get query params
		var queryParams map[string][]string
		if isEmpty(ra.CreateDTO.RawQueryParams) {
			queryParams = r.URL.Query()
		} else {
			// Keep existing query params
			json.Unmarshal(ra.CreateDTO.RawQueryParams, &queryParams)
		}

		// Use constructor to ensure proper initialization
		initialized, err := dtos.NewFeedMessageCreateFromHTTP(
			ra.CreateDTO.FeedID,
			bodyData,
			headers,
			queryParams,
		)
		if err != nil {
			return err
		}

		// Merge initialized values with user-provided values
		if isEmpty(ra.CreateDTO.RawRequest) {
			ra.CreateDTO.RawRequest = initialized.RawRequest
		}
		if isEmpty(ra.CreateDTO.RawHeaders) {
			ra.CreateDTO.RawHeaders = initialized.RawHeaders
		}
		if isEmpty(ra.CreateDTO.RawQueryParams) {
			ra.CreateDTO.RawQueryParams = initialized.RawQueryParams
		}
		if len(ra.CreateDTO.Logs) == 0 {
			ra.CreateDTO.Logs = initialized.Logs
		}
		if isEmpty(ra.CreateDTO.Metadata) {
			ra.CreateDTO.Metadata = initialized.Metadata
		}
	}

	return nil
}

// AsFeedMessage returns the FeedMessageCreate DTO
func (ra *RawAdapter) AsFeedMessage() dtos.FeedMessageCreate {
	return ra.CreateDTO
}
