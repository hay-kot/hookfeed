-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    packing_lists
ADD
    COLUMN tags JSONB NOT NULL DEFAULT '[]',
ADD
    COLUMN days INTEGER NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    packing_lists DROP COLUMN tags,
    DROP COLUMN days;

-- +goose StatementEnd
