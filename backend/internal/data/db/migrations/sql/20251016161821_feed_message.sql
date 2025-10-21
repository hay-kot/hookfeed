-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Feed Messages table
-- Note: feeds are stored in cache, not in the database
CREATE TABLE IF NOT EXISTS feed_messages (
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY,
    feed_slug VARCHAR(255) NOT NULL,
    raw_request JSONB NOT NULL,
    raw_headers JSONB NOT NULL,
    title VARCHAR(500),
    message TEXT,
    priority INTEGER DEFAULT 3,
    logs TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}'::jsonb,
    state VARCHAR(20) DEFAULT 'new',
    state_changed_at TIMESTAMP,
    received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feed_messages_feed_slug ON feed_messages(feed_slug);
CREATE INDEX IF NOT EXISTS idx_feed_messages_received_at ON feed_messages(received_at DESC);
CREATE INDEX IF NOT EXISTS idx_feed_messages_priority ON feed_messages(priority);
CREATE INDEX IF NOT EXISTS idx_feed_messages_state ON feed_messages(state);
CREATE INDEX IF NOT EXISTS idx_feed_messages_feed_received ON feed_messages(feed_slug, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_feed_messages_feed_state ON feed_messages(feed_slug, state);

-- Full-text search support
ALTER TABLE feed_messages ADD COLUMN search_vector tsvector;
CREATE INDEX feed_messages_search_idx ON feed_messages USING GIN(search_vector);

CREATE OR REPLACE FUNCTION feed_messages_search_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.message, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(array_to_string(NEW.logs, ' '), '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER feed_messages_search_update
    BEFORE INSERT OR UPDATE ON feed_messages
    FOR EACH ROW EXECUTE FUNCTION feed_messages_search_trigger();

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_feed_messages_updated_at BEFORE UPDATE ON feed_messages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS feed_messages_search_update ON feed_messages;
DROP FUNCTION IF EXISTS feed_messages_search_trigger();
DROP TRIGGER IF EXISTS update_feed_messages_updated_at ON feed_messages;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS feed_messages;

-- +goose StatementEnd
