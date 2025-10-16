-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feed_messages (
    -- table column defaults
    id UUID DEFAULT uuid_generate_v7 () PRIMARY KEY,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feed_messages;

-- +goose StatementEnd
