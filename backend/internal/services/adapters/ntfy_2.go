package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
)

// ntfyMessage represents a parsed ntfy-compatible message
type ntfyMessage struct {
	Topic    string            `json:"topic"`
	Message  string            `json:"message"`
	Title    string            `json:"title"`
	Priority int32             `json:"priority"`
	Tags     []string          `json:"tags"`
	Click    string            `json:"click"`
	Icon     string            `json:"icon"`
	Actions  []json.RawMessage `json:"actions"`
	Markdown bool              `json:"markdown"`
}

// ParseNtfyMessage parses a ntfy compatible http request and transforms it into
// a validate creaton object or returns an error.
func ParseNtfyMessage(r *http.Request) (dtos.FeedMessageCreate, error) {
	data := dtos.FeedMessageCreateNew()

	data, err := copyHTTPInto(data, r)
	if err != nil {
		return dtos.FeedMessageCreate{}, err
	}

	// Parse ntfy request

	return data, nil
}
