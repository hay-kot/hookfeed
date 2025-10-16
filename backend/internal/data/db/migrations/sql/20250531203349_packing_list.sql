-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS packing_lists (
    id UUID DEFAULT uuid_generate_v7 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Relations
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    -- Fields
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS packing_list_items (
    id UUID DEFAULT uuid_generate_v7 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Relations
    packing_list_id UUID NOT NULL REFERENCES packing_lists(id) ON DELETE CASCADE,
    -- Fields
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    is_packed BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT NOT NULL
);

CREATE INDEX idx_packing_lists_user_id ON packing_lists(user_id);
CREATE INDEX idx_packing_lists_created_at ON packing_lists(created_at DESC);
CREATE INDEX idx_items_packing_list_id ON packing_list_items(packing_list_id);
CREATE INDEX idx_items_category ON packing_list_items(category);
CREATE INDEX idx_items_is_packed ON packing_list_items(is_packed);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_packing_lists_updated_at
    BEFORE UPDATE ON packing_lists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_items_updated_at
    BEFORE UPDATE ON packing_list_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS packing_lists;

DROP TABLE IF EXISTS packing_list_items;

-- +goose StatementEnd
