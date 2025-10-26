-- +goose Up
-- +goose StatementBegin
-- Drop the view first
DROP VIEW IF EXISTS feed_messages_view;

-- Add the new column to the table
ALTER TABLE feed_messages ADD COLUMN raw_query_params JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Recreate the view with the new column
CREATE VIEW feed_messages_view AS
SELECT
    id,
    feed_slug,
    raw_request,
    raw_headers,
    raw_query_params,
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
-- Drop the view
DROP VIEW IF EXISTS feed_messages_view;

-- Remove the column
ALTER TABLE feed_messages DROP COLUMN IF EXISTS raw_query_params;

-- Recreate the view without raw_query_params
CREATE VIEW feed_messages_view AS
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
