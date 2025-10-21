-- +goose Up
-- +goose StatementBegin
-- View that excludes search_vector for application queries
CREATE OR REPLACE VIEW feed_messages_view AS
SELECT
    id,
    feed_slug,
    raw_request,
    raw_headers,
    title,
    message,
    priority,
    logs,
    metadata,
    state,
    state_changed_at,
    received_at,
    processed_at,
    created_at,
    updated_at
FROM feed_messages;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS feed_messages_view;
-- +goose StatementEnd
